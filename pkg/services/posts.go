package services

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
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
