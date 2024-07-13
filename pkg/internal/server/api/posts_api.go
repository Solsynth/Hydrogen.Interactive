package api

import (
	"fmt"
	"strings"
	"time"

	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/server/exts"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
)

func getPost(c *fiber.Ctx) error {
	alias := c.Params("post")

	item, err := services.GetPostWithAlias(services.FilterPostDraft(database.C), alias)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	item.Metric = models.PostMetric{
		ReplyCount:    services.CountPostReply(item.ID),
		ReactionCount: services.CountPostReactions(item.ID),
	}
	item.Metric.ReactionList, err = services.ListResourceReactions(database.C.Where("post_id = ?", item.ID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(item)
}

func listPost(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)
	realmId := c.QueryInt("realmId", 0)

	tx := services.FilterPostDraft(database.C)
	if realmId > 0 {
		if realm, err := services.GetRealmWithExtID(uint(realmId)); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("realm was not found: %v", err))
		} else {
			tx = services.FilterPostWithRealm(tx, realm.ID)
		}
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

	countTx := tx
	count, err := services.CountPost(countTx)
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

func listDraftPost(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	tx := services.FilterPostWithAuthorDraft(database.C, user.ID)

	count, err := services.CountPost(tx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	items, err := services.ListPost(tx, take, offset, true)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  items,
	})
}

func createPost(c *fiber.Ctx) error {
	if err := gap.H.EnsureGrantedPerm(c, "CreatePosts", true); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Alias       string            `json:"alias"`
		Content     string            `json:"content" validate:"required,max=4096"`
		Tags        []models.Tag      `json:"tags"`
		Categories  []models.Category `json:"categories"`
		Attachments []uint            `json:"attachments"`
		IsDraft     bool              `json:"is_draft"`
		PublishedAt *time.Time        `json:"published_at"`
		RealmAlias  string            `json:"realm"`
		ReplyTo     *uint             `json:"reply_to"`
		RepostTo    *uint             `json:"repost_to"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	} else if len(data.Alias) == 0 {
		data.Alias = strings.ReplaceAll(uuid.NewString(), "-", "")
	}

	for _, attachment := range data.Attachments {
		if !services.CheckAttachmentByIDExists(attachment, "i.attachment") {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("attachment %d not found", attachment))
		}
	}

	item := models.Post{
		Alias:       data.Alias,
		Content:     &data.Content,
		Tags:        data.Tags,
		Categories:  data.Categories,
		Attachments: data.Attachments,
		IsDraft:     data.IsDraft,
		PublishedAt: data.PublishedAt,
		AuthorID:    user.ID,
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
	id, _ := c.ParamsInt("postId", 0)
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Alias       string            `json:"alias"`
		Content     string            `json:"content" validate:"required,max=4096"`
		IsDraft     bool              `json:"is_draft"`
		PublishedAt *time.Time        `json:"published_at"`
		Tags        []models.Tag      `json:"tags"`
		Categories  []models.Category `json:"categories"`
		Attachments []uint            `json:"attachments"`
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

	for _, attachment := range data.Attachments {
		if !services.CheckAttachmentByIDExists(attachment, "i.attachment") {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("attachment %d not found", attachment))
		}
	}

	item.Content = &data.Content
	item.Alias = data.Alias
	item.IsDraft = data.IsDraft
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
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
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
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Symbol   string                  `json:"symbol"`
		Attitude models.ReactionAttitude `json:"attitude"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	reaction := models.Reaction{
		Symbol:    data.Symbol,
		Attitude:  data.Attitude,
		AccountID: user.ID,
	}

	alias := c.Params("post")

	var res models.Post
	if err := database.C.Where("alias = ?", alias).Select("id").First(&res).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to find post to react: %v", err))
	} else {
		reaction.PostID = &res.ID
	}

	if positive, reaction, err := services.ReactPost(user, reaction); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.Status(lo.Ternary(positive, fiber.StatusCreated, fiber.StatusNoContent)).JSON(reaction)
	}
}
