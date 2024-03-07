package server

import (
	"strings"
	"time"

	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"code.smartsheep.studio/hydrogen/interactive/pkg/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func contextArticle() *services.PostTypeContext {
	return &services.PostTypeContext{
		Tx:         database.C,
		TableName:  "articles",
		ColumnName: "article",
		CanReply:   false,
		CanRepost:  false,
	}
}

func createArticle(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		Alias       string              `json:"alias" form:"alias"`
		Title       string              `json:"title" form:"title" validate:"required"`
		Description string              `json:"description" form:"description"`
		Content     string              `json:"content" form:"content" validate:"required"`
		Hashtags    []models.Tag        `json:"hashtags" form:"hashtags"`
		Categories  []models.Category   `json:"categories" form:"categories"`
		Attachments []models.Attachment `json:"attachments" form:"attachments"`
		PublishedAt *time.Time          `json:"published_at" form:"published_at"`
		RealmID     *uint               `json:"realm_id" form:"realm_id"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	} else if len(data.Alias) == 0 {
		data.Alias = strings.ReplaceAll(uuid.NewString(), "-", "")
	}

	item := &models.Article{
		PostBase: models.PostBase{
			Alias:       data.Alias,
			Attachments: data.Attachments,
			PublishedAt: data.PublishedAt,
			AuthorID:    user.ID,
		},
		Hashtags:    data.Hashtags,
		Categories:  data.Categories,
		Title:       data.Title,
		Description: data.Description,
		Content:     data.Content,
		RealmID:     data.RealmID,
	}

	var realm *models.Realm
	if data.RealmID != nil {
		if err := database.C.Where(&models.Realm{
			BaseModel: models.BaseModel{ID: *data.RealmID},
		}).First(&realm).Error; err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	if item, err := services.NewPost(item); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(item)
	}
}

func editArticle(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("articleId", 0)

	var data struct {
		Alias       string              `json:"alias" form:"alias" validate:"required"`
		Title       string              `json:"title" form:"title" validate:"required"`
		Description string              `json:"description" form:"description"`
		Content     string              `json:"content" form:"content" validate:"required"`
		PublishedAt *time.Time          `json:"published_at" form:"published_at"`
		Hashtags    []models.Tag        `json:"hashtags" form:"hashtags"`
		Categories  []models.Category   `json:"categories" form:"categories"`
		Attachments []models.Attachment `json:"attachments" form:"attachments"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	var item *models.Article
	if err := database.C.Where(models.Article{
		PostBase: models.PostBase{
			BaseModel: models.BaseModel{ID: uint(id)},
			AuthorID:  user.ID,
		},
	}).First(&item).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	item.Alias = data.Alias
	item.Title = data.Title
	item.Description = data.Description
	item.Content = data.Content
	item.PublishedAt = data.PublishedAt
	item.Hashtags = data.Hashtags
	item.Categories = data.Categories
	item.Attachments = data.Attachments

	if item, err := services.EditPost(item); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(item)
	}
}

func deleteArticle(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("articleId", 0)

	var item *models.Article
	if err := database.C.Where(models.Article{
		PostBase: models.PostBase{
			BaseModel: models.BaseModel{ID: uint(id)},
			AuthorID:  user.ID,
		},
	}).First(&item).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.DeletePost(item); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
