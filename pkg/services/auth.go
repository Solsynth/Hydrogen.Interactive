package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.solsynth.dev/hydrogen/identity/pkg/grpc/proto"
	"git.solsynth.dev/hydrogen/interactive/pkg/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/grpc"
	"git.solsynth.dev/hydrogen/interactive/pkg/models"
	"gorm.io/gorm"
)

func LinkAccount(userinfo *proto.Userinfo) (models.Account, error) {
	var account models.Account
	if userinfo == nil {
		return account, fmt.Errorf("remote userinfo was not found")
	}
	if err := database.C.Where(&models.Account{
		ExternalID: uint(userinfo.Id),
	}).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			account = models.Account{
				Name:         userinfo.Name,
				Nick:         userinfo.Nick,
				Avatar:       userinfo.Avatar,
				Banner:       userinfo.Banner,
				Description:  userinfo.GetDescription(),
				EmailAddress: userinfo.Email,
				PowerLevel:   0,
				ExternalID:   uint(userinfo.Id),
			}
			return account, database.C.Save(&account).Error
		}
		return account, err
	}

	account.Name = userinfo.Name
	account.Nick = userinfo.Nick
	account.Avatar = userinfo.Avatar
	account.Banner = userinfo.Banner
	account.Description = userinfo.GetDescription()
	account.EmailAddress = userinfo.Email

	err := database.C.Save(&account).Error

	return account, err
}

func Authenticate(atk, rtk string) (models.Account, string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var err error
	var user models.Account
	reply, err := grpc.Auth.Authenticate(ctx, &proto.AuthRequest{
		AccessToken:  atk,
		RefreshToken: &rtk,
	})
	if err != nil {
		return user, reply.GetAccessToken(), reply.GetRefreshToken(), err
	} else if !reply.IsValid {
		return user, reply.GetAccessToken(), reply.GetRefreshToken(), fmt.Errorf("invalid authorization context")
	}

	user, err = LinkAccount(reply.Userinfo)

	return user, reply.GetAccessToken(), reply.GetRefreshToken(), err
}
