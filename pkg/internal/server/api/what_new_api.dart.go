package api

import (
	"fmt"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func getWhatsNew(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	pivot := c.QueryInt("pivot", 0)
	if pivot < 0 {
		return fiber.NewError(fiber.StatusBadRequest, "pivot must be greater than zero")
	}

	realm := c.Query("realm")

	tx := services.FilterPostDraft(database.C)
	tx = services.FilterPostWithUserContext(tx, &user)

	friends, _ := services.ListAccountFriends(user)
	friendList := lo.Map(friends, func(item models.Account, index int) uint {
		return item.ID
	})

	tx = tx.Where("id > ?", pivot)
	tx = tx.Where("author_id IN ?", friendList)

	if len(realm) > 0 {
		if realm, err := services.GetRealmWithAlias(realm); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("realm was not found: %v", err))
		} else {
			tx = services.FilterPostWithRealm(tx, realm.ID)
		}
	}

	countTx := tx
	count, err := services.CountPost(countTx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	order := "published_at DESC"
	if c.QueryBool("featured", false) {
		order = "published_at DESC, (COALESCE(total_upvote, 0) - COALESCE(total_downvote, 0)) DESC"
	}

	items, err := services.ListPost(tx, 10, 0, order)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  items,
	})
}
