package database

import (
	"git.solsynth.dev/hydrogen/interactive/pkg/models"
	"gorm.io/gorm"
)

var DatabaseAutoActionRange = []any{
	&models.Account{},
	&models.Realm{},
	&models.RealmMember{},
	&models.Category{},
	&models.Tag{},
	&models.Moment{},
	&models.Article{},
	&models.Comment{},
	&models.Reaction{},
	&models.Attachment{},
}

func RunMigration(source *gorm.DB) error {
	if err := source.AutoMigrate(
		append([]any{
			&models.AccountMembership{},
		}, DatabaseAutoActionRange...)...,
	); err != nil {
		return err
	}

	return nil
}
