package server

import (
	"fmt"
	"strings"
	"time"

	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"code.smartsheep.studio/hydrogen/interactive/pkg/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func contextComment() *services.PostTypeContext {
	return &services.PostTypeContext{
		Tx:         database.C,
		TableName:  "comments",
		ColumnName: "comment",
		CanReply:   false,
		CanRepost:  true,
	}
}

func createComment(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		Alias       string              `json:"alias"`
		Content     string              `json:"content" validate:"required"`
		Hashtags    []models.Tag        `json:"hashtags"`
		Categories  []models.Category   `json:"categories"`
		Attachments []models.Attachment `json:"attachments"`
		PublishedAt *time.Time          `json:"published_at"`
		ArticleID   *uint               `json:"article_id"`
		MomentID    *uint               `json:"moment_id"`
		ReplyTo     uint                `json:"reply_to"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	} else if len(data.Alias) == 0 {
		data.Alias = strings.ReplaceAll(uuid.NewString(), "-", "")
	}

	item := &models.Comment{
		PostBase: models.PostBase{
			Alias:       data.Alias,
			Attachments: data.Attachments,
			PublishedAt: data.PublishedAt,
			AuthorID:    user.ID,
		},
		Hashtags:   data.Hashtags,
		Categories: data.Categories,
		Content:    data.Content,
	}

	if data.ArticleID == nil && data.MomentID == nil {
		return fiber.NewError(fiber.StatusBadRequest, "comment must belongs to a resource")
	}
	if data.ArticleID != nil {
		var article models.Article
		if err := database.C.Where("id = ?", data.ArticleID).First(&article).Error; err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("belongs to resource was not found: %v", err))
		}
	}

	var relatedCount int64
	if data.ReplyTo > 0 {
		if err := database.C.Where("id = ?", data.ReplyTo).
			Model(&models.Comment{}).Count(&relatedCount).Error; err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else if relatedCount <= 0 {
			return fiber.NewError(fiber.StatusNotFound, "related post was not found")
		} else {
			item.ReplyID = &data.ReplyTo
		}
	}

	item, err := services.NewPost(item)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(item)
}

func editComment(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("commentId", 0)

	var data struct {
		Alias       string              `json:"alias" validate:"required"`
		Content     string              `json:"content" validate:"required"`
		PublishedAt *time.Time          `json:"published_at"`
		Hashtags    []models.Tag        `json:"hashtags"`
		Categories  []models.Category   `json:"categories"`
		Attachments []models.Attachment `json:"attachments"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	var item *models.Comment
	if err := database.C.Where(models.Comment{
		PostBase: models.PostBase{
			BaseModel: models.BaseModel{ID: uint(id)},
			AuthorID:  user.ID,
		},
	}).First(&item).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	item.Alias = data.Alias
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

func deleteComment(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("commentId", 0)

	var item *models.Comment
	if err := database.C.Where(models.Comment{
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
