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

func contextArticle() *services.PostTypeContext[models.Article] {
	return &services.PostTypeContext[models.Article]{
		Tx:        database.C,
		TypeName:  "Article",
		CanReply:  false,
		CanRepost: false,
	}
}

func getArticle(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("articleId", 0)

	mx := contextArticle().FilterPublishedAt(time.Now())

	item, err := mx.Get(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(item)
}

func listArticle(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)
	realmId := c.QueryInt("realmId", 0)

	mx := contextArticle().
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

func createArticle(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		Alias       string              `json:"alias"`
		Title       string              `json:"title" validate:"required"`
		Description string              `json:"description"`
		Content     string              `json:"content" validate:"required"`
		Hashtags    []models.Tag        `json:"hashtags"`
		Categories  []models.Category   `json:"categories"`
		Attachments []models.Attachment `json:"attachments"`
		PublishedAt *time.Time          `json:"published_at"`
		RealmID     *uint               `json:"realm_id"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	} else if len(data.Alias) == 0 {
		data.Alias = strings.ReplaceAll(uuid.NewString(), "-", "")
	}

	mx := contextArticle()

	item := models.Article{
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

	item, err := mx.New(item)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(item)
}

func editArticle(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("articleId", 0)

	var data struct {
		Alias       string              `json:"alias" validate:"required"`
		Title       string              `json:"title" validate:"required"`
		Description string              `json:"description"`
		Content     string              `json:"content" validate:"required"`
		PublishedAt *time.Time          `json:"published_at"`
		Hashtags    []models.Tag        `json:"hashtags"`
		Categories  []models.Category   `json:"categories"`
		Attachments []models.Attachment `json:"attachments"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	mx := contextArticle().FilterAuthor(user.ID)

	item, err := mx.Get(uint(id))
	if err != nil {
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

	item, err = mx.Edit(item)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(item)
}

func reactArticle(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("articleId", 0)

	var data struct {
		Symbol   string                  `json:"symbol" validate:"required"`
		Attitude models.ReactionAttitude `json:"attitude" validate:"required"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	mx := contextArticle()

	item, err := mx.Get(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	reaction := models.Reaction{
		Symbol:    data.Symbol,
		Attitude:  data.Attitude,
		AccountID: user.ID,
		ArticleID: &item.ID,
	}

	if positive, reaction, err := mx.React(reaction); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.Status(lo.Ternary(positive, fiber.StatusCreated, fiber.StatusNoContent)).JSON(reaction)
	}

}

func deleteArticle(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("articleId", 0)

	mx := contextArticle().FilterAuthor(user.ID)

	item, err := mx.Get(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := mx.Delete(item); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
