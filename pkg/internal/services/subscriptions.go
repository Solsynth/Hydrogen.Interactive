package services

import (
	"errors"
	"fmt"

	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/models"
	"gorm.io/gorm"
)

func GetSubscriptionOnUser(user models.Account, target models.Account) (*models.Subscription, error) {
	var subscription models.Subscription
	if err := database.C.Where("follower_id = ? AND account_id = ?", user.ID, target.ID).First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to get subscription: %v", err)
	}
	return &subscription, nil
}

func GetSubscriptionOnTag(user models.Account, target models.Tag) (*models.Subscription, error) {
	var subscription models.Subscription
	if err := database.C.Where("follower_id = ? AND tag_id = ?", user.ID, target.ID).First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to get subscription: %v", err)
	}
	return &subscription, nil
}

func GetSubscriptionOnCategory(user models.Account, target models.Category) (*models.Subscription, error) {
	var subscription models.Subscription
	if err := database.C.Where("follower_id = ? AND category_id = ?", user.ID, target.ID).First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to get subscription: %v", err)
	}
	return &subscription, nil
}

func SubscribeToUser(user models.Account, target models.Account) (models.Subscription, error) {
	var subscription models.Subscription
	if err := database.C.Where("follower_id = ? AND account_id = ?", user.ID, target.ID).First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return subscription, fmt.Errorf("subscription already exists")
		}
		return subscription, fmt.Errorf("unable to check subscription is exists or not: %v", err)
	}

	subscription = models.Subscription{
		FollowerID: user.ID,
		AccountID:  &target.ID,
	}

	err := database.C.Save(&subscription).Error
	return subscription, err
}

func SubscribeToTag(user models.Account, target models.Tag) (models.Subscription, error) {
	var subscription models.Subscription
	if err := database.C.Where("follower_id = ? AND tag_id = ?", user.ID, target.ID).First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return subscription, fmt.Errorf("subscription already exists")
		}
		return subscription, fmt.Errorf("unable to check subscription is exists or not: %v", err)
	}

	subscription = models.Subscription{
		FollowerID: user.ID,
		TagID:      &target.ID,
	}

	err := database.C.Save(&subscription).Error
	return subscription, err
}

func SubscribeToCategory(user models.Account, target models.Category) (models.Subscription, error) {
	var subscription models.Subscription
	if err := database.C.Where("follower_id = ? AND category_id = ?", user.ID, target.ID).First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return subscription, fmt.Errorf("subscription already exists")
		}
		return subscription, fmt.Errorf("unable to check subscription is exists or not: %v", err)
	}

	subscription = models.Subscription{
		FollowerID: user.ID,
		CategoryID: &target.ID,
	}

	err := database.C.Save(&subscription).Error
	return subscription, err
}

func UnsubscribeFromUser(user models.Account, target models.Account) error {
	var subscription models.Subscription
	if err := database.C.Where("follower_id = ? AND account_id = ?", user.ID, target.ID).First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("subscription does not exist")
		}
		return fmt.Errorf("unable to check subscription is exists or not: %v", err)
	}

	err := database.C.Delete(&subscription).Error
	return err
}

func UnsubscribeFromTag(user models.Account, target models.Tag) error {
	var subscription models.Subscription
	if err := database.C.Where("follower_id = ? AND tag_id = ?", user.ID, target.ID).First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("subscription does not exist")
		}
		return fmt.Errorf("unable to check subscription is exists or not: %v", err)
	}

	err := database.C.Delete(&subscription).Error
	return err
}

func UnsubscribeFromCategory(user models.Account, target models.Category) error {
	var subscription models.Subscription
	if err := database.C.Where("follower_id = ? AND category_id = ?", user.ID, target.ID).First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("subscription does not exist")
		}
		return fmt.Errorf("unable to check subscription is exists or not: %v", err)
	}

	err := database.C.Delete(&subscription).Error
	return err
}
