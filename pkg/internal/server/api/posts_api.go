package api

import (
	"fmt"

	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/server/exts"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func getPost(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("postId")

	item, err := services.GetPost(services.FilterPostDraft(database.C), uint(id))
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

	if user, authenticated := c.Locals("user").(models.Account); authenticated {
		tx = services.FilterPostWithUserContext(tx, &user)
	} else {
		tx = services.FilterPostWithUserContext(tx, nil)
	}

	if realmId > 0 {
		if realm, err := services.GetRealmWithExtID(uint(realmId)); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("realm was not found: %v", err))
		} else {
			tx = services.FilterPostWithRealm(tx, realm.ID)
		}
	}

	if len(c.Query("author")) > 0 {
		var author models.Account
		if err := database.C.Where(&models.Account{Name: c.Query("author")}).First(&author).Error; err != nil {
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

	items, err := services.ListPost(tx, take, offset, "published_at DESC")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  items,
	})
}

func listPostMinimal(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)
	realmId := c.QueryInt("realmId", 0)

	tx := services.FilterPostDraft(database.C)

	if user, authenticated := c.Locals("user").(models.Account); authenticated {
		tx = services.FilterPostWithUserContext(tx, &user)
	} else {
		tx = services.FilterPostWithUserContext(tx, nil)
	}

	if realmId > 0 {
		if realm, err := services.GetRealmWithExtID(uint(realmId)); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("realm was not found: %v", err))
		} else {
			tx = services.FilterPostWithRealm(tx, realm.ID)
		}
	}

	if len(c.Query("author")) > 0 {
		var author models.Account
		if err := database.C.Where(&models.Account{Name: c.Query("author")}).First(&author).Error; err != nil {
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

	items, err := services.ListPostMinimal(tx, take, offset, "published_at DESC")
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

	items, err := services.ListPost(tx, take, offset, "created_at DESC", true)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  items,
	})
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
	if err := gap.H.EnsureGrantedPerm(c, "CreateReactions", true); err != nil {
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

	var res models.Post
	if err := database.C.Where("id = ?", c.Params("postId")).Select("id").First(&res).Error; err != nil {
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

func pinPost(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var res models.Post
	if err := database.C.Where("id = ? AND author_id = ?", c.Params("postId"), user.ID).First(&res).Error; err != nil {
		return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("unable to find post in your posts to pin: %v", err))
	}

	if status, err := services.PinPost(res); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else if status {
		return c.SendStatus(fiber.StatusOK)
	} else {
		return c.SendStatus(fiber.StatusNoContent)
	}
}
