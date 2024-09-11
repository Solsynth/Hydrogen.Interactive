package services

import (
	"context"
	"fmt"
	"git.solsynth.dev/hydrogen/dealer/pkg/hyper"
	"git.solsynth.dev/hydrogen/dealer/pkg/proto"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
	"time"
)

func ListAccountFriends(user models.Account) ([]models.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pc, err := gap.H.GetServiceGrpcConn(hyper.ServiceTypeAuthProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to listing account friends: %v", err)
	}
	result, err := proto.NewAuthClient(pc).ListUserFriends(ctx, &proto.ListUserRelativeRequest{
		UserId:    uint64(user.ID),
		IsRelated: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to listing account friends: %v", err)
	}

	out := lo.Map(result.Data, func(item *proto.SimpleUserInfo, index int) uint {
		return uint(item.Id)
	})

	var accounts []models.Account
	if err = database.C.Where("id IN ?", out).Find(&accounts).Error; err != nil {
		return nil, fmt.Errorf("failed to linking listed account friends: %v", err)
	}

	return accounts, nil
}

func ListAccountBlockedUsers(user models.Account) ([]models.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pc, err := gap.H.GetServiceGrpcConn(hyper.ServiceTypeAuthProvider)
	if err != nil {
		return nil, fmt.Errorf("failed to listing account blocked users: %v", err)
	}
	result, err := proto.NewAuthClient(pc).ListUserBlocklist(ctx, &proto.ListUserRelativeRequest{
		UserId:    uint64(user.ID),
		IsRelated: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to listing account blocked users: %v", err)
	}

	out := lo.Map(result.Data, func(item *proto.SimpleUserInfo, index int) uint {
		return uint(item.Id)
	})

	var accounts []models.Account
	if err = database.C.Where("id IN ?", out).Find(&accounts).Error; err != nil {
		return nil, fmt.Errorf("failed to linking listed blocked users: %v", err)
	}

	return accounts, nil
}

func ModifyPosterVoteCount(user models.Account, isUpvote bool, delta int) error {
	if isUpvote {
		user.TotalUpvote += delta
	} else {
		user.TotalDownvote += delta
	}

	return database.C.Save(&user).Error
}

func NotifyPosterAccount(user models.Account, title, body string, subtitle *string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	pc, err := gap.H.GetServiceGrpcConn(hyper.ServiceTypeAuthProvider)
	if err != nil {
		return err
	}
	_, err = proto.NewNotifierClient(pc).NotifyUser(ctx, &proto.NotifyUserRequest{
		UserId: uint64(user.ID),
		Notify: &proto.NotifyRequest{
			Topic:       "interactive.feedback",
			Title:       title,
			Subtitle:    subtitle,
			Body:        body,
			IsRealtime:  false,
			IsForcePush: true,
		},
	})
	if err != nil {
		log.Warn().Err(err).Msg("An error occurred when notify account...")
	} else {
		log.Debug().Uint("uid", user.ID).Msg("Notified account.")
	}

	return err
}
