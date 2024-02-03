package server

import (
	"strings"

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
	authorId := c.QueryInt("authorId", 0)

	tx := database.C.Where(&models.Post{RealmID: nil}).Order("created_at desc")

	if authorId > 0 {
		tx = tx.Where(&models.Post{AuthorID: uint(authorId)})
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

func createPost(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		Alias      string            `json:"alias"`
		Title      string            `json:"title"`
		Content    string            `json:"content" validate:"required"`
		Tags       []models.Tag      `json:"tags"`
		Categories []models.Category `json:"categories"`
		RepostTo   uint              `json:"repost_to"`
		ReplyTo    uint              `json:"reply_to"`
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

	post, err := services.NewPost(user, data.Alias, data.Title, data.Content, data.Categories, data.Tags, replyTo, repostTo)
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
