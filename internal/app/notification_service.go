package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-chat/internal/adapter/dto"
	"go-chat/internal/port"
	"log"
	"net/http"

	"github.com/SherClockHolmes/webpush-go"
)

type NotificationService struct {
	repo            port.NotificationRepository
	vapidPublicKey  string
	vapidPrivateKey string
	vapidSubject    string
}

func NewNotificationService(r port.NotificationRepository, vapidPublicKey, vapidPrivateKey, vapidSubject string) *NotificationService {
	return &NotificationService{
		r,
		vapidPublicKey,
		vapidPrivateKey,
		vapidSubject,
	}
}

func (s *NotificationService) UpdateUserSubscription(ctx context.Context, userID uint, req dto.PushSubscriptionRequest) error {
	return s.repo.UpdateUserSubscription(ctx, userID, req)
}

func (s *NotificationService) UserUnsubscribe(ctx context.Context, endpoint string) error {
	return s.repo.UserUnsubscribe(ctx, endpoint)
}

// VAPIDPublicKey exposes the public key so frontend can pass it to
// PushManager.subscribe({ applicationServerKey: ... }).
func (s *NotificationService) VAPIDPublicKey() string {
	return s.vapidPublicKey
}

// SendToUser fans a notification out to every device/browser the user is
// currently subscribed on. It's best-effort: one failing endpoint doesn't
// stop delivery to the rest. Subscriptions the browser reports as gone
// (410/404) are deleted automatically so they stop being retried.
func (s *NotificationService) SendToUser(ctx context.Context, userID uint, payload dto.PushPayload) error {
	subs, err := s.repo.GetSubscriptionsByUser(ctx, userID)

	if err != nil {
		return err
	}

	if len(subs) == 0 {
		return nil
	}

	message, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	var sendErrs []error

	for _, sub := range subs {
		resp, err := webpush.SendNotification(message, &webpush.Subscription{
			Endpoint: sub.Endpoint,
			Keys: webpush.Keys{
				P256dh: sub.PublicKey,
				Auth:   sub.AuthToken,
			},
		}, &webpush.Options{
			Subscriber:      s.vapidSubject,
			VAPIDPublicKey:  s.vapidPublicKey,
			VAPIDPrivateKey: s.vapidPrivateKey,
			TTL:             30,
		})

		if err != nil {
			log.Printf("push: failed to send to %s: %v", sub.Endpoint, err)
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusGone || resp.StatusCode == http.StatusNotFound {
			err := s.repo.UserUnsubscribe(ctx, sub.Endpoint)

			if err != nil {
				log.Printf("push: failed to clean up dead subscription %s: %v", sub.Endpoint, err)
			}
			continue
		}

		// Anything outside the 200-299 range that isn't a known "gone"
		// status (e.g. 400/401/403 from bad VAPID keys) still counts as
		// a failure the caller should know about.
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			sendErrs = append(sendErrs, fmt.Errorf("endpoint %s: push service returned status %d", sub.Endpoint, resp.StatusCode))
		}
	}

	return errors.Join(sendErrs...)
}
