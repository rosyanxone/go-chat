package app

import (
	"context"
	"go-chat/internal/domain"
	"go-chat/internal/port"
)

type NotificationService struct {
	repo port.NotificationRepository
}

func NewNotificationService(r port.NotificationRepository) *NotificationService {
	return &NotificationService{repo: r}
}

func (s *NotificationService) UpdateUserSubscription(ctx context.Context, userID uint, req domain.PushSubscriptionRequest) error {
	return s.repo.UpdateUserSubscription(ctx, userID, req)
}

func (s *NotificationService) UserUnsubscribe(ctx context.Context, endpoint string) error {
	return s.repo.UserUnsubscribe(ctx, endpoint)
}
