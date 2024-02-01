package services

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
	"code.smartsheep.studio/hydrogen/interactive/pkg/security"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type PassportUserinfo struct {
	Sub               string `json:"sub"`
	Email             string `json:"email"`
	Picture           string `json:"picture"`
	PreferredUsername string `json:"preferred_username"`
}

func LinkAccount(userinfo PassportUserinfo) (models.Account, error) {
	id, _ := strconv.Atoi(userinfo.Sub)

	var account models.Account
	if err := database.C.Where(&models.Account{
		ExternalID: uint(id),
	}).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			account = models.Account{
				Name:         userinfo.PreferredUsername,
				Avatar:       userinfo.Picture,
				EmailAddress: userinfo.Email,
				PowerLevel:   0,
				ExternalID:   uint(id),
			}
			return account, database.C.Save(&account).Error
		}
		return account, err
	}

	return account, nil
}

func GetToken(account models.Account) (string, string, error) {
	var err error
	var refresh, access string

	sub := strconv.Itoa(int(account.ID))
	access, err = security.EncodeJwt(
		uuid.NewString(),
		security.JwtAccessType,
		sub,
		[]string{"interactive"},
		time.Now().Add(30*time.Minute),
	)
	if err != nil {
		return refresh, access, err
	}
	refresh, err = security.EncodeJwt(
		uuid.NewString(),
		security.JwtRefreshType,
		sub,
		[]string{"interactive"},
		time.Now().Add(30*24*time.Hour),
	)
	if err != nil {
		return refresh, access, err
	}

	return access, refresh, nil
}

func RefreshToken(token string) (string, string, error) {
	parseInt := func(str string) int {
		val, _ := strconv.Atoi(str)
		return val
	}

	var account models.Account
	if claims, err := security.DecodeJwt(token); err != nil {
		return "404", "403", err
	} else if claims.Type != security.JwtRefreshType {
		return "404", "403", fmt.Errorf("invalid token type, expected refresh token")
	} else if err := database.C.Where(models.Account{
		BaseModel: models.BaseModel{ID: uint(parseInt(claims.Subject))},
	}).First(&account).Error; err != nil {
		return "404", "403", err
	}

	return GetToken(account)
}
