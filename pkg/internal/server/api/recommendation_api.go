package api

import (
	"fmt"

	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func listRecommendationFeatured(c *fiber.Ctx) error {
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

	countTx := tx
	count, err := services.CountPost(countTx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	items, err := services.ListPost(tx, take, offset, "published_at DESC, (COALESCE(total_upvote, 0) - COALESCE(total_downvote, 0)) DESC")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  items,
	})
}

func listRecommendationNews(c *fiber.Ctx) error {
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

	fmt.Println(database.C.ToSQL(func(tx *gorm.DB) *gorm.DB {
		tx = services.FilterPostDraft(tx)

		if user, authenticated := c.Locals("user").(models.Account); authenticated {
			tx = services.FilterPostWithUserContext(tx, &user)
		} else {
			tx = services.FilterPostWithUserContext(tx, nil)
		}

		services.CountPost(tx)
		return tx
	}))

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

	countTx := tx
	count, err := services.CountPost(countTx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	items, err := services.ListPost(tx, take, offset, "RANDOM()")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  items,
	})
}
