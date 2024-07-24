package api

import (
	"fmt"

	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func listRecommendationDefault(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)
	realmId := c.QueryInt("realmId", 0)
	maxDownVote := c.QueryInt("maxDownVote", 3)

	tx := database.C.Joins("Author").Where("\"Author\".total_downvote <= ?", maxDownVote)
	tx = services.FilterPostDraft(tx)
	if realmId > 0 {
		if realm, err := services.GetRealmWithExtID(uint(realmId)); err != nil {
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

	items, err := services.ListPost(tx, take, offset, "published_at DESC")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  items,
	})
}

func listRecommendationShuffle(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)
	realmId := c.QueryInt("realmId", 0)
	maxDownVote := c.QueryInt("maxDownVote", 3)

	tx := database.C.Joins("Author").Where("\"Author\".total_downvote <= ?", maxDownVote)
	tx = services.FilterPostDraft(tx)
	if realmId > 0 {
		if realm, err := services.GetRealmWithExtID(uint(realmId)); err != nil {
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

	items, err := services.ListPost(tx, take, offset, gorm.Expr("RANDOM()"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  items,
	})
}
