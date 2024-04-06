package server

import (
	"git.solsynth.dev/hydrogen/interactive/pkg/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/models"
	"git.solsynth.dev/hydrogen/interactive/pkg/services"
	"github.com/gofiber/fiber/v2"
)

func listRealmMembers(c *fiber.Ctx) error {
	realmId, _ := c.ParamsInt("realmId", 0)

	if members, err := services.ListRealmMember(uint(realmId)); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	} else {
		return c.JSON(members)
	}
}

func inviteRealm(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	realmId, _ := c.ParamsInt("realmId", 0)

	var data struct {
		AccountName string `json:"account_name" validate:"required"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	var realm models.Realm
	if err := database.C.Where(&models.Realm{
		BaseModel: models.BaseModel{ID: uint(realmId)},
		AccountID: user.ID,
	}).First(&realm).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	var account models.Account
	if err := database.C.Where(&models.Account{
		Name: data.AccountName,
	}).First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.InviteRealmMember(account, realm); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.SendStatus(fiber.StatusOK)
	}
}

func kickRealm(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	realmId, _ := c.ParamsInt("realmId", 0)

	var data struct {
		AccountName string `json:"account_name" validate:"required"`
	}

	if err := BindAndValidate(c, &data); err != nil {
		return err
	}

	var realm models.Realm
	if err := database.C.Where(&models.Realm{
		BaseModel: models.BaseModel{ID: uint(realmId)},
		AccountID: user.ID,
	}).First(&realm).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	var account models.Account
	if err := database.C.Where(&models.Account{
		Name: data.AccountName,
	}).First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.KickRealmMember(account, realm); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.SendStatus(fiber.StatusOK)
	}
}

func leaveRealm(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)
	realmId, _ := c.ParamsInt("realmId", 0)

	var realm models.Realm
	if err := database.C.Where(&models.Realm{
		BaseModel: models.BaseModel{ID: uint(realmId)},
	}).First(&realm).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else if user.ID == realm.AccountID {
		return fiber.NewError(fiber.StatusBadRequest, "you cannot leave your own realm")
	}

	var account models.Account
	if err := database.C.Where(&models.Account{
		BaseModel: models.BaseModel{ID: user.ID},
	}).First(&account).Error; err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.KickRealmMember(account, realm); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	} else {
		return c.SendStatus(fiber.StatusOK)
	}
}
