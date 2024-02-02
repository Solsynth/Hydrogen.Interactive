package server

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"code.smartsheep.studio/hydrogen/interactive/pkg/security"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/keyauth"
	"strconv"
)

var auth = keyauth.New(keyauth.Config{
	KeyLookup:  "header:Authorization",
	AuthScheme: "Bearer",
	Validator: func(c *fiber.Ctx, token string) (bool, error) {
		claims, err := security.DecodeJwt(token)
		if err != nil {
			return false, err
		}

		id, _ := strconv.Atoi(claims.Subject)

		var user models.Account
		if err := database.C.Where(&models.Account{
			BaseModel: models.BaseModel{ID: uint(id)},
		}).First(&user).Error; err != nil {
			return false, err
		}

		c.Locals("principal", user)

		return true, nil
	},
	ContextKey: "token",
})
