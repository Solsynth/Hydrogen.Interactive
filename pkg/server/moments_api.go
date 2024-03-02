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

func getMomentContext() *services.PostTypeContext[models.Moment] {
	return &services.PostTypeContext[models.Moment]{
		Tx:        database.C,
		TypeName:  "Moment",
		CanReply:  false,
		CanRepost: true,
	}
}

func getMoment(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("momentId", 0)

	mx := getMomentContext().FilterPublishedAt(time.Now())

	item, err := mx.Get(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(item)
}

func listMoment(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)
	realmId := c.QueryInt("realmId", 0)

	mx := getMomentContext().
		FilterPublishedAt(time.Now()).
		FilterRealm(uint(realmId)).
		SortCreatedAt("desc")

	var author models.Account
	if len(c.Query("authorId")) > 0 {
		if err := database.C.Where(&models.Account{Name: c.Query("authorId")}).First(&author).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		mx = mx.FilterAuthor(author.ID)
	}

	if len(c.Query("category")) > 0 {
		mx = mx.FilterWithCategory(c.Query("category"))
	}
	if len(c.Query("tag")) > 0 {
		mx = mx.FilterWithTag(c.Query("tag"))
	}

	if !c.QueryBool("reply", true) {
		mx = mx.FilterReply(true)
	}

	count, err := mx.Count()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	items, err := mx.List(take, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  items,
	})
}

func createMoment(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		Alias       string              `json:"alias"`
		Title       string              `json:"title"`
		Content     string              `json:"content" validate:"required"`
		Hashtags    []models.Tag        `json:"hashtags"`
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

	mx := getMomentContext()

	item := models.Moment{
		PostBase: models.PostBase{
			Alias:       data.Alias,
			Attachments: data.Attachments,
			PublishedAt: data.PublishedAt,
			AuthorID:    user.ID,
		},
		Hashtags:   data.Hashtags,
		Categories: data.Categories,
		Content:    data.Content,
		RealmID:    data.RealmID,
	}

	var relatedCount int64
	if data.RepostTo > 0 {
		if err := database.C.Where("id = ?", data.RepostTo).
			Model(&models.Moment{}).Count(&relatedCount).Error; err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else if relatedCount <= 0 {
			return fiber.NewError(fiber.StatusNotFound, "related post was not found")
		} else {
			item.RepostID = &data.RepostTo
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

	item, err := mx.New(item)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(item)
}

func editMoment(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("momentId", 0)

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

	mx := getMomentContext().FilterAuthor(user.ID)

	item, err := mx.Get(uint(id), true)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	item.Alias = data.Alias
	item.Content = data.Content
	item.PublishedAt = data.PublishedAt
	item.Hashtags = data.Hashtags
	item.Categories = data.Categories
	item.Attachments = data.Attachments

	item, err = mx.Edit(item)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(item)
}

func reactMoment(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("momentId", 0)

	mx := getMomentContext()

	item, err := mx.Get(uint(id), true)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	switch strings.ToLower(c.Params("reactType")) {
	case "like":
		if positive, err := mx.ReactLike(user, item.ID); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else {
			return c.SendStatus(lo.Ternary(positive, fiber.StatusCreated, fiber.StatusNoContent))
		}
	case "dislike":
		if positive, err := mx.ReactDislike(user, item.ID); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else {
			return c.SendStatus(lo.Ternary(positive, fiber.StatusCreated, fiber.StatusNoContent))
		}
	default:
		return fiber.NewError(fiber.StatusBadRequest, "unsupported reaction")
	}
}

func deleteMoment(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("momentId", 0)

	mx := getMomentContext().FilterAuthor(user.ID)

	item, err := mx.Get(uint(id), true)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := mx.Delete(item); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
