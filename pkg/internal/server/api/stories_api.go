package api

import (
	"fmt"
	"time"

	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/server/exts"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
)

func createStory(c *fiber.Ctx) error {
	if err := gap.H.EnsureGrantedPerm(c, "CreatePosts", true); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Title          *string           `json:"title"`
		Content        string            `json:"content" validate:"required,max=4096"`
		Location       *string           `json:"location"`
		Attachments    []uint            `json:"attachments"`
		Tags           []models.Tag      `json:"tags"`
		Categories     []models.Category `json:"categories"`
		PublishedAt    *time.Time        `json:"published_at"`
		PublishedUntil *time.Time        `json:"published_until"`
		IsDraft        bool              `json:"is_draft"`
		RealmAlias     *string           `json:"realm"`
		ReplyTo        *uint             `json:"reply_to"`
		RepostTo       *uint             `json:"repost_to"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	body := models.PostStoryBody{
		Title:       data.Title,
		Content:     data.Content,
		Location:    data.Location,
		Attachments: data.Attachments,
	}

	var bodyMapping map[string]any
	rawBody, _ := jsoniter.Marshal(body)
	_ = jsoniter.Unmarshal(rawBody, &bodyMapping)

	item := models.Post{
		Type:           models.PostTypeStory,
		Body:           bodyMapping,
		Language:       services.DetectLanguage(data.Content),
		Tags:           data.Tags,
		Categories:     data.Categories,
		PublishedAt:    data.PublishedAt,
		PublishedUntil: data.PublishedUntil,
		IsDraft:        data.IsDraft,
		AuthorID:       user.ID,
	}

	if data.ReplyTo != nil {
		var replyTo models.Post
		if err := database.C.Where("id = ?", data.ReplyTo).First(&replyTo).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("related post was not found: %v", err))
		} else {
			item.ReplyID = &replyTo.ID
		}
	}
	if data.RepostTo != nil {
		var repostTo models.Post
		if err := database.C.Where("id = ?", data.RepostTo).First(&repostTo).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("related post was not found: %v", err))
		} else {
			item.RepostID = &repostTo.ID
		}
	}

	if data.RealmAlias != nil {
		if realm, err := services.GetRealmWithAlias(*data.RealmAlias); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else if _, err = services.GetRealmMember(realm.ExternalID, user.ExternalID); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to post in the realm, access denied: %v", err))
		} else {
			item.RealmID = &realm.ID
		}
	}

	item, err := services.NewPost(user, item)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(item)
}

func editStory(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("postId", 0)
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Title          *string           `json:"title"`
		Content        string            `json:"content" validate:"required,max=4096"`
		Location       *string           `json:"location"`
		Attachments    []uint            `json:"attachments"`
		Tags           []models.Tag      `json:"tags"`
		Categories     []models.Category `json:"categories"`
		PublishedAt    *time.Time        `json:"published_at"`
		PublishedUntil *time.Time        `json:"published_until"`
		IsDraft        bool              `json:"is_draft"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	var item models.Post
	if err := database.C.Where(models.Post{
		BaseModel: models.BaseModel{ID: uint(id)},
		AuthorID:  user.ID,
	}).First(&item).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if item.IsDraft && !data.IsDraft {
		item.PublishedAt = lo.ToPtr(time.Now())
	}

	body := models.PostStoryBody{
		Title:       data.Title,
		Content:     data.Content,
		Location:    data.Location,
		Attachments: data.Attachments,
	}

	var bodyMapping map[string]any
	rawBody, _ := jsoniter.Marshal(body)
	_ = jsoniter.Unmarshal(rawBody, &bodyMapping)

	item.Body = bodyMapping
	item.Language = services.DetectLanguage(data.Content)
	item.Tags = data.Tags
	item.Categories = data.Categories
	item.PublishedAt = data.PublishedAt
	item.PublishedUntil = data.PublishedUntil
	item.IsDraft = data.IsDraft

	if item, err := services.EditPost(item); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(item)
	}
}
