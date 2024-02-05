package server

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"code.smartsheep.studio/hydrogen/interactive/pkg/services"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"time"
)

func listRealms(c *fiber.Ctx) error {
	user := c.Locals("principal").(models.Account)

	realms, err := services.ListRealms(user)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(realms)
}

func listPostInRealm(c *fiber.Ctx) error {
	take := c.QueryInt("take", 0)
	offset := c.QueryInt("offset", 0)
	authorId := c.QueryInt("authorId", 0)

	realmId, _ := c.ParamsInt("realmId", 0)

	tx := database.C.
		Where(&models.Post{RealmID: lo.ToPtr(uint(realmId))}).
		Where("published_at <= ? OR published_at IS NULL", time.Now()).
		Order("created_at desc")

	if authorId > 0 {
		tx = tx.Where(&models.Post{AuthorID: uint(authorId)})
	}

	var count int64
	if err := tx.
		Model(&models.Post{}).
		Count(&count).Error; err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	posts, err := services.ListPost(tx, take, offset)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{
		"count": count,
		"data":  posts,
	})
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
