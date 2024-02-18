package services

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func FollowAccount(followerId, followingId uint) error {
	relationship := models.AccountMembership{
		FollowerID:  followerId,
		FollowingID: followingId,
	}
	return database.C.Create(&relationship).Error
}

func UnfollowAccount(followerId, followingId uint) error {
	return database.C.Where(models.AccountMembership{
		FollowerID:  followerId,
		FollowingID: followingId,
	}).Delete(&models.AccountMembership{}).Error
}

func GetAccountFollowed(user models.Account, target models.Account) (models.AccountMembership, bool) {
	var relationship models.AccountMembership
	err := database.C.Model(&models.AccountMembership{}).
		Where(&models.AccountMembership{FollowerID: user.ID, FollowingID: target.ID}).
		First(&relationship).
		Error
	return relationship, err == nil
}

func NotifyAccount(user models.Account, subject, content string, links ...fiber.Map) error {
	agent := fiber.Post(viper.GetString("identity.endpoint") + "/api/dev/notify")
	agent.JSON(fiber.Map{
		"client_id":     viper.GetString("identity.client_id"),
		"client_secret": viper.GetString("identity.client_secret"),
		"subject":       subject,
		"content":       content,
		"links":         links,
		"user_id":       user.ExternalID,
	})

	if status, body, errs := agent.Bytes(); len(errs) > 0 {
		return errs[0]
	} else if status != 200 {
		return fmt.Errorf(string(body))
	}

	return nil
}
