package models

import (
	"time"
)

type PersonalAccessToken struct {
	// gorm.Model // Adds ID, CreatedAt, UpdatedAt, DeletedAt
	ID uint `gorm:"primaryKey"`

	TokenableType string `gorm:"default:App\\Models\\User"`
	TokenableID   uint64 `gorm:"index;not null"`

	Name      string `gorm:"size:100;not null"`
	Token     string `gorm:"size:64;uniqueIndex;not null"`
	Abilities string `gorm:"type:json;default:['*']"`

	LastUsedAt *time.Time
	ExpiresAt  *time.Time

	// RevokedAt  *time.Time
	// IPAddress string `gorm:"size:100"`
	// UserAgent string `gorm:"type:text"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
