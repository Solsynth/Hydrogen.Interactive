package api

import (
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

func listRecommendationNews(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	tx := database.C

	var err error
	if tx, err = universalPostFilter(c, tx); err != nil {
		return err
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

	items, err := services.ListPost(tx, take, offset, order)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if c.QueryBool("truncate", true) {
		for _, item := range items {
			if item != nil {
				item = lo.ToPtr(services.TruncatePostContent(*item))
			}
		}
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  items,
	})
}

func listRecommendationFriends(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	tx := database.C

	var err error
	if tx, err = universalPostFilter(c, tx); err != nil {
		return err
	}

	friends, _ := services.ListAccountFriends(user)
	friendList := lo.Map(friends, func(item models.Account, index int) uint {
		return item.ID
	})

	tx = tx.Where("author_id IN ?", friendList)

	countTx := tx
	count, err := services.CountPost(countTx)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	order := "published_at DESC"
	if c.QueryBool("featured", false) {
		order = "published_at DESC, (COALESCE(total_upvote, 0) - COALESCE(total_downvote, 0)) DESC"
	}

	items, err := services.ListPost(tx, take, offset, order)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if c.QueryBool("truncate", true) {
		for _, item := range items {
			if item != nil {
				item = lo.ToPtr(services.TruncatePostContent(*item))
			}
		}
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  items,
	})
}

func listRecommendationShuffle(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)

	tx := database.C

	var err error
	if tx, err = universalPostFilter(c, tx); err != nil {
		return err
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

	if c.QueryBool("truncate", true) {
		for _, item := range items {
			if item != nil {
				item = lo.ToPtr(services.TruncatePostContent(*item))
			}
		}
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  items,
	})
}
