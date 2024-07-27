package services

import (
	"errors"
	"fmt"
	"time"

	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func FilterPostWithCategory(tx *gorm.DB, alias string) *gorm.DB {
	prefix := viper.GetString("database.prefix")
	return tx.Joins(fmt.Sprintf("JOIN %spost_categories ON %sposts.id = %spost_categories.post_id", prefix, prefix, prefix)).
		Joins(fmt.Sprintf("JOIN %scategories ON %scategories.id = %spost_categories.category_id", prefix, prefix, prefix)).
		Where(fmt.Sprintf("%scategories.alias = ?", prefix), alias)
}

func FilterPostWithTag(tx *gorm.DB, alias string) *gorm.DB {
	prefix := viper.GetString("database.prefix")
	return tx.Joins(fmt.Sprintf("JOIN %spost_tags ON %sposts.id = %spost_tags.post_id", prefix, prefix, prefix)).
		Joins(fmt.Sprintf("JOIN %stags ON %stags.id = %spost_tags.tag_id", prefix, prefix, prefix)).
		Where(fmt.Sprintf("%stags.alias = ?", prefix), alias)
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
	return tx.
		Where("published_at <= ? OR published_at IS NULL", date).
		Where("published_until > ? OR published_until IS NULL", date)
}

func FilterPostWithAuthorDraft(tx *gorm.DB, uid uint) *gorm.DB {
	return tx.Where("author_id = ? AND is_draft = ?", uid, true)
}

func FilterPostDraft(tx *gorm.DB) *gorm.DB {
	return tx.Where("is_draft = ? OR is_draft IS NULL", false)
}

func PreloadGeneral(tx *gorm.DB) *gorm.DB {
	return tx.
		Preload("Tags").
		Preload("Categories").
		Preload("Realm").
		Preload("Author").
		Preload("ReplyTo").
		Preload("ReplyTo.Author").
		Preload("ReplyTo.Tags").
		Preload("ReplyTo.Categories").
		Preload("RepostTo").
		Preload("RepostTo.Author").
		Preload("RepostTo.Tags").
		Preload("RepostTo.Categories")
}

func GetPost(tx *gorm.DB, id uint, ignoreLimitation ...bool) (models.Post, error) {
	if len(ignoreLimitation) == 0 || !ignoreLimitation[0] {
		tx = FilterPostWithPublishedAt(tx, time.Now())
	}

	var item models.Post
	if err := PreloadGeneral(tx).
		Where("id = ?", id).
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

func ListPost(tx *gorm.DB, take int, offset int, order any, noReact ...bool) ([]*models.Post, error) {
	if take > 100 {
		take = 100
	}

	var items []*models.Post
	if err := PreloadGeneral(tx).
		Limit(take).Offset(offset).
		Order(order).
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
					post.Metric = models.PostMetric{
						ReactionList: v,
					}
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
				post.Metric = models.PostMetric{
					ReactionList: post.Metric.ReactionList,
					ReplyCount:   v,
				}
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
	if !item.IsDraft && item.PublishedAt == nil {
		item.PublishedAt = lo.ToPtr(time.Now())
	}

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

	item.EditedAt = lo.ToPtr(time.Now())

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
				err = NotifyPosterAccount(
					op.Author,
					"Post got replied",
					fmt.Sprintf("%s (%s) replied your post.", user.Nick, user.Name),
					lo.ToPtr(fmt.Sprintf("%s replied you", user.Nick)),
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
	var op models.Post
	if err := database.C.
		Where("id = ?", reaction.PostID).
		Preload("Author").
		First(&op).Error; err != nil {
		return true, reaction, err
	}

	if err := database.C.Where(reaction).First(&reaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if op.Author.ID != user.ID {
				err = NotifyPosterAccount(
					op.Author,
					"Post got reacted",
					fmt.Sprintf("%s (%s) reacted your post a %s.", user.Nick, user.Name, reaction.Symbol),
					lo.ToPtr(fmt.Sprintf("%s reacted you", user.Nick)),
				)
				if err != nil {
					log.Error().Err(err).Msg("An error occurred when notifying user...")
				}
			}

			err = database.C.Save(&reaction).Error
			if err == nil && reaction.Attitude != models.AttitudeNeutral {
				_ = ModifyPosterVoteCount(op.Author, reaction.Attitude == models.AttitudePositive, 1)

				if reaction.Attitude == models.AttitudePositive {
					op.TotalUpvote++
				} else {
					op.TotalDownvote++
				}
				database.C.Save(&op)
			}

			return true, reaction, err
		} else {
			return true, reaction, err
		}
	} else {
		err = database.C.Delete(&reaction).Error
		if err == nil && reaction.Attitude != models.AttitudeNeutral {
			_ = ModifyPosterVoteCount(op.Author, reaction.Attitude == models.AttitudePositive, -1)

			if reaction.Attitude == models.AttitudePositive {
				op.TotalUpvote--
			} else {
				op.TotalDownvote--
			}
			database.C.Save(&op)
		}

		return false, reaction, err
	}
}

func PinPost(post models.Post) (bool, error) {
	if post.PinnedAt != nil {
		post.PinnedAt = nil
	} else {
		post.PinnedAt = lo.ToPtr(time.Now())
	}

	if err := database.C.Save(&post).Error; err != nil {
		return post.PinnedAt != nil, err
	}
	return post.PinnedAt != nil, nil
}
