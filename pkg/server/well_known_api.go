package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func getMetadata(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"name":   viper.GetString("name"),
		"domain": viper.GetString("domain"),
		"components": fiber.Map{
			"passport": viper.GetString("passport.endpoint"),
		},
	})
}
