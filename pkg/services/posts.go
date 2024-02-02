package services

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"errors"
	"gorm.io/gorm"
)

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
