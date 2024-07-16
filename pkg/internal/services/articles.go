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

func FilterArticleWithCategory(tx *gorm.DB, alias string) *gorm.DB {
	prefix := viper.GetString("database.prefix")
	return tx.Joins(fmt.Sprintf("JOIN %sarticle_categories ON %sarticles.id = %sarticle_categories.article_id", prefix, prefix, prefix)).
		Joins(fmt.Sprintf("JOIN %scategories ON %scategories.id = %sarticle_categories.category_id", prefix, prefix, prefix)).
		Where(fmt.Sprintf("%scategories.alias = ?", prefix), alias)
}

func FilterArticleWithTag(tx *gorm.DB, alias string) *gorm.DB {
	prefix := viper.GetString("database.prefix")
	return tx.Joins(fmt.Sprintf("JOIN %sarticle_tags ON %sarticles.id = %sarticle_tags.article_id", prefix, prefix, prefix)).
		Joins(fmt.Sprintf("JOIN %stags ON %stags.id = %sarticle_tags.tag_id", prefix, prefix, prefix)).
		Where(fmt.Sprintf("%stags.alias = ?", prefix), alias)
}

func FilterArticleWithRealm(tx *gorm.DB, id uint) *gorm.DB {
	if id > 0 {
		return tx.Where("realm_id = ?", id)
	} else {
		return tx.Where("realm_id IS NULL")
	}
}

func FilterArticleWithPublishedAt(tx *gorm.DB, date time.Time) *gorm.DB {
	return tx.Where("published_at <= ? OR published_at IS NULL", date)
}

func FilterArticleWithAuthorDraft(tx *gorm.DB, uid uint) *gorm.DB {
	return tx.Where("author_id = ? AND is_draft = ?", uid, true)
}

func FilterArticleDraft(tx *gorm.DB) *gorm.DB {
	return tx.Where("is_draft = ? OR is_draft IS NULL", false)
}

func GetArticleWithAlias(tx *gorm.DB, alias string, ignoreLimitation ...bool) (models.Article, error) {
	if len(ignoreLimitation) == 0 || !ignoreLimitation[0] {
		tx = FilterArticleWithPublishedAt(tx, time.Now())
	}

	var item models.Article
	if err := tx.
		Where("alias = ?", alias).
		Preload("Tags").
		Preload("Categories").
		Preload("Realm").
		Preload("Author").
		First(&item).Error; err != nil {
		return item, err
	}

	return item, nil
}

func GetArticle(tx *gorm.DB, id uint, ignoreLimitation ...bool) (models.Article, error) {
	if len(ignoreLimitation) == 0 || !ignoreLimitation[0] {
		tx = FilterArticleWithPublishedAt(tx, time.Now())
	}

	var item models.Article
	if err := tx.
		Where("id = ?", id).
		Preload("Tags").
		Preload("Categories").
		Preload("Realm").
		Preload("Author").
		First(&item).Error; err != nil {
		return item, err
	}

	return item, nil
}

func CountArticle(tx *gorm.DB) (int64, error) {
	var count int64
	if err := tx.Model(&models.Article{}).Count(&count).Error; err != nil {
		return count, err
	}

	return count, nil
}

func CountArticleReactions(id uint) int64 {
	var count int64
	if err := database.C.Model(&models.Reaction{}).
		Where("article_id = ?", id).
		Count(&count).Error; err != nil {
		return 0
	}

	return count
}

func ListArticle(tx *gorm.DB, take int, offset int, noReact ...bool) ([]*models.Article, error) {
	if take > 100 {
		take = 100
	}

	var items []*models.Article
	if err := tx.
		Limit(take).Offset(offset).
		Order("created_at DESC").
		Preload("Tags").
		Preload("Categories").
		Preload("Realm").
		Preload("Author").
		Find(&items).Error; err != nil {
		return items, err
	}

	idx := lo.Map(items, func(item *models.Article, index int) uint {
		return item.ID
	})

	// Load reactions
	if len(noReact) <= 0 || !noReact[0] {
		if mapping, err := BatchListResourceReactions(database.C.Where("article_id IN ?", idx), "article_id"); err != nil {
			return items, err
		} else {
			itemMap := lo.SliceToMap(items, func(item *models.Article) (uint, *models.Article) {
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

	return items, nil
}

func EnsureArticleCategoriesAndTags(item models.Article) (models.Article, error) {
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

func NewArticle(user models.Account, item models.Article) (models.Article, error) {
	item.Language = DetectLanguage(&item.Content)

	item, err := EnsureArticleCategoriesAndTags(item)
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

	return item, nil
}

func EditArticle(item models.Article) (models.Article, error) {
	item.Language = DetectLanguage(&item.Content)
	item, err := EnsureArticleCategoriesAndTags(item)
	if err != nil {
		return item, err
	}

	err = database.C.Save(&item).Error

	return item, err
}

func DeleteArticle(item models.Article) error {
	return database.C.Delete(&item).Error
}

func ReactArticle(user models.Account, reaction models.Reaction) (bool, models.Reaction, error) {
	if err := database.C.Where(reaction).First(&reaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var op models.Article
			if err := database.C.
				Where("id = ?", reaction.ArticleID).
				Preload("Author").
				First(&op).Error; err == nil {
				if op.Author.ID != user.ID {
					err = NotifyPosterAccount(
						op.Author,
						"Article got reacted",
						fmt.Sprintf("%s (%s) reacted your article a %s", user.Nick, user.Name, reaction.Symbol),
						lo.ToPtr(fmt.Sprintf("%s reacted your article", user.Nick)),
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
