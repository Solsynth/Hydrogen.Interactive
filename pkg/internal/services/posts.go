package services

import (
	"errors"
	"fmt"
	"time"

	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/proto"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func FilterPostWithCategory(tx *gorm.DB, alias string) *gorm.DB {
	return tx.Joins("JOIN post_categories ON posts.id = post_categories.post_id").
		Joins("JOIN post_categories ON post_categories.id = post_categories.category_id").
		Where("post_categories.alias = ?", alias)
}

func FilterPostWithTag(tx *gorm.DB, alias string) *gorm.DB {
	return tx.Joins("JOIN post_tags ON posts.id = post_tags.post_id").
		Joins("JOIN post_tags ON post_tags.id = post_tags.category_id").
		Where("post_tags.alias = ?", alias)
}

func FilterPostWithRealm(tx *gorm.DB, id uint) *gorm.DB {
	if id > 0 {
		return tx.Where("realm_id = ?", id)
	} else {
		return tx.Where("realm_id IS NULL")
	}
}

func FilterPostReply(tx *gorm.DB, replyTo ...uint) *gorm.DB {
	if len(replyTo) > 0 && replyTo[0] > 0 {
		return tx.Where("reply_id = ?", replyTo[0])
	} else {
		return tx.Where("reply_id IS NULL")
	}
}

func FilterPostWithPublishedAt(tx *gorm.DB, date time.Time) *gorm.DB {
	return tx.Where("published_at <= ? OR published_at IS NULL", date)
}

func FilterPostWithAuthorDraft(tx *gorm.DB, uid uint) *gorm.DB {
	return tx.Where("author_id = ? AND is_draft = ?", uid, true)
}

func FilterPostDraft(tx *gorm.DB) *gorm.DB {
	return tx.Where("is_draft = ? OR is_draft IS NULL", false)
}

func GetPostWithAlias(tx *gorm.DB, alias string, ignoreLimitation ...bool) (models.Post, error) {
	if len(ignoreLimitation) == 0 || !ignoreLimitation[0] {
		tx = FilterPostWithPublishedAt(tx, time.Now())
	}

	var item models.Post
	if err := tx.
		Where("alias = ?", alias).
		Preload("Tags").
		Preload("Categories").
		Preload("Realm").
		Preload("Author").
		Preload("ReplyTo").
		Preload("ReplyTo.Author").
		Preload("RepostTo").
		Preload("RepostTo.Author").
		First(&item).Error; err != nil {
		return item, err
	}

	return item, nil
}

func GetPost(tx *gorm.DB, id uint, ignoreLimitation ...bool) (models.Post, error) {
	if len(ignoreLimitation) == 0 || !ignoreLimitation[0] {
		tx = FilterPostWithPublishedAt(tx, time.Now())
	}

	var item models.Post
	if err := tx.
		Where("id = ?", id).
		Preload("Tags").
		Preload("Categories").
		Preload("Realm").
		Preload("Author").
		Preload("ReplyTo").
		Preload("ReplyTo.Author").
		Preload("RepostTo").
		Preload("RepostTo.Author").
		First(&item).Error; err != nil {
		return item, err
	}

	return item, nil
}

func CountPost(tx *gorm.DB) (int64, error) {
	var count int64
	if err := tx.Model(&models.Post{}).Count(&count).Error; err != nil {
		return count, err
	}

	return count, nil
}

func CountPostReply(id uint) int64 {
	var count int64
	if err := database.C.Model(&models.Post{}).
		Where("reply_id = ?", id).
		Count(&count).Error; err != nil {
		return 0
	}

	return count
}

func CountPostReactions(id uint) int64 {
	var count int64
	if err := database.C.Model(&models.Reaction{}).
		Where("post_id = ?", id).
		Count(&count).Error; err != nil {
		return 0
	}

	return count
}

