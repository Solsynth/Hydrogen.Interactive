package services

import (
	"errors"
	"fmt"
	"git.solsynth.dev/hydrogen/dealer/pkg/proto"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"gorm.io/gorm"
	"reflect"
)

func LinkAccount(userinfo *proto.UserInfo) (models.Account, error) {
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

	prev := account
	account.Name = userinfo.Name
	account.Nick = userinfo.Nick
	account.Avatar = userinfo.Avatar
	account.Banner = userinfo.Banner
	account.Description = userinfo.GetDescription()
	account.EmailAddress = userinfo.Email

	var err error
	if !reflect.DeepEqual(prev, account) {
		err = database.C.Save(&account).Error
	}

	return account, err
}
