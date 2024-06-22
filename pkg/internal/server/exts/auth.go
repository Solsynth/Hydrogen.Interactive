package exts

import (
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/services"
	"git.solsynth.dev/hydrogen/passport/pkg/proto"
	"github.com/gofiber/fiber/v2"
)

func LinkAccountMiddleware(c *fiber.Ctx) error {
	if val, ok := c.Locals("p_user").(*proto.Userinfo); ok {
		if account, err := services.LinkAccount(val); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		} else {
			c.Locals("user", account)
		}
	}

	return c.Next()
}
