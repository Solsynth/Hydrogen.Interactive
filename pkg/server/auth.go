package server

import (
	"git.solsynth.dev/hydrogen/interactive/pkg/services"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func authMiddleware(c *fiber.Ctx) error {
	var token string
	if cookie := c.Cookies(services.CookieAccessKey); len(cookie) > 0 {
		token = cookie
	}
	if header := c.Get(fiber.HeaderAuthorization); len(header) > 0 {
		tk := strings.Replace(header, "Bearer", "", 1)
		token = strings.TrimSpace(tk)
	}

	c.Locals("token", token)

	if err := authFunc(c); err != nil {
		return err
	}

	return c.Next()
}

func authFunc(c *fiber.Ctx, overrides ...string) error {
	var token string
	if len(overrides) > 0 {
		token = overrides[0]
	} else {
		if tk, ok := c.Locals("token").(string); !ok {
			return fiber.NewError(fiber.StatusUnauthorized)
		} else {
			token = tk
		}
	}

	rtk := c.Cookies(services.CookieRefreshKey)
	if user, atk, rtk, err := services.Authenticate(token, rtk); err == nil {
		if atk != token {
			services.SetJwtCookieSet(c, atk, rtk)
		}
		c.Locals("principal", user)
		return nil
	} else {
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	}
}
