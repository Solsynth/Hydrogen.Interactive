package database

import (
	"git.solsynth.dev/hydrogen/interactive/pkg/models"
	"gorm.io/gorm"
)

var AutoMaintainRange = []any{
	&models.Account{},
	&models.Realm{},
	&models.Category{},
	&models.Tag{},
	&models.Post{},
	&models.Reaction{},
	&models.Attachment{},
}

func RunMigration(source *gorm.DB) error {
	if err := source.AutoMigrate(
		AutoMaintainRange...,
	); err != nil {
		return err
	}

	return nil
}
