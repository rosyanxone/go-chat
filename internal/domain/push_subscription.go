package domain

import "time"

type PushSubscription struct {
	ID               uint   `gorm:"primaryKey"`
	SubscribableType string `gorm:"default:App\\Models\\User;index"`
	SubscribableID   uint   `gorm:"index"`
	Endpoint         string `gorm:"type:text"`
	PublicKey        string
	AuthToken        string
	ContentEncoding  string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
