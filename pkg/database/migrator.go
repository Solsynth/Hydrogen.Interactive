package database

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"gorm.io/gorm"
)

func RunMigration(source *gorm.DB) error {
	if err := source.AutoMigrate(
		&models.Account{},
	); err != nil {
		return err
	}

	return nil
}
