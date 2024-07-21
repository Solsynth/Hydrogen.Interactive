package api

import (
	"fmt"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/server/exts"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"time"
)

func createArticle(c *fiber.Ctx) error {
	if err := gap.H.EnsureGrantedPerm(c, "CreatePosts", true); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Title       string            `json:"title" validate:"required,max=1024"`
		Description *string           `json:"description" validate:"max=2048"`
		Content     string            `json:"content" validate:"required"`
		Attachments []uint            `json:"attachments"`
		PublishedAt *time.Time        `json:"published_at"`
		IsDraft     bool              `json:"is_draft"`
		RealmAlias  *string           `json:"realm"`
		Tags        []models.Tag      `json:"tags"`
		Categories  []models.Category `json:"categories"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	body := models.PostArticleBody{
		Title:       data.Title,
		Description: data.Description,
		Content:     data.Content,
		Attachments: data.Attachments,
	}

	var bodyMapping map[string]any
	rawBody, _ := jsoniter.Marshal(body)
	_ = jsoniter.Unmarshal(rawBody, &bodyMapping)

	item := models.Post{
		Body:        bodyMapping,
		Tags:        data.Tags,
		Categories:  data.Categories,
		IsDraft:     data.IsDraft,
		PublishedAt: data.PublishedAt,
		AuthorID:    user.ID,
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

func editArticle(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("postId", 0)
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Title       string            `json:"title" validate:"required,max=1024"`
		Description *string           `json:"description" validate:"max=2048"`
		Content     string            `json:"content" validate:"required"`
		Attachments []uint            `json:"attachments"`
		IsDraft     bool              `json:"is_draft"`
		PublishedAt *time.Time        `json:"published_at"`
		Tags        []models.Tag      `json:"tags"`
		Categories  []models.Category `json:"categories"`
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

	body := models.PostArticleBody{
		Title:       data.Title,
		Content:     data.Content,
		Attachments: data.Attachments,
	}

	var bodyMapping map[string]any
	rawBody, _ := jsoniter.Marshal(body)
	_ = jsoniter.Unmarshal(rawBody, &bodyMapping)

	item.Body = bodyMapping
	item.IsDraft = data.IsDraft
	item.PublishedAt = data.PublishedAt
	item.Tags = data.Tags
	item.Categories = data.Categories

	if item, err := services.EditPost(item); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(item)
	}
}
