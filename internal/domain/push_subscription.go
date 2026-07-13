package domain

import "time"

type PushSubscription struct {
	ID                uint   `gorm:"primaryKey"`
	SubscribeableType string `gorm:"default:App\\Models\\User;index"`
	SubscribeableID   uint   `gorm:"index"`
	Endpoint          string `gorm:"type:text"`
	PublicKey         string
	AuthToken         string
	ContentEncoding   string
	CreatedAt         time.Time
	UpdatedAtAt       time.Time
}

type PushSubscriptionRequest struct {
	// UserID          uint   `json:"user_id" binding:"required"`
	Endpoint        string `json:"endpoint" binding:"required"`
	ContentEncoding string `json:"content_encoding"`
	Keys            struct {
		P256dh string `json:"p256dh"`
		Auth   string `json:"auth"`
	} `json:"keys"`
}

// PushPayload is the JSON body delivered to the browser's service worker
// on the "push" event (lookup public/sw.js -> event.data.json()).
type PushPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Icon  string `json:"icon,omitempty"`
	Url   string `json:"url,omitempty"`
}
