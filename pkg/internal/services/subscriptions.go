package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"git.solsynth.dev/hydrogen/dealer/pkg/hyper"
	"git.solsynth.dev/hydrogen/dealer/pkg/proto"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/database"
	"git.solsynth.dev/hydrogen/interactive/pkg/internal/gap"
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

func GetSubscriptionOnRealm(user models.Account, target models.Realm) (*models.Subscription, error) {
	var subscription models.Subscription
	if err := database.C.Where("follower_id = ? AND realm_id = ?", user.ID, target.ID).First(&subscription).Error; err != nil {
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
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return subscription, fmt.Errorf("subscription already exists")
		}
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
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return subscription, fmt.Errorf("subscription already exists")
		}
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
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return subscription, fmt.Errorf("subscription already exists")
		}
	}

	subscription = models.Subscription{
		FollowerID: user.ID,
		CategoryID: &target.ID,
	}

	err := database.C.Save(&subscription).Error
	return subscription, err
}

func SubscribeToRealm(user models.Account, target models.Realm) (models.Subscription, error) {
	var subscription models.Subscription
	if err := database.C.Where("follower_id = ? AND realm_id = ?", user.ID, target.ID).First(&subscription).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return subscription, fmt.Errorf("subscription already exists")
		}
	}

	subscription = models.Subscription{
		FollowerID: user.ID,
		RealmID:    &target.ID,
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

func UnsubscribeFromRealm(user models.Account, target models.Realm) error {
	var subscription models.Subscription
	if err := database.C.Where("follower_id = ? AND realm_id = ?", user.ID, target.ID).First(&subscription).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("subscription does not exist")
		}
		return fmt.Errorf("unable to check subscription is exists or not: %v", err)
	}

	err := database.C.Delete(&subscription).Error
	return err
}

func NotifyUserSubscription(poster models.Account, content string, title *string) error {
	var subscriptions []models.Subscription
	if err := database.C.Where("account_id = ?", poster.ID).Preload("Follower").Find(&subscriptions).Error; err != nil {
		return fmt.Errorf("unable to get subscriptions: %v", err)
	}

	nTitle := fmt.Sprintf("New post from %s (%s)", poster.Nick, poster.Name)
	nSubtitle := "From your subscription"

	body := TruncatePostContentShort(content)
	if title != nil {
		body = fmt.Sprintf("%s\n%s", *title, body)
	}

	userIDs := make([]uint64, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		userIDs = append(userIDs, uint64(subscription.Follower.ID))
	}

	pc, err := gap.H.GetServiceGrpcConn(hyper.ServiceTypeAuthProvider)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = proto.NewNotifierClient(pc).NotifyUserBatch(ctx, &proto.NotifyUserBatchRequest{
		UserId: userIDs,
		Notify: &proto.NotifyRequest{
			Topic:       "interactive.subscription",
			Title:       nTitle,
			Subtitle:    &nSubtitle,
			Body:        body,
			IsRealtime:  false,
			IsForcePush: true,
		},
	})

	return err
}

func NotifyTagSubscription(poster models.Tag, og models.Account, content string, title *string) error {
	var subscriptions []models.Subscription
	if err := database.C.Where("tag_id = ?", poster.ID).Preload("Follower").Find(&subscriptions).Error; err != nil {
		return fmt.Errorf("unable to get subscriptions: %v", err)
	}

	nTitle := fmt.Sprintf("New post in %s by %s (%s)", poster.Name, og.Nick, og.Name)
	nSubtitle := "From your subscription"

	body := TruncatePostContentShort(content)
	if title != nil {
		body = fmt.Sprintf("%s\n%s", *title, body)
	}

	userIDs := make([]uint64, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		userIDs = append(userIDs, uint64(subscription.Follower.ID))
	}

	pc, err := gap.H.GetServiceGrpcConn(hyper.ServiceTypeAuthProvider)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = proto.NewNotifierClient(pc).NotifyUserBatch(ctx, &proto.NotifyUserBatchRequest{
		UserId: userIDs,
		Notify: &proto.NotifyRequest{
			Topic:       "interactive.subscription",
			Title:       nTitle,
			Subtitle:    &nSubtitle,
			Body:        body,
			IsRealtime:  false,
			IsForcePush: true,
		},
	})

	return err
}

func NotifyCategorySubscription(poster models.Category, og models.Account, content string, title *string) error {
	var subscriptions []models.Subscription
	if err := database.C.Where("category_id = ?", poster.ID).Preload("Follower").Find(&subscriptions).Error; err != nil {
		return fmt.Errorf("unable to get subscriptions: %v", err)
	}

	nTitle := fmt.Sprintf("New post in %s by %s (%s)", poster.Name, og.Nick, og.Name)
	nSubtitle := "From your subscription"

	body := TruncatePostContentShort(content)
	if title != nil {
		body = fmt.Sprintf("%s\n%s", *title, body)
	}

	userIDs := make([]uint64, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		userIDs = append(userIDs, uint64(subscription.Follower.ID))
	}

	pc, err := gap.H.GetServiceGrpcConn(hyper.ServiceTypeAuthProvider)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = proto.NewNotifierClient(pc).NotifyUserBatch(ctx, &proto.NotifyUserBatchRequest{
		UserId: userIDs,
		Notify: &proto.NotifyRequest{
			Topic:       "interactive.subscription",
			Title:       nTitle,
			Subtitle:    &nSubtitle,
			Body:        body,
			IsRealtime:  false,
			IsForcePush: true,
		},
	})

	return err
}

func NotifyRealmSubscription(poster models.Realm, og models.Account, content string, title *string) error {
	var subscriptions []models.Subscription
	if err := database.C.Where("realm_id = ?", poster.ID).Preload("Follower").Find(&subscriptions).Error; err != nil {
		return fmt.Errorf("unable to get subscriptions: %v", err)
	}

	nTitle := fmt.Sprintf("New post in %s by %s (%s)", poster.Name, og.Nick, og.Name)
	nSubtitle := "From your subscription"

	body := TruncatePostContentShort(content)
	if title != nil {
		body = fmt.Sprintf("%s\n%s", *title, body)
	}

	userIDs := make([]uint64, 0, len(subscriptions))
	for _, subscription := range subscriptions {
		userIDs = append(userIDs, uint64(subscription.Follower.ID))
	}

	pc, err := gap.H.GetServiceGrpcConn(hyper.ServiceTypeAuthProvider)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = proto.NewNotifierClient(pc).NotifyUserBatch(ctx, &proto.NotifyUserBatchRequest{
		UserId: userIDs,
		Notify: &proto.NotifyRequest{
			Topic:       "interactive.subscription",
			Title:       nTitle,
			Subtitle:    &nSubtitle,
			Body:        body,
			IsRealtime:  false,
			IsForcePush: true,
		},
	})

	return err
}
