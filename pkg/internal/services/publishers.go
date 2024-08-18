package services

import (
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"github.com/gofiber/fiber/v2"
)

func GetPublisher(alias string) (any, error) {
	realm, err := GetRealmWithAlias(alias)
	if err == nil {
		return fiber.Map{
			"type": "realm",
			"data": realm,
		}, nil
	}

	var account models.Account
	if err = database.C.Where("name = ?", alias).First(&account).Error; err != nil {
		return nil, err
	}
	return fiber.Map{
		"type": "account",
		"data": account,
	}, nil
}
