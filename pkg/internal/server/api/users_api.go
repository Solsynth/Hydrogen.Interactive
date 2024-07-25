package api

import (
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func getUserinfo(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)

	var data models.Account
	if err := database.C.
		Where(&models.Account{BaseModel: models.BaseModel{ID: user.ID}}).
		First(&data).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(data)
}

func getOthersInfo(c *fiber.Ctx) error {
	account := c.Params("account")

	var data models.Account
	if err := database.C.
		Where(&models.Account{Name: account}).
		First(&data).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(data)
}

func listOthersPinnedPost(c *fiber.Ctx) error {
	account := c.Params("account")

	var user models.Account
	if err := database.C.
		Where(&models.Account{Name: account}).
		First(&user).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	tx := services.FilterPostDraft(database.C)
	tx = tx.Where("author_id = ?", user.ID)
	tx = tx.Where("pinned_at IS NOT NULL")

	items, err := services.ListPost(tx, 100, 0, "published_at DESC")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(items)
}
