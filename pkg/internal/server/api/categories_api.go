package api

import (
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/server/exts"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func getCategory(c *fiber.Ctx) error {
	alias := c.Params("category")

	category, err := services.GetCategory(alias)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(category)
}

func listCategories(c *fiber.Ctx) error {
	categories, err := services.ListCategory()
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	return c.JSON(categories)
}

func newCategory(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	if user.PowerLevel <= 55 {
		return fiber.NewError(fiber.StatusForbidden, "require power level 55 to create categories")
	}

	var data struct {
		Alias       string `json:"alias" validate:"required"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	category, err := services.NewCategory(data.Alias, data.Name, data.Description)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(category)
}

func editCategory(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	if user.PowerLevel <= 55 {
		return fiber.NewError(fiber.StatusForbidden, "require power level 55 to edit categories")
	}

	id, _ := c.ParamsInt("categoryId", 0)
	category, err := services.GetCategoryWithID(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	var data struct {
		Alias       string `json:"alias" validate:"required"`
		Name        string `json:"name" validate:"required"`
		Description string `json:"description"`
	}

	if err := exts.BindAndValidate(c, &data); err != nil {
		return err
	}

	category, err = services.EditCategory(category, data.Alias, data.Name, data.Description)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(category)
}

func deleteCategory(c *fiber.Ctx) error {
	if err := gap.H.EnsureAuthenticated(c); err != nil {
		return err
	}
	user := c.Locals("user").(models.Account)
	if user.PowerLevel <= 55 {
		return fiber.NewError(fiber.StatusForbidden, "require power level 55 to delete categories")
	}

	id, _ := c.ParamsInt("categoryId", 0)
	category, err := services.GetCategoryWithID(uint(id))
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	if err := services.DeleteCategory(category); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(category)
}
