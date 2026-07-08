package db

import (
	"context"
	"errors"
	"fmt"
	"go-chat/internal/domain"
	"go-chat/internal/port"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) port.NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r NotificationRepository) UpdateUserSubscription(ctx context.Context, userID uint, req domain.PushSubscriptionRequest) error {
	tx := r.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return tx.Error
	}

	// Automatically rollback if the function exits before committing
	defer tx.Rollback()

	// Delete existing subscriptions for this specific user
	err := tx.Where("subscribeable_type = ? AND subscribeable_ID = ?", "App\\Models\\User", userID).
		Delete(&domain.PushSubscription{}).
		Error

	if err != nil {
		return err
	}

	newSub := domain.PushSubscription{
		SubscribeableID: userID,
		Endpoint:        req.Endpoint,
		PublicKey:       req.Keys.P256dh,
		AuthToken:       req.Keys.Auth,
		ContentEncoding: req.ContentEncoding,
	}

	err = tx.Create(&newSub).Error

	if err != nil {
		return err
	}

	return tx.Commit().Error
}

func (r NotificationRepository) UserUnsubscribe(ctx context.Context, endpoint string) error {
	var pushSubs domain.PushSubscription

	tx := r.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return tx.Error
	}

	// Automatically rollback if the function exits before committing
	defer tx.Rollback()

	err := tx.Where("endpoint = ?", endpoint).First(&pushSubs).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid or deleted endpoint")
		}

		return fmt.Errorf("database error during endpoint lookup: %w", err)
	}

	err = tx.Delete(&pushSubs).Error

	if err != nil {
		return err
	}

	return tx.Commit().Error
}
