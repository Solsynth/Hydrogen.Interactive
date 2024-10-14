package api

import (
	"fmt"
	"strconv"
	"time"

	"git.solsynth.dev/hydrogen/dealer/pkg/hyper"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/server/exts"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	jsoniter "github.com/json-iterator/go"
	"github.com/samber/lo"
)

func createArticle(c *fiber.Ctx) error {
	if err := gap.H.EnsureGrantedPerm(c, "CreatePosts", true); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Alias          *string           `json:"alias"`
		Title          string            `json:"title" validate:"required,max=1024"`
		Description    *string           `json:"description"`
		Content        string            `json:"content" validate:"required"`
		Thumbnail      *uint             `json:"thumbnail"`
		Attachments    []string          `json:"attachments"`
		Tags           []models.Tag      `json:"tags"`
		Categories     []models.Category `json:"categories"`
		PublishedAt    *time.Time        `json:"published_at"`
		PublishedUntil *time.Time        `json:"published_until"`
		VisibleUsers   []uint            `json:"visible_users_list"`
		InvisibleUsers []uint            `json:"invisible_users_list"`
		Visibility     *int8             `json:"visibility"`
		IsDraft        bool              `json:"is_draft"`
		RealmAlias     *string           `json:"realm"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	body := models.PostArticleBody{
		Thumbnail:   data.Thumbnail,
		Title:       data.Title,
		Description: data.Description,
		Content:     data.Content,
		Attachments: data.Attachments,
	}

	var bodyMapping map[string]any
	rawBody, _ := jsoniter.Marshal(body)
	_ = jsoniter.Unmarshal(rawBody, &bodyMapping)

	item := models.Post{
		Alias:          data.Alias,
		Type:           models.PostTypeArticle,
		Body:           bodyMapping,
		Language:       services.DetectLanguage(data.Content),
		Tags:           data.Tags,
		Categories:     data.Categories,
		IsDraft:        data.IsDraft,
		PublishedAt:    data.PublishedAt,
		PublishedUntil: data.PublishedUntil,
		VisibleUsers:   data.VisibleUsers,
		InvisibleUsers: data.InvisibleUsers,
		AuthorID:       user.ID,
	}

	if item.PublishedAt == nil {
		item.PublishedAt = lo.ToPtr(time.Now())
	}

	if data.Visibility != nil {
		item.Visibility = *data.Visibility
	} else {
		item.Visibility = models.PostVisibilityAll
	}

	if data.RealmAlias != nil {
		if realm, err := services.GetRealmWithAlias(*data.RealmAlias); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else if _, err = services.GetRealmMember(realm.ID, user.ID); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to post in the realm, access denied: %v", err))
		} else {
			item.RealmID = &realm.ID
			item.Realm = &realm
		}
	}

	item, err := services.NewPost(user, item)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		_ = gap.H.RecordAuditLog(
			user.ID,
			"posts.new",
			strconv.Itoa(int(item.ID)),
			c.IP(),
			c.Get(fiber.HeaderUserAgent),
		)
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
		Alias          *string           `json:"alias"`
		Title          string            `json:"title" validate:"required,max=1024"`
		Description    *string           `json:"description"`
		Content        string            `json:"content" validate:"required"`
		Thumbnail      *uint             `json:"thumbnail"`
		Attachments    []string          `json:"attachments"`
		Tags           []models.Tag      `json:"tags"`
		Categories     []models.Category `json:"categories"`
		PublishedAt    *time.Time        `json:"published_at"`
		PublishedUntil *time.Time        `json:"published_until"`
		VisibleUsers   []uint            `json:"visible_users_list"`
		InvisibleUsers []uint            `json:"invisible_users_list"`
		Visibility     *int8             `json:"visibility"`
		IsDraft        bool              `json:"is_draft"`
		RealmAlias     *string           `json:"realm"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	var item models.Post
	if err := database.C.Where(models.Post{
		BaseModel: hyper.BaseModel{ID: uint(id)},
		AuthorID:  user.ID,
	}).First(&item).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if item.LockedAt != nil {
		return fiber.NewError(fiber.StatusForbidden, "post was locked")
	}

	if !item.IsDraft && !data.IsDraft {
		item.EditedAt = lo.ToPtr(time.Now())
	}

	if item.IsDraft && !data.IsDraft && data.PublishedAt == nil {
		item.PublishedAt = lo.ToPtr(time.Now())
	} else {
		item.PublishedAt = data.PublishedAt
	}

	body := models.PostArticleBody{
		Thumbnail:   data.Thumbnail,
		Title:       data.Title,
		Description: data.Description,
		Content:     data.Content,
		Attachments: data.Attachments,
	}

	var bodyMapping map[string]any
	rawBody, _ := jsoniter.Marshal(body)
	_ = jsoniter.Unmarshal(rawBody, &bodyMapping)

	item.Alias = data.Alias
	item.Body = bodyMapping
	item.Language = services.DetectLanguage(data.Content)
	item.Tags = data.Tags
	item.Categories = data.Categories
	item.IsDraft = data.IsDraft
	item.PublishedUntil = data.PublishedUntil
	item.VisibleUsers = data.VisibleUsers
	item.InvisibleUsers = data.InvisibleUsers
	item.Author = user

	if data.Visibility != nil {
		item.Visibility = *data.Visibility
	}

	if data.RealmAlias != nil {
		if realm, err := services.GetRealmWithAlias(*data.RealmAlias); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else if _, err = services.GetRealmMember(realm.ID, user.ID); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to post in the realm, access denied: %v", err))
		} else {
			item.RealmID = &realm.ID
			item.Realm = &realm
		}
	}

	var err error
	if item, err = services.EditPost(item); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		_ = gap.H.RecordAuditLog(
			user.ID,
			"posts.edit",
			strconv.Itoa(int(item.ID)),
			c.IP(),
			c.Get(fiber.HeaderUserAgent),
		)
	}

	return c.JSON(item)
}