func ListPost(tx *gorm.DB, take int, offset int, noReact ...bool) ([]*models.Post, error) {
	if take > 20 {
		take = 20
	}

	var items []*models.Post
	if err := tx.
		Limit(take).Offset(offset).
		Order("created_at DESC").
		Preload("Realm").
		Preload("Author").
		Preload("ReplyTo").
		Preload("ReplyTo.Author").
		Preload("RepostTo").
		Preload("RepostTo.Author").
		Find(&items).Error; err != nil {
		return items, err
	}

	idx := lo.Map(items, func(item *models.Post, index int) uint {
		return item.ID
	})

	// Load reactions
	if len(noReact) <= 0 || !noReact[0] {
		if mapping, err := BatchListResourceReactions(database.C.Where("post_id IN ?", idx), "post_id"); err != nil {
			return items, err
		} else {
			itemMap := lo.SliceToMap(items, func(item *models.Post) (uint, *models.Post) {
				return item.ID, item
			})

			for k, v := range mapping {
				if post, ok := itemMap[k]; ok {
					post.ReactionList = v
				}
			}
		}
	}

	// Load replies
	if len(noReact) <= 0 || !noReact[0] {
		var replies []struct {
			PostID uint
			Count  int64
		}

		if err := database.C.Model(&models.Post{}).
			Select("id as post_id, COUNT(id) as count").
			Where("reply_id IN (?)", idx).
			Group("post_id").
			Scan(&replies).Error; err != nil {
			return items, err
		}

		itemMap := lo.SliceToMap(items, func(item *models.Post) (uint, *models.Post) {
			return item.ID, item
		})

		list := map[uint]int64{}
		for _, info := range replies {
			list[info.PostID] = info.Count
		}

		for k, v := range list {
			if post, ok := itemMap[k]; ok {
				post.ReplyCount = v
			}
		}
	}

	return items, nil
}

func EnsurePostCategoriesAndTags(item models.Post) (models.Post, error) {
	var err error
	for idx, category := range item.Categories {
		item.Categories[idx], err = GetCategory(category.Alias)
		if err != nil {
			return item, err
		}
	}
	for idx, tag := range item.Tags {
		item.Tags[idx], err = GetTagOrCreate(tag.Alias, tag.Name)
		if err != nil {
			return item, err
		}
	}
	return item, nil
}

func NewPost(user models.Account, item models.Post) (models.Post, error) {
	item, err := EnsurePostCategoriesAndTags(item)
	if err != nil {
		return item, err
	}

	if item.RealmID != nil {
		_, err := GetRealmMember(*item.RealmID, user.ExternalID)
		if err != nil {
			return item, fmt.Errorf("you aren't a part of that realm: %v", err)
		}
	}

	if err := database.C.Save(&item).Error; err != nil {
		return item, err
	}

	// Notify the original poster its post has been replied
	if item.ReplyID != nil {
		var op models.Post
		if err := database.C.
			Where("id = ?", item.ReplyID).
			Preload("Author").
			First(&op).Error; err == nil {
			if op.Author.ID != user.ID {
				postUrl := fmt.Sprintf("https://%s/posts/%s", viper.GetString("domain"), item.Alias)
				err := NotifyPosterAccount(
					op.Author,
					fmt.Sprintf("%s replied you", user.Nick),
					fmt.Sprintf("%s (%s) replied your post #%s.", user.Nick, user.Name, op.Alias),
					&proto.NotifyLink{Label: "Related post", Url: postUrl},
				)
				if err != nil {
					log.Error().Err(err).Msg("An error occurred when notifying user...")
				}
			}
		}
	}

	return item, nil
}

func EditPost(item models.Post) (models.Post, error) {
	item, err := EnsurePostCategoriesAndTags(item)
	if err != nil {
		return item, err
	}

	err = database.C.Save(&item).Error

	return item, err
}

func DeletePost(item models.Post) error {
	return database.C.Delete(&item).Error
}

func ReactPost(user models.Account, reaction models.Reaction) (bool, models.Reaction, error) {
	if err := database.C.Where(reaction).First(&reaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var op models.Post
			if err := database.C.
				Where("id = ?", reaction.PostID).
				Preload("Author").
				First(&op).Error; err == nil {
				if op.Author.ID != user.ID {
					postUrl := fmt.Sprintf("https://%s/posts/%s", viper.GetString("domain"), op.Alias)
					err := NotifyPosterAccount(
						op.Author,
						fmt.Sprintf("%s reacted your post", user.Nick),
						fmt.Sprintf("%s (%s) reacted your post a %s", user.Nick, user.Name, reaction.Symbol),
						&proto.NotifyLink{Label: "Related post", Url: postUrl},
					)
					if err != nil {
						log.Error().Err(err).Msg("An error occurred when notifying user...")
					}
				}
			}

			return true, reaction, database.C.Save(&reaction).Error
		} else {
			return true, reaction, err
		}
	} else {
		return false, reaction, database.C.Delete(&reaction).Error
	}
}
