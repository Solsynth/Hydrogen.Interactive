package server

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"

	"git.solsynth.dev/hydrogen/interactive/pkg/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/models"
	"git.solsynth.dev/hydrogen/interactive/pkg/services"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func getPost(c *fiber.Ctx) error {
	alias := c.Params("postId")

	item, err := services.GetPostWithAlias(alias)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	item.ReplyCount = services.CountPostReply(item.ID)
	item.ReactionCount = services.CountPostReactions(item.ID)
	item.ReactionList, err = services.ListPostReactions(item.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(item)
}

func listPost(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)
	realmId := c.QueryInt("realmId", 0)

	tx := database.C
	if realmId > 0 {
		tx = services.FilterWithRealm(tx, uint(realmId))
	}

	if len(c.Query("authorId")) > 0 {
		var author models.Account
		if err := database.C.Where(&models.Account{Name: c.Query("authorId")}).First(&author).Error; err != nil {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		tx = tx.Where("author_id = ?", author.ID)
	}

	if len(c.Query("category")) > 0 {
		tx = services.FilterPostWithCategory(tx, c.Query("category"))
	}
	if len(c.Query("tag")) > 0 {
		tx = services.FilterPostWithTag(tx, c.Query("tag"))
	}

	count, err := services.CountPost(tx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	items, err := services.ListPost(tx, take, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  items,
	})
}

func createPost(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		Alias       string              `json:"alias" form:"alias"`
		Content     string              `json:"content" form:"content" validate:"required,max=4096"`
		Tags        []models.Tag        `json:"tags" form:"tags"`
		Categories  []models.Category   `json:"categories" form:"categories"`
		Attachments []models.Attachment `json:"attachments" form:"attachments"`
		PublishedAt *time.Time          `json:"published_at" form:"published_at"`
		RealmAlias  string              `json:"realm" form:"realm"`
		ReplyTo     *uint               `json:"reply_to" form:"reply_to"`
		RepostTo    *uint               `json:"repost_to" form:"repost_to"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	} else if len(data.Alias) == 0 {
		data.Alias = strings.ReplaceAll(uuid.NewString(), "-", "")
	}

	item := models.Post{
		Alias:       data.Alias,
		PublishedAt: data.PublishedAt,
		AuthorID:    user.ID,
		Tags:        data.Tags,
		Categories:  data.Categories,
		Attachments: data.Attachments,
		Content:     data.Content,
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

	if len(data.RealmAlias) > 0 {
		if realm, err := services.GetRealmWithAlias(data.RealmAlias); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		} else if _, err := services.GetRealmMember(realm.ExternalID, user.ExternalID); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("you aren't a part of related realm: %v", err))
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

func editPost(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("postId", 0)

	var data struct {
		Alias       string              `json:"alias" form:"alias" validate:"required"`
		Content     string              `json:"content" form:"content" validate:"required,max=1024"`
		PublishedAt *time.Time          `json:"published_at" form:"published_at"`
		Tags        []models.Tag        `json:"tags" form:"tags"`
		Categories  []models.Category   `json:"categories" form:"categories"`
		Attachments []models.Attachment `json:"attachments" form:"attachments"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	var item models.Post
	if err := database.C.Where(models.Post{
		BaseModel: models.BaseModel{ID: uint(id)},
		AuthorID:  user.ID,
	}).First(&item).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	item.Alias = data.Alias
	item.Content = data.Content
	item.PublishedAt = data.PublishedAt
	item.Tags = data.Tags
	item.Categories = data.Categories
	item.Attachments = data.Attachments

	if item, err := services.EditPost(item); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(item)
	}
}

func deletePost(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("postId", 0)

	var item models.Post
	if err := database.C.Where(models.Post{
		BaseModel: models.BaseModel{ID: uint(id)},
		AuthorID:  user.ID,
	}).First(&item).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.DeletePost(item); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func reactPost(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data struct {
		Symbol   string                  `json:"symbol" form:"symbol" validate:"required"`
		Attitude models.ReactionAttitude `json:"attitude" form:"attitude"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	reaction := models.Reaction{
		Symbol:    data.Symbol,
		Attitude:  data.Attitude,
		AccountID: user.ID,
	}

	alias := c.Params("postId")

	var res models.Post

	if err := database.C.Where("id = ?", alias).Select("id").First(&res).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to find post to react: %v", err))
	} else {
		reaction.PostID = &res.ID
	}

	if positive, reaction, err := services.ReactPost(reaction); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.Status(lo.Ternary(positive, fiber.StatusCreated, fiber.StatusNoContent)).JSON(reaction)
	}
}
