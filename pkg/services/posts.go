package services

import (
	"code.smartsheep.studio/hydrogen/identity/pkg/grpc/proto"
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"errors"
	"fmt"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"strings"
	"time"
)

type PostTypeContext[T models.PostInterface] struct {
	Tx *gorm.DB

	TypeName  string
	CanReply  bool
	CanRepost bool
}

var pluralizeHelper = pluralize.NewClient()

func (v *PostTypeContext[T]) GetTableName(plural ...bool) string {
	if len(plural) <= 0 || !plural[0] {
		return strings.ToLower(v.TypeName)
	} else {
		return pluralizeHelper.Plural(strings.ToLower(v.TypeName))
	}
}

func (v *PostTypeContext[T]) Preload() *PostTypeContext[T] {
	v.Tx.Preload("Author").
		Preload("Attachments").
		Preload("Categories").
		Preload("Hashtags")

	if v.CanReply {
		v.Tx.Preload("ReplyTo")
	}
	if v.CanRepost {
		v.Tx.Preload("RepostTo")
	}

	return v
}

func (v *PostTypeContext[T]) FilterWithCategory(alias string) *PostTypeContext[T] {
	table := v.GetTableName()
	v.Tx.Joins(fmt.Sprintf("JOIN %s_categories ON %s.id = %s_categories.%s_id", table, v.GetTableName(true), table, v.GetTableName())).
		Joins(fmt.Sprintf("JOIN %s_categories ON %s_categories.id = %s_categories.category_id", table, table, table)).
		Where(table+"_categories.alias = ?", alias)
	return v
}

func (v *PostTypeContext[T]) FilterWithTag(alias string) *PostTypeContext[T] {
	table := v.GetTableName()
	v.Tx.Joins(fmt.Sprintf("JOIN %s_tags ON %s.id = %s_tags.%s_id", table, v.GetTableName(true), table, v.GetTableName())).
		Joins(fmt.Sprintf("JOIN %s_tags ON %s_tags.id = %s_tags.category_id", table, table, table)).
		Where(table+"_tags.alias = ?", alias)
	return v
}

func (v *PostTypeContext[T]) FilterPublishedAt(date time.Time) *PostTypeContext[T] {
	v.Tx.Where("published_at <= ? AND published_at IS NULL", date)
	return v
}

func (v *PostTypeContext[T]) FilterRealm(id uint) *PostTypeContext[T] {
	if id > 0 {
		v.Tx = v.Tx.Where("realm_id = ?", id)
	} else {
		v.Tx = v.Tx.Where("realm_id IS NULL")
	}
	return v
}

func (v *PostTypeContext[T]) FilterAuthor(id uint) *PostTypeContext[T] {
	v.Tx = v.Tx.Where("author_id = ?", id)
	return v
}

func (v *PostTypeContext[T]) FilterReply(condition bool) *PostTypeContext[T] {
	if condition {
		v.Tx = v.Tx.Where("reply_id IS NOT NULL")
	} else {
		v.Tx = v.Tx.Where("reply_id IS NULL")
	}
	return v
}

func (v *PostTypeContext[T]) SortCreatedAt(order string) *PostTypeContext[T] {
	v.Tx.Order(fmt.Sprintf("created_at %s", order))
	return v
}

func (v *PostTypeContext[T]) GetViaAlias(alias string) (T, error) {
	var item T
	if err := v.Preload().Tx.Where("alias = ?", alias).First(&item).Error; err != nil {
		return item, err
	}

	return item, nil
}

func (v *PostTypeContext[T]) Get(id uint) (T, error) {
	var item T
	if err := v.Preload().Tx.Where("id = ?", id).First(&item).Error; err != nil {
		return item, err
	}

	return item, nil
}

func (v *PostTypeContext[T]) Count() (int64, error) {
	var count int64
	table := viper.GetString("database.prefix") + v.GetTableName(true)
	if err := v.Tx.Table(table).Count(&count).Error; err != nil {
		return count, err
	}

	return count, nil
}

func (v *PostTypeContext[T]) List(take int, offset int, noReact ...bool) ([]T, error) {
	if take > 20 {
		take = 20
	}

	var items []T
	if err := v.Preload().Tx.Limit(take).Offset(offset).Find(&items).Error; err != nil {
		return items, err
	}

	return items, nil
}

func (v *PostTypeContext[T]) MapCategoriesAndTags(item T) (T, error) {
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

func (v *PostTypeContext[T]) New(item T) (T, error) {
	item, err := v.MapCategoriesAndTags(item)
	if err != nil {
		return item, err
	}

	if item.GetRealm() != nil {
		if !item.GetRealm().IsPublic {
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

func (v *PostTypeContext[T]) Edit(item T) (T, error) {
	item, err := v.MapCategoriesAndTags(item)
	if err != nil {
		return item, err
	}

	err = database.C.Save(&item).Error

	return item, err
}

func (v *PostTypeContext[T]) Delete(item T) error {
	return database.C.Delete(&item).Error
}

func (v *PostTypeContext[T]) React(reaction models.Reaction) (bool, models.Reaction, error) {
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
