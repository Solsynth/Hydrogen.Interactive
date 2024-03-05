package server

import (
	"time"

	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"code.smartsheep.studio/hydrogen/interactive/pkg/services"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

var postContextKey = "ptx"

func useDynamicContext(c *fiber.Ctx) error {
	postType := c.Params("postType")
	switch postType {
	case "articles":
		c.Locals(postContextKey, contextArticle())
	case "moments":
		c.Locals(postContextKey, contextMoment())
	case "comments":
		c.Locals(postContextKey, contextComment())
	default:
		return fiber.NewError(fiber.StatusBadRequest, "invalid dataset")
	}

	return c.Next()
}

func getPost(c *fiber.Ctx) error {
	alias := c.Params("postId")

	mx := c.Locals(postContextKey).(*services.PostTypeContext)

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

func listPost(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)
	realmId := c.QueryInt("realmId", 0)

	mx := c.Locals(postContextKey).(*services.PostTypeContext).
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

func reactPost(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("articleId", 0)

	var data struct {
		Symbol   string                  `json:"symbol" validate:"required"`
		Attitude models.ReactionAttitude `json:"attitude" validate:"required"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	mx := c.Locals(postContextKey).(*services.PostTypeContext)

	item, err := mx.Get(uint(id), true)
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
