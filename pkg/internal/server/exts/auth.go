package exts

import (
	"git.solsynth.dev/hydrogen/dealer/pkg/proto"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/services"
	"github.com/gofiber/fiber/v2"
)

func LinkAccountMiddleware(c *fiber.Ctx) error {
	if val, ok := c.Locals("p_user").(*proto.UserInfo); ok {
		if account, err := services.LinkAccount(val); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		} else {
			c.Locals("user", account)
		}
	}

	return c.Next()
}
