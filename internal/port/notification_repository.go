package port

import (
	"context"
	"go-chat/internal/domain"
)

type NotificationRepository interface {
	UpdateUserSubscription(ctx context.Context, userID uint, req domain.PushSubscriptionRequest) error
	UserUnsubscribe(ctx context.Context, endpoint string) error
}
