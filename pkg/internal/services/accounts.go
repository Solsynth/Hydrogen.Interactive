package services

import (
	"context"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"git.solsynth.dev/hydrogen/passport/pkg/proto"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"time"
)

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

	pc, err := gap.H.DiscoverServiceGRPC("Hydrogen.Passport")
	if err != nil {
		return nil, err
	}
	return proto.NewFriendshipsClient(pc).GetFriendship(ctx, &proto.FriendshipTwoSideLookupRequest{
		AccountId: uint64(user.ExternalID),
		RelatedId: uint64(related.ExternalID),
		Status:    uint32(status),
	})
}

func NotifyPosterAccount(user models.Account, subject, content string, links ...*proto.NotifyLink) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pc, err := gap.H.DiscoverServiceGRPC("Hydrogen.Passport")
	if err != nil {
		return err
	}
	_, err = proto.NewNotifyClient(pc).NotifyUser(ctx, &proto.NotifyRequest{
		ClientId:     viper.GetString("passport.client_id"),
		ClientSecret: viper.GetString("passport.client_secret"),
		Type:         "interactive.feedback",
		Subject:      subject,
		Content:      content,
		Links:        links,
		RecipientId:  uint64(user.ExternalID),
		IsRealtime:   false,
		IsForcePush:  true,
	})
	if err != nil {
		log.Warn().Err(err).Msg("An error occurred when notify account...")
	} else {
		log.Debug().Uint("eid", user.ExternalID).Msg("Notified account.")
	}

	return err
}
