package server

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"code.smartsheep.studio/hydrogen/interactive/pkg/services"
	"github.com/gofiber/fiber/v2"
)

func getRealm(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("realmId", 0)

	var realm models.Realm
	if err := database.C.Where(&models.Realm{
		BaseModel: models.BaseModel{ID: uint(id)},
	}).First(&realm).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(realm)
}

func listRealm(c *fiber.Ctx) error {
	realms, err := services.ListRealm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(realms)
}

func listOwnedRealm(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	realms, err := services.ListRealmWithUser(user)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(realms)
}

func createRealm(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	if user.PowerLevel < 10 {
		return fiber.NewError(fiber.StatusForbidden, "require power level 10 to create realm")
	}

	var data struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	realm, err := services.NewRealm(user, data.Name, data.Description)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(realm)
}

func editRealm(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("realmId", 0)

	var data struct {
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	var realm models.Realm
	if err := database.C.Where(&models.Realm{
		BaseModel: models.BaseModel{ID: uint(id)},
		AccountID: user.ID,
	}).First(&realm).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	realm, err := services.EditRealm(realm, data.Name, data.Description)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(realm)
}

func deleteRealm(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	id, _ := c.ParamsInt("realmId", 0)

	var realm models.Realm
	if err := database.C.Where(&models.Realm{
		BaseModel: models.BaseModel{ID: uint(id)},
		AccountID: user.ID,
	}).First(&realm).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.DeleteRealm(realm); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
