package services

import (
	"code.smartsheep.studio/hydrogen/identity/pkg/grpc/proto"
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/grpc"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"context"
	"github.com/spf13/viper"
	"time"
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

func NotifyAccount(user models.Account, subject, content string, links ...*proto.NotifyLink) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := grpc.Notify.NotifyUser(ctx, &proto.NotifyRequest{
		ClientId:     viper.GetString("identity.client_id"),
		ClientSecret: viper.GetString("identity.client_secret"),
		Subject:      subject,
		Content:      content,
		Links:        links,
		RecipientId:  uint64(user.ID),
		IsImportant:  false,
	})

	return err
}
