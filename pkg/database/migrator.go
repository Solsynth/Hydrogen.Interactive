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
		&models.MomentLike{},
		&models.MomentDislike{},
		&models.Article{},
		&models.ArticleLike{},
		&models.ArticleDislike{},
		&models.Comment{},
		&models.CommentLike{},
		&models.CommentDislike{},
		&models.Attachment{},
	); err != nil {
		return err
	}

	return nil
}
