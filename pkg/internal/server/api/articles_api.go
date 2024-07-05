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

func getArticle(c *fiber.Ctx) error {
	alias := c.Params("article")

	item, err := services.GetArticleWithAlias(services.FilterPostDraft(database.C), alias)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	item.ReactionCount = services.CountArticleReactions(item.ID)
	item.ReactionList, err = services.ListResourceReactions(database.C.Where("article_id = ?", item.ID))
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(item)
}

func listArticle(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)
	realmId := c.QueryInt("realmId", 0)

	tx := services.FilterPostDraft(database.C)
	if realmId > 0 {
		if realm, err := services.GetRealmWithExtID(uint(realmId)); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("realm was not found: %v", err))
		} else {
			tx = services.FilterArticleWithRealm(tx, realm.ID)
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
		tx = services.FilterArticleWithCategory(tx, c.Query("category"))
	}
	if len(c.Query("tag")) > 0 {
		tx = services.FilterArticleWithTag(tx, c.Query("tag"))
	}

	counTx := tx
	count, err := services.CountArticle(counTx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	items, err := services.ListArticle(tx, take, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  items,
	})
}

func listDraftArticle(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	tx := services.FilterArticleWithAuthorDraft(database.C, user.ID)

	count, err := services.CountArticle(tx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	items, err := services.ListArticle(tx, take, offset, true)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  items,
	})
}

func createArticle(c *fiber.Ctx) error {
	if err := gap.H.EnsureGrantedPerm(c, "CreateArticles", true); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Alias       string            `json:"alias"`
		Title       string            `json:"title" validate:"required"`
		Description string            `json:"description"`
		Content     string            `json:"content"`
		Tags        []models.Tag      `json:"tags"`
		Categories  []models.Category `json:"categories"`
		Attachments []uint            `json:"attachments"`
		IsDraft     bool              `json:"is_draft"`
		PublishedAt *time.Time        `json:"published_at"`
		RealmAlias  string            `json:"realm"`
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

	item := models.Article{
		Alias:       data.Alias,
		Title:       data.Title,
		Description: data.Description,
		Content:     data.Content,
		IsDraft:     data.IsDraft,
		PublishedAt: data.PublishedAt,
		AuthorID:    user.ID,
		Tags:        data.Tags,
		Categories:  data.Categories,
		Attachments: data.Attachments,
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

	item, err := services.NewArticle(user, item)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(item)
}

func editArticle(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("articleId", 0)
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data struct {
		Alias       string            `json:"alias"`
		Title       string            `json:"title"`
		Description string            `json:"description"`
		Content     string            `json:"content"`
		IsDraft     bool              `json:"is_draft"`
		PublishedAt *time.Time        `json:"published_at"`
		Tags        []models.Tag      `json:"tags"`
		Categories  []models.Category `json:"categories"`
		Attachments []uint            `json:"attachments"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	var item models.Article
	if err := database.C.Where(models.Article{
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

	item.Alias = data.Alias
	item.Title = data.Title
	item.Description = data.Description
	item.Content = data.Content
	item.IsDraft = data.IsDraft
	item.PublishedAt = data.PublishedAt
	item.Tags = data.Tags
	item.Categories = data.Categories
	item.Attachments = data.Attachments

	if item, err := services.EditArticle(item); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.JSON(item)
	}
}

func deleteArticle(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	id, _ := c.ParamsInt("articleId", 0)

	var item models.Article
	if err := database.C.Where(models.Article{
		BaseModel: models.BaseModel{ID: uint(id)},
		AuthorID:  user.ID,
	}).First(&item).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.DeleteArticle(item); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}

func reactArticle(c *fiber.Ctx) error {
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

	alias := c.Params("article")

	var res models.Article
	if err := database.C.Where("alias = ?", alias).Select("id").First(&res).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to find article to react: %v", err))
	} else {
		reaction.ArticleID = &res.ID
	}

	if positive, reaction, err := services.ReactArticle(user, reaction); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.Status(lo.Ternary(positive, fiber.StatusCreated, fiber.StatusNoContent)).JSON(reaction)
	}
}
