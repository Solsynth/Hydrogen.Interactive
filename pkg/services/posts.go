package services

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"errors"
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func ListPost(take int, offset int) ([]*models.Post, error) {
	var posts []*models.Post
	if err := database.C.
		Where(&models.Post{RealmID: nil}).
		Limit(take).
		Offset(offset).
		Order("created_at desc").
		Preload("Author").
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
	alias, title, content string,
	categories []models.Category,
	tags []models.Tag,
) (models.Post, error) {
	return NewPostWithRealm(user, nil, alias, title, content, categories, tags)
}

func NewPostWithRealm(
	user models.Account,
	realm *models.Realm,
	alias, title, content string,
	categories []models.Category,
	tags []models.Tag,
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

	post = models.Post{
		Alias:      alias,
		Title:      title,
		Content:    content,
		Tags:       tags,
		Categories: categories,
		AuthorID:   user.ID,
		RealmID:    realmId,
	}

	if err := database.C.Save(&post).Error; err != nil {
		return post, err
	}

	return post, nil
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
