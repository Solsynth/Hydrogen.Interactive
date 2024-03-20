package services

import (
	"errors"
	"fmt"
	"time"

	"git.solsynth.dev/hydrogen/identity/pkg/grpc/proto"
	"git.solsynth.dev/hydrogen/interactive/pkg/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/models"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type PostTypeContext struct {
	Tx *gorm.DB

	TableName  string
	ColumnName string
	CanReply   bool
	CanRepost  bool
}

func (v *PostTypeContext) FilterWithCategory(alias string) *PostTypeContext {
	name := v.ColumnName
	v.Tx.Joins(fmt.Sprintf("JOIN %s_categories ON %s.id = %s_categories.%s_id", name, v.TableName, name, name)).
		Joins(fmt.Sprintf("JOIN %s_categories ON %s_categories.id = %s_categories.category_id", name, name, name)).
		Where(name+"_categories.alias = ?", alias)
	return v
}

func (v *PostTypeContext) FilterWithTag(alias string) *PostTypeContext {
	name := v.ColumnName
	v.Tx.Joins(fmt.Sprintf("JOIN %s_tags ON %s.id = %s_tags.%s_id", name, v.TableName, name, name)).
		Joins(fmt.Sprintf("JOIN %s_tags ON %s_tags.id = %s_tags.category_id", name, name, name)).
		Where(name+"_tags.alias = ?", alias)
	return v
}

func (v *PostTypeContext) FilterPublishedAt(date time.Time) *PostTypeContext {
	v.Tx.Where("published_at <= ? AND published_at IS NULL", date)
	return v
}

func (v *PostTypeContext) FilterRealm(id uint) *PostTypeContext {
	if id > 0 {
		v.Tx = v.Tx.Where("realm_id = ?", id)
	} else {
		v.Tx = v.Tx.Where("realm_id IS NULL")
	}
	return v
}

func (v *PostTypeContext) FilterAuthor(id uint) *PostTypeContext {
	v.Tx = v.Tx.Where("author_id = ?", id)
	return v
}

func (v *PostTypeContext) FilterReply(condition bool) *PostTypeContext {
	if condition {
		v.Tx = v.Tx.Where("reply_id IS NOT NULL")
	} else {
		v.Tx = v.Tx.Where("reply_id IS NULL")
	}
	return v
}

func (v *PostTypeContext) SortCreatedAt(order string) *PostTypeContext {
	v.Tx.Order(fmt.Sprintf("created_at %s", order))
	return v
}

func (v *PostTypeContext) GetViaAlias(alias string) (models.Feed, error) {
	var item models.Feed
	table := viper.GetString("database.prefix") + v.TableName
	userTable := viper.GetString("database.prefix") + "accounts"
	if err := v.Tx.
		Table(table).
		Select("*, ? as model_type", v.ColumnName).
		Joins(fmt.Sprintf("INNER JOIN %s AS author ON author_id = author.id", userTable)).
		Where("alias = ?", alias).
		First(&item).Error; err != nil {
		return item, err
	}

	var attachments []models.Attachment
	if err := database.C.
		Model(&models.Attachment{}).
		Where(v.ColumnName+"_id = ?", item.ID).
		Scan(&attachments).Error; err != nil {
		return item, err
	} else {
		item.Attachments = attachments
	}

	return item, nil
}

func (v *PostTypeContext) Get(id uint, noComments ...bool) (models.Feed, error) {
	var item models.Feed
	table := viper.GetString("database.prefix") + v.TableName
	userTable := viper.GetString("database.prefix") + "accounts"
	if err := v.Tx.
		Table(table).
		Select("*, ? as model_type", v.ColumnName).
		Joins(fmt.Sprintf("INNER JOIN %s AS author ON author_id = author.id", userTable)).
		Where("id = ?", id).First(&item).Error; err != nil {
		return item, err
	}

	var attachments []models.Attachment
	if err := database.C.
		Model(&models.Attachment{}).
		Where(v.ColumnName+"_id = ?", id).
		Scan(&attachments).Error; err != nil {
		return item, err
	} else {
		item.Attachments = attachments
	}

	return item, nil
}

func (v *PostTypeContext) Count() (int64, error) {
	var count int64
	table := viper.GetString("database.prefix") + v.TableName
	if err := v.Tx.Table(table).Count(&count).Error; err != nil {
		return count, err
	}

	return count, nil
}

func (v *PostTypeContext) CountReactions(id uint) (map[string]int64, error) {
	var reactions []struct {
		Symbol string
		Count  int64
	}

	if err := database.C.Model(&models.Reaction{}).
		Select("symbol, COUNT(id) as count").
		Where(v.ColumnName+"_id = ?", id).
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

func (v *PostTypeContext) List(take int, offset int, noReact ...bool) ([]*models.Feed, error) {
	if take > 20 {
		take = 20
	}

	var items []*models.Feed
	table := viper.GetString("database.prefix") + v.TableName
	if err := v.Tx.
		Table(table).
		Select("*, ? as model_type", v.ColumnName).
		Limit(take).Offset(offset).Find(&items).Error; err != nil {
		return items, err
	}

	idx := lo.Map(items, func(item *models.Feed, index int) uint {
		return item.ID
	})

	if len(noReact) <= 0 || !noReact[0] {
		var reactions []struct {
			PostID uint
			Symbol string
			Count  int64
		}

		if err := database.C.Model(&models.Reaction{}).
			Select(v.ColumnName+"_id as post_id, symbol, COUNT(id) as count").
			Where(v.ColumnName+"_id IN (?)", idx).
			Group("post_id, symbol").
			Scan(&reactions).Error; err != nil {
			return items, err
		}

		itemMap := lo.SliceToMap(items, func(item *models.Feed) (uint, *models.Feed) {
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

	{
		var attachments []struct {
			models.Attachment

			PostID uint `json:"post_id"`
		}

		itemMap := lo.SliceToMap(items, func(item *models.Feed) (uint, *models.Feed) {
			return item.ID, item
		})

		idx := lo.Map(items, func(item *models.Feed, index int) uint {
			return item.ID
		})

		if err := database.C.
			Model(&models.Attachment{}).
			Select(v.ColumnName+"_id as post_id, *").
			Where(v.ColumnName+"_id IN (?)", idx).
			Scan(&attachments).Error; err != nil {
			return items, err
		}

		list := map[uint][]models.Attachment{}
		for _, info := range attachments {
			list[info.PostID] = append(list[info.PostID], info.Attachment)
		}

		for k, v := range list {
			if post, ok := itemMap[k]; ok {
				post.Attachments = v
			}
		}
	}

	return items, nil
}

func MapCategoriesAndTags[T models.PostInterface](item T) (T, error) {
	var err error
	categories := item.GetCategories()
	for idx, category := range categories {
		categories[idx], err = GetCategory(category.Alias)
		if err != nil {
			return item, err
		}
	}
	item.SetCategories(categories)
	tags := item.GetHashtags()
	for idx, tag := range tags {
		tags[idx], err = GetTagOrCreate(tag.Alias, tag.Name)
		if err != nil {
			return item, err
		}
	}
	item.SetHashtags(tags)
	return item, nil
}

func NewPost[T models.PostInterface](item T) (T, error) {
	item, err := MapCategoriesAndTags(item)
	if err != nil {
		return item, err
	}

	if item.GetRealm() != nil {
		if item.GetRealm().RealmType != models.RealmTypePublic {
			var member models.RealmMember
			if err := database.C.Where(&models.RealmMember{
				RealmID:   item.GetRealm().ID,
				AccountID: item.GetAuthor().ID,
			}).First(&member).Error; err != nil {
				return item, fmt.Errorf("you aren't a part of that realm")
			}
		}
	}

	if err := database.C.Save(&item).Error; err != nil {
		return item, err
	}

	if item.GetReplyTo() != nil {
		go func() {
			var op models.Moment
			if err := database.C.Where("id = ?", item.GetReplyTo()).Preload("Author").First(&op).Error; err == nil {
				if op.Author.ID != item.GetAuthor().ID {
					postUrl := fmt.Sprintf("https://%s/posts/%d", viper.GetString("domain"), item.GetID())
					err := NotifyAccount(
						op.Author,
						fmt.Sprintf("%s replied you", item.GetAuthor().Name),
						fmt.Sprintf("%s replied your post. Check it out!", item.GetAuthor().Name),
						&proto.NotifyLink{Label: "Related post", Url: postUrl},
					)
					if err != nil {
						log.Error().Err(err).Msg("An error occurred when notifying user...")
					}
				}
			}
		}()
	}

	var subscribers []models.AccountMembership
	if err := database.C.Where(&models.AccountMembership{
		FollowingID: item.GetAuthor().ID,
	}).Preload("Follower").Find(&subscribers).Error; err == nil && len(subscribers) > 0 {
		go func() {
			accounts := lo.Map(subscribers, func(item models.AccountMembership, index int) models.Account {
				return item.Follower
			})

			for _, account := range accounts {
				postUrl := fmt.Sprintf("https://%s/posts/%d", viper.GetString("domain"), item.GetID())
				err := NotifyAccount(
					account,
					fmt.Sprintf("%s just posted a post", item.GetAuthor().Name),
					"Account you followed post a brand new post. Check it out!",
					&proto.NotifyLink{Label: "Related post", Url: postUrl},
				)
				if err != nil {
					log.Error().Err(err).Msg("An error occurred when notifying user...")
				}
			}
		}()
	}

	return item, nil
}

func EditPost[T models.PostInterface](item T) (T, error) {
	item, err := MapCategoriesAndTags(item)
	if err != nil {
		return item, err
	}

	err = database.C.Save(&item).Error

	return item, err
}

func DeletePost[T models.PostInterface](item T) error {
	return database.C.Delete(&item).Error
}

func (v *PostTypeContext) React(reaction models.Reaction) (bool, models.Reaction, error) {
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
