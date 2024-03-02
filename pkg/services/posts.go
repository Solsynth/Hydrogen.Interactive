package services

import (
	"code.smartsheep.studio/hydrogen/identity/pkg/grpc/proto"
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"fmt"
	pluralize "github.com/gertd/go-pluralize"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"strings"
	"time"
)

const (
	reactUnionSelect = `SELECT t.id AS post_id,
       COALESCE(l.like_count, 0)    AS like_count,
       COALESCE(d.dislike_count, 0) AS dislike_count--!COMMA!--
       --!REPLY_UNION_COLUMN!-- --!BOTH_COMMA!--
       --!REPOST_UNION_COLUMN!-- 
	FROM %s t
         LEFT JOIN (SELECT %s_id, COUNT(*) AS like_count
                    FROM %s_likes
                    GROUP BY %s_id) l ON t.id = l.%s_id
         LEFT JOIN (SELECT %s_id, COUNT(*) AS dislike_count
                    FROM %s_likes
                    GROUP BY %s_id) d ON t.id = d.%s_id
        --!REPLY_UNION_SELECT!--
		--!REPOST_UNION_SELECT!--
	WHERE t.id = ?`
	// TODO Solve for the cross table query(like articles -> comments)
	replyUnionColumn = `COALESCE(r.reply_count, 0) AS reply_count`
	replyUnionSelect = `LEFT JOIN (SELECT reply_id, COUNT(*) AS reply_count
		FROM %s
		WHERE reply_id IS NOT NULL
        GROUP BY reply_id) r ON t.id = r.reply_id`
	repostUnionColumn = `COALESCE(rp.repost_count, 0) AS repost_count`
	repostUnionSelect = `LEFT JOIN (SELECT repost_id, COUNT(*) AS repost_count
		FROM %s
		WHERE repost_id IS NOT NULL
		GROUP BY repost_id) rp ON t.id = rp.repost_id`
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
	v.Tx.Preload("Author").Preload("Attachments").Preload("Categories").Preload("Hashtags")

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

func (v *PostTypeContext[T]) BuildReactInfoSql() string {
	column := strings.ToLower(v.TypeName)
	table := viper.GetString("database.prefix") + v.GetTableName()
	pluralTable := viper.GetString("database.prefix") + v.GetTableName(true)
	sql := fmt.Sprintf(reactUnionSelect, pluralTable, column, table, column, column, column, table, column, column)

	if v.CanReply {
		sql = strings.Replace(sql, "--!REPLY_UNION_COLUMN!--", replyUnionColumn, 1)
		sql = strings.Replace(sql, "--!REPLY_UNION_SELECT!--", fmt.Sprintf(replyUnionSelect, pluralTable), 1)
	}
	if v.CanRepost {
		sql = strings.Replace(sql, "--!REPOST_UNION_COLUMN!--", repostUnionColumn, 1)
		sql = strings.Replace(sql, "--!REPOST_UNION_SELECT!--", fmt.Sprintf(repostUnionSelect, pluralTable), 1)
	}
	if v.CanReply || v.CanRepost {
		sql = strings.ReplaceAll(sql, "--!COMMA!--", ",")
	}
	if v.CanReply && v.CanRepost {
		sql = strings.ReplaceAll(sql, "--!BOTH_COMMA!--", ",")
	}

	return sql
}

func (v *PostTypeContext[T]) Get(id uint, noReact ...bool) (T, error) {
	var item T
	if err := v.Preload().Tx.Where("id = ?", id).First(&item).Error; err != nil {
		return item, err
	}

	var reactInfo models.PostReactInfo

	if len(noReact) <= 0 || !noReact[0] {
		sql := v.BuildReactInfoSql()
		database.C.Raw(sql, item.GetID()).Scan(&reactInfo)
	}

	item.SetReactInfo(reactInfo)

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

	idx := lo.Map(items, func(item T, _ int) uint {
		return item.GetID()
	})

	var reactInfo []struct {
		PostID       uint  `json:"post_id"`
		LikeCount    int64 `json:"like_count"`
		DislikeCount int64 `json:"dislike_count"`
		ReplyCount   int64 `json:"reply_count"`
		RepostCount  int64 `json:"repost_count"`
	}

	if len(noReact) <= 0 || !noReact[0] {
		sql := v.BuildReactInfoSql()
		database.C.Raw(sql, idx).Scan(&reactInfo)
	}

	itemMap := lo.SliceToMap(items, func(item T) (uint, T) {
		return item.GetID(), item
	})

	for _, info := range reactInfo {
		if item, ok := itemMap[info.PostID]; ok {
			item.SetReactInfo(info)
		}
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

func (v *PostTypeContext[T]) ReactLike(user models.Account, id uint) (bool, error) {
	var count int64
	table := viper.GetString("database.prefix") + v.GetTableName() + "_likes"
	tx := database.C.Where("account_id = ?", user.ID).Where(v.GetTableName()+"id = ?", id)
	if tx.Count(&count); count <= 0 {
		return true, database.C.Table(table).Create(map[string]any{
			"AccountID":       user.ID,
			v.TypeName + "ID": id,
		}).Error
	} else {
		column := strings.ToLower(v.TypeName)
		return false, tx.Raw(fmt.Sprintf("DELETE FROM %s WHERE account_id = ? AND %s_id = ?", table, column), user.ID, id).Error
	}
}

func (v *PostTypeContext[T]) ReactDislike(user models.Account, id uint) (bool, error) {
	var count int64
	table := viper.GetString("database.prefix") + v.GetTableName() + "_dislikes"
	tx := database.C.Where("account_id = ?", user.ID).Where(v.GetTableName()+"id = ?", id)
	if tx.Count(&count); count <= 0 {
		return true, database.C.Table(table).Create(map[string]any{
			"AccountID":       user.ID,
			v.TypeName + "ID": id,
		}).Error
	} else {
		column := strings.ToLower(v.TypeName)
		return false, tx.Raw(fmt.Sprintf("DELETE FROM %s WHERE account_id = ? AND %s_id = ?", table, column), user.ID, id).Error
	}
}
