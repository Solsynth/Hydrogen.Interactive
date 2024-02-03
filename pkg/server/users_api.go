package server

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"code.smartsheep.studio/hydrogen/interactive/pkg/services"
	"github.com/gofiber/fiber/v2"
)

func getUserinfo(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	var data models.Account
	if err := database.C.
		Where(&models.Account{BaseModel: models.BaseModel{ID: user.ID}}).
		First(&data).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(data)
}

func doFollowAccount(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("accountId", 0)

	var account models.Account
	if err := database.C.Where(&models.Account{
		BaseModel: models.BaseModel{ID: uint(id)},
	}).First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if _, ok := services.GetAccountFollowed(user, account); ok {
		if err := services.UnfollowAccount(user.ID, account.ID); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return c.SendStatus(fiber.StatusNoContent)
	} else {
		if err := services.FollowAccount(user.ID, account.ID); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		return c.SendStatus(fiber.StatusCreated)
	}
}
