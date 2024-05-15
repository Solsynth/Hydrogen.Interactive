package services

import (
	"errors"
	"fmt"
	"git.solsynth.dev/hydrogen/interactive/pkg/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/models"
	"git.solsynth.dev/hydrogen/passport/pkg/grpc/proto"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"time"
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

func FilterWithRealm(tx *gorm.DB, id uint) *gorm.DB {
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
	return tx.Where("published_at <= ? AND published_at IS NULL", date)
}

func GetPostWithAlias(alias string, ignoreLimitation ...bool) (models.Post, error) {
	tx := database.C
	if len(ignoreLimitation) == 0 || !ignoreLimitation[0] {
		tx = FilterPostWithPublishedAt(tx, time.Now())
	}

	var item models.Post
	if err := tx.
		Where("alias = ?", alias).
		Preload("Author").
		Preload("Attachments").
		Preload("ReplyTo").
		Preload("ReplyTo.Author").
		Preload("ReplyTo.Attachments").
		Preload("RepostTo").
		Preload("RepostTo.Author").
		Preload("RepostTo.Attachments").
		First(&item).Error; err != nil {
		return item, err
	}

	return item, nil
}

func GetPost(id uint, ignoreLimitation ...bool) (models.Post, error) {
	tx := database.C
	if len(ignoreLimitation) == 0 || !ignoreLimitation[0] {
		tx = FilterPostWithPublishedAt(tx, time.Now())
	}

	var item models.Post
	if err := tx.
		Where("id = ?", id).
		Preload("Author").
		Preload("Attachments").
		Preload("ReplyTo").
		Preload("ReplyTo.Author").
		Preload("ReplyTo.Attachments").
		Preload("RepostTo").
		Preload("RepostTo.Author").
		Preload("RepostTo.Attachments").
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

func ListPostReactions(id uint) (map[string]int64, error) {
	var reactions []struct {
		Symbol string
		Count  int64
	}

	if err := database.C.Model(&models.Reaction{}).
		Select("symbol, COUNT(id) as count").
		Where("post_id = ?", id).
		Group("symbol").
		Scan(&reactions).Error; err != nil {
		return map[string]int64{}, err
	}

	return lo.SliceToMap(reactions, func(item struct {
		Symbol string
		Count  int64
	},
	) (string, int64) {
		return item.Symbol, item.Count
	}), nil
}

func ListPost(tx *gorm.DB, take int, offset int, noReact ...bool) ([]models.Post, error) {
	if take > 20 {
		take = 20
	}

	var items []models.Post
	if err := tx.
		Limit(take).Offset(offset).
		Order("created_at DESC").
		Preload("Author").
		Preload("Attachments").
		Preload("ReplyTo").
		Preload("ReplyTo.Author").
		Preload("ReplyTo.Attachments").
		Preload("RepostTo").
		Preload("RepostTo.Author").
		Preload("RepostTo.Attachments").
		Find(&items).Error; err != nil {
		return items, err
	}

	idx := lo.Map(items, func(item models.Post, index int) uint {
		return item.ID
	})

	if len(noReact) <= 0 || !noReact[0] {
		var reactions []struct {
			PostID uint
			Symbol string
			Count  int64
		}

		if err := database.C.Model(&models.Reaction{}).
			Select("post_id as post_id, symbol, COUNT(id) as count").
			Where("post_id IN (?)", idx).
			Group("post_id, symbol").
			Scan(&reactions).Error; err != nil {
			return items, err
		}

		itemMap := lo.SliceToMap(items, func(item models.Post) (uint, models.Post) {
			return item.ID, item
		})

		list := map[uint]map[string]int64{}
		for _, info := range reactions {
			if _, ok := list[info.PostID]; !ok {
				list[info.PostID] = make(map[string]int64)
			}
			list[info.PostID][info.Symbol] = info.Count
		}

		for k, v := range list {
			if post, ok := itemMap[k]; ok {
				post.ReactionList = v
			}
		}
	}

	return items, nil
}

func InitPostCategoriesAndTags(item models.Post) (models.Post, error) {
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
	item, err := InitPostCategoriesAndTags(item)
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
		go func() {
			var op models.Post
			if err := database.C.
				Where("id = ?", item.ReplyID).
				Preload("Author").
				First(&op).Error; err == nil {
				if op.Author.ID != user.ID {
					postUrl := fmt.Sprintf("https://%s/posts/%s", viper.GetString("domain"), item.Alias)
					err := NotifyAccount(
						op.Author,
						fmt.Sprintf("%s replied you", user.Nick),
						fmt.Sprintf("%s (%s) replied your post #%s.", user.Nick, user.Name, op.Alias),
						false,
						&proto.NotifyLink{Label: "Related post", Url: postUrl},
					)
					if err != nil {
						log.Error().Err(err).Msg("An error occurred when notifying user...")
					}
				}
			}
		}()
	}

	return item, nil
}

func EditPost(item models.Post) (models.Post, error) {
	item, err := InitPostCategoriesAndTags(item)
	if err != nil {
		return item, err
	}

	err = database.C.Save(&item).Error

	return item, err
}

func DeletePost(item models.Post) error {
	return database.C.Delete(&item).Error
}

func ReactPost(reaction models.Reaction) (bool, models.Reaction, error) {
	if err := database.C.Where(reaction).First(&reaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, reaction, database.C.Save(&reaction).Error
		} else {
			return true, reaction, err
		}
	} else {
		return false, reaction, database.C.Delete(&reaction).Error
	}
}
