package api

import (
	"fmt"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func listFeed(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)
	realmId := c.QueryInt("realmId", 0)

	tx := database.C
	if realmId > 0 {
		if realm, err := services.GetRealmWithExtID(uint(realmId)); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, fmt.Sprintf("realm was not found: %v", err))
		} else {
			tx = services.FilterWithRealm(tx, realm.ID)
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
