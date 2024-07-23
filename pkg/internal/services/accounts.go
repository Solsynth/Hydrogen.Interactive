package services

import (
	"context"
	"git.solsynth.dev/hydrogen/dealer/pkg/hyper"
	"git.solsynth.dev/hydrogen/dealer/pkg/proto"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/gap"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"github.com/rs/zerolog/log"
	"time"
)

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
		UserId: uint64(user.ExternalID),
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
		log.Debug().Uint("uid", user.ExternalID).Msg("Notified account.")
	}

	return err
}
