package database

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"gorm.io/gorm"
)

func RunMigration(source *gorm.DB) error {
	if err := source.AutoMigrate(
		&models.Account{},
		&models.AccountMembership{},
		&models.Realm{},
		&models.RealmMember{},
		&models.Category{},
		&models.Tag{},
		&models.Moment{},
		&models.Article{},
		&models.Comment{},
		&models.Reaction{},
		&models.Attachment{},
	); err != nil {
		return err
	}

	return nil
}
