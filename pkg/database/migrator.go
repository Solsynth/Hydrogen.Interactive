package database

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"gorm.io/gorm"
)

func RunMigration(source *gorm.DB) error {
	if err := source.AutoMigrate(
		&models.Account{},
		&models.Realm{},
		&models.Category{},
		&models.Tag{},
		&models.Post{},
	); err != nil {
		return err
	}

	return nil
}
