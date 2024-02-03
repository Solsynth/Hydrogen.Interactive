package services

import (
	"code.smartsheep.studio/hydrogen/interactive/pkg/database"
	"code.smartsheep.studio/hydrogen/interactive/pkg/models"
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
