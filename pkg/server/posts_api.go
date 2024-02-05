package server

import (
	"strings"
	"time"

	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"code.smartsheep.studio/hydrogen/interactive/pkg/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func listPost(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)
	authorId := c.Query("authorId")

	tx := database.C.
		Where("realm_id IS NULL").
		Where("published_at <= ? OR published_at IS NULL", time.Now()).
		Order("created_at desc")

	var author models.Account
	if len(authorId) > 0 {
		if err := database.C.Where(&models.Account{Name: authorId}).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		tx = tx.Where(&models.Post{AuthorID: author.ID})
	}

	var count int64
	if err := tx.
		Model(&models.Post{}).
		Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	posts, err := services.ListPost(tx, take, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  posts,
	})
}

func getPost(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("postId", 0)
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	var post models.Post
	if err := services.PreloadRelatedPost(database.C.Where(&models.Post{
		BaseModel: models.BaseModel{ID: uint(id)},
	})).First(&post).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	tx := database.C.
		Where(&models.Post{ReplyID: &post.ID}).
		Where("published_at <= ? OR published_at IS NULL", time.Now()).
		Order("created_at desc")

	var count int64
	if err := tx.
		Model(&models.Post{}).
		Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	posts, err := services.ListPost(tx, take, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"data":    post,
		"count":   count,
		"related": posts,
	})
}

func createPost(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		Alias       string              `json:"alias"`
		Title       string              `json:"title"`
		Content     string              `json:"content" validate:"required"`
		Tags        []models.Tag        `json:"tags"`
		Categories  []models.Category   `json:"categories"`
		Attachments []models.Attachment `json:"attachments"`
		PublishedAt *time.Time          `json:"published_at"`
		RealmID     *uint               `json:"realm_id"`
		RepostTo    uint                `json:"repost_to"`
		ReplyTo     uint                `json:"reply_to"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	} else if len(data.Alias) == 0 {
		data.Alias = strings.ReplaceAll(uuid.NewString(), "-", "")
	}

	var repostTo *uint = nil
	var replyTo *uint = nil
	var relatedCount int64
	if data.RepostTo > 0 {
		if err := database.C.Where(&models.Post{
			BaseModel: models.BaseModel{ID: data.RepostTo},
		}).Model(&models.Post{}).Count(&relatedCount).Error; err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else if relatedCount <= 0 {
			return fiber.NewError(fiber.StatusNotFound, "related post was not found")
		} else {
			repostTo = &data.RepostTo
		}
	} else if data.ReplyTo > 0 {
		if err := database.C.Where(&models.Post{
			BaseModel: models.BaseModel{ID: data.ReplyTo},
		}).Model(&models.Post{}).Count(&relatedCount).Error; err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else if relatedCount <= 0 {
			return fiber.NewError(fiber.StatusNotFound, "related post was not found")
		} else {
			replyTo = &data.ReplyTo
		}
	}

	var realm *models.Realm
	if data.RealmID != nil {
		if err := database.C.Where(&models.Realm{
			BaseModel: models.BaseModel{ID: *data.RealmID},
		}).First(&realm).Error; err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	post, err := services.NewPost(
		user,
		realm,
		data.Alias,
		data.Title,
		data.Content,
		data.Attachments,
		data.Categories,
		data.Tags,
		data.PublishedAt,
		replyTo,
		repostTo,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(post)
}

func editPost(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("postId", 0)

	var data struct {
		Alias       string            `json:"alias" validate:"required"`
		Title       string            `json:"title"`
		Content     string            `json:"content" validate:"required"`
		PublishedAt *time.Time        `json:"published_at"`
		Tags        []models.Tag      `json:"tags"`
		Categories  []models.Category `json:"categories"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	var post models.Post
	if err := database.C.Where(&models.Post{
		BaseModel: models.BaseModel{ID: uint(id)},
		AuthorID:  user.ID,
	}).First(&post).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	post, err := services.EditPost(
		post,
		data.Alias,
		data.Title,
		data.Content,
		data.PublishedAt,
		data.Categories,
		data.Tags,
	)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(post)
}

func reactPost(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("postId", 0)

	var post models.Post
	if err := database.C.Where(&models.Post{
		BaseModel: models.BaseModel{ID: uint(id)},
	}).First(&post).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	switch strings.ToLower(c.Params("reactType")) {
	case "like":
		if positive, err := services.LikePost(user, post); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else {
			return c.SendStatus(lo.Ternary(positive, fiber.StatusCreated, fiber.StatusNoContent))
		}
	case "dislike":
		if positive, err := services.DislikePost(user, post); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else {
			return c.SendStatus(lo.Ternary(positive, fiber.StatusCreated, fiber.StatusNoContent))
		}
	default:
		return fiber.NewError(fiber.StatusBadRequest, "unsupported reaction")
	}
}

func deletePost(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("postId", 0)

	var post models.Post
	if err := database.C.Where(&models.Post{
		BaseModel: models.BaseModel{ID: uint(id)},
		AuthorID:  user.ID,
	}).First(&post).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.DeletePost(post); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
