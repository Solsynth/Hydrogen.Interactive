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
	"github.com/samber/lo"
)

func contextComment() *services.PostTypeContext[*models.Comment] {
	return &services.PostTypeContext[*models.Comment]{
		Tx:        database.C,
		TypeName:  "Comment",
		CanReply:  false,
		CanRepost: true,
	}
}

func getComment(c *fiber.Ctx) error {
	alias := c.Params("commentId")

	mx := contextComment().FilterPublishedAt(time.Now())

	item, err := mx.GetViaAlias(alias)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	item.ReactionList, err = mx.CountReactions(item.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(item)
}

func listComment(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)
	realmId := c.QueryInt("realmId", 0)

	mx := contextComment().
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

	mx := contextComment()

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

	item, err := mx.New(item)
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

	mx := contextComment().FilterAuthor(user.ID)

	item, err := mx.Get(uint(id))
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

func reactComment(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("commentId", 0)

	var data struct {
		Symbol   string                  `json:"symbol" validate:"required"`
		Attitude models.ReactionAttitude `json:"attitude" validate:"required"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	mx := contextComment()

	item, err := mx.Get(uint(id), true)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	reaction := models.Reaction{
		Symbol:    data.Symbol,
		Attitude:  data.Attitude,
		AccountID: user.ID,
		CommentID: &item.ID,
	}

	if positive, reaction, err := mx.React(reaction); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.Status(lo.Ternary(positive, fiber.StatusCreated, fiber.StatusNoContent)).JSON(reaction)
	}

}

func deleteComment(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("commentId", 0)

	mx := contextComment().FilterAuthor(user.ID)

	item, err := mx.Get(uint(id), true)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := mx.Delete(item); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
