package services

import (
	"errors"
	"fmt"
	"time"

	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func PreloadRelatedPost(tx *gorm.DB) *gorm.DB {
	return tx.
		Preload("Author").
		Preload("Attachments").
		Preload("Categories").
		Preload("Tags").
		Preload("RepostTo").
		Preload("ReplyTo").
		Preload("RepostTo.Author").
		Preload("ReplyTo.Author").
		Preload("RepostTo.Attachments").
		Preload("ReplyTo.Attachments").
		Preload("RepostTo.Categories").
		Preload("ReplyTo.Categories").
		Preload("RepostTo.Tags").
		Preload("ReplyTo.Tags")
}

func ListPost(tx *gorm.DB, take int, offset int) ([]*models.Post, error) {
	var posts []*models.Post
	if err := PreloadRelatedPost(tx).
		Limit(take).
		Offset(offset).
		Find(&posts).Error; err != nil {
		return posts, err
	}

	postIds := lo.Map(posts, func(item *models.Post, _ int) uint {
		return item.ID
	})

	var reactInfo []struct {
		PostID       uint
		LikeCount    int64
		DislikeCount int64
	}

	prefix := viper.GetString("database.prefix")
	database.C.Raw(fmt.Sprintf(`SELECT t.id                         as post_id,
       COALESCE(l.like_count, 0)    AS like_count,
       COALESCE(d.dislike_count, 0) AS dislike_count
FROM %sposts t
         LEFT JOIN (SELECT post_id, COUNT(*) AS like_count
                    FROM %spost_likes
                    GROUP BY post_id) l ON t.id = l.post_id
         LEFT JOIN (SELECT post_id, COUNT(*) AS dislike_count
                    FROM %spost_dislikes
                    GROUP BY post_id) d ON t.id = d.post_id
WHERE t.id IN (?)`, prefix, prefix, prefix), postIds).Scan(&reactInfo)

	postMap := lo.SliceToMap(posts, func(item *models.Post) (uint, *models.Post) {
		return item.ID, item
	})

	for _, info := range reactInfo {
		if post, ok := postMap[info.PostID]; ok {
			post.LikeCount = info.LikeCount
			post.DislikeCount = info.DislikeCount
		}
	}

	return posts, nil
}

func NewPost(
	user models.Account,
	realm *models.Realm,
	alias, title, content string,
	attachments []models.Attachment,
	categories []models.Category,
	tags []models.Tag,
	publishedAt *time.Time,
	replyTo, repostTo *uint,
) (models.Post, error) {
	var err error
	var post models.Post
	for idx, category := range categories {
		categories[idx], err = GetCategory(category.Alias, category.Name)
		if err != nil {
			return post, err
		}
	}
	for idx, tag := range tags {
		tags[idx], err = GetTag(tag.Alias, tag.Name)
		if err != nil {
			return post, err
		}
	}

	var realmId *uint
	if realm != nil {
		realmId = &realm.ID
	}

	if publishedAt == nil {
		publishedAt = lo.ToPtr(time.Now())
	}

	post = models.Post{
		Alias:       alias,
		Title:       title,
		Content:     content,
		Attachments: attachments,
		Tags:        tags,
		Categories:  categories,
		AuthorID:    user.ID,
		RealmID:     realmId,
		PublishedAt: *publishedAt,
		RepostID:    repostTo,
		ReplyID:     replyTo,
	}

	if err := database.C.Save(&post).Error; err != nil {
		return post, err
	}

	return post, nil
}

func EditPost(
	post models.Post,
	alias, title, content string,
	publishedAt *time.Time,
	categories []models.Category,
	tags []models.Tag,
) (models.Post, error) {
	var err error
	for idx, category := range categories {
		categories[idx], err = GetCategory(category.Alias, category.Name)
		if err != nil {
			return post, err
		}
	}
	for idx, tag := range tags {
		tags[idx], err = GetTag(tag.Alias, tag.Name)
		if err != nil {
			return post, err
		}
	}

	if publishedAt == nil {
		publishedAt = lo.ToPtr(time.Now())
	}

	post.Alias = alias
	post.Title = title
	post.Content = content
	post.PublishedAt = *publishedAt

	err = database.C.Save(&post).Error

	return post, err
}

func LikePost(user models.Account, post models.Post) (bool, error) {
	var like models.PostLike
	if err := database.C.Where(&models.PostLike{
		AccountID: user.ID,
		PostID:    post.ID,
	}).First(&like).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return true, err
		}
		like = models.PostLike{
			AccountID: user.ID,
			PostID:    post.ID,
		}
		return true, database.C.Save(&like).Error
	} else {
		return false, database.C.Delete(&like).Error
	}
}

func DislikePost(user models.Account, post models.Post) (bool, error) {
	var dislike models.PostDislike
	if err := database.C.Where(&models.PostDislike{
		AccountID: user.ID,
		PostID:    post.ID,
	}).First(&dislike).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return true, err
		}
		dislike = models.PostDislike{
			AccountID: user.ID,
			PostID:    post.ID,
		}
		return true, database.C.Save(&dislike).Error
	} else {
		return false, database.C.Delete(&dislike).Error
	}
}

func DeletePost(post models.Post) error {
	return database.C.Delete(&post).Error
}
