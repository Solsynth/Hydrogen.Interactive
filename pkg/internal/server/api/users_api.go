package api

import (
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
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
	accountId := c.Params("accountId")

	var data models.Account
	if err := database.C.
		Where(&models.Account{Name: accountId}).
		First(&data).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(data)
}
