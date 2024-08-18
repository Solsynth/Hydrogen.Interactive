package api

import (
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func getPublisher(c *fiber.Ctx) error {
	alias := c.Params("name")
	if out, err := services.GetPublisher(alias); err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	} else {
		return c.JSON(out)
	}
}
