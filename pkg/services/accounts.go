package services

import (
	"context"
	"git.solsynth.dev/hydrogen/interactive/pkg/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/grpc"
	"git.solsynth.dev/hydrogen/interactive/pkg/models"
	"git.solsynth.dev/hydrogen/passport/pkg/grpc/proto"
	"github.com/rs/zerolog/log"
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

func GetAccountFriend(userId, relatedId uint, status int) (*proto.FriendshipResponse, error) {
	var user models.Account
	if err := database.C.Where("id = ?", userId).First(&user).Error; err != nil {
		return nil, err
	}
	var related models.Account
	if err := database.C.Where("id = ?", relatedId).First(&related).Error; err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	return grpc.Friendships.GetFriendship(ctx, &proto.FriendshipTwoSideLookupRequest{
		AccountId: uint64(user.ExternalID),
		RelatedId: uint64(related.ExternalID),
		Status:    uint32(status),
	})
}

func NotifyAccount(user models.Account, subject, content string, realtime bool, links ...*proto.NotifyLink) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := grpc.Notify.NotifyUser(ctx, &proto.NotifyRequest{
		ClientId:     viper.GetString("passport.client_id"),
		ClientSecret: viper.GetString("passport.client_secret"),
		Subject:      subject,
		Content:      content,
		Links:        links,
		RecipientId:  uint64(user.ExternalID),
		IsRealtime:   realtime,
		IsImportant:  false,
	})
	if err != nil {
		log.Warn().Err(err).Msg("An error occurred when notify account...")
	} else {
		log.Debug().Uint("external", user.ExternalID).Msg("Notified account.")
	}

	return err
}
