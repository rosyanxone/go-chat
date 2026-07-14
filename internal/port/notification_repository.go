package port

import (
	"context"
	"go-chat/internal/adapter/dto"
	"go-chat/internal/domain"
)

type NotificationRepository interface {
	UpdateUserSubscription(ctx context.Context, userID uint, req dto.PushSubscriptionRequest) error
	GetSubscriptionsByUser(ctx context.Context, userID uint) ([]domain.PushSubscription, error)
	UserUnsubscribe(ctx context.Context, endpoint string) error
}
