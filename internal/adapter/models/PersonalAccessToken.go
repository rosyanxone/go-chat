package models

import (
	"time"

	"gorm.io/gorm"
)

type PersonalAccessToken struct {
	gorm.Model // Adds ID, CreatedAt, UpdatedAt, DeletedAt

	UserID uint64 `gorm:"index;not null"`

	Name      string `gorm:"size:100;not null"`
	TokenHash string `gorm:"size:64;uniqueIndex;not null"`

	Abilities string `gorm:"type:json"`

	LastUsedAt *time.Time
	ExpiresAt  *time.Time
	RevokedAt  *time.Time

	IPAddress string `gorm:"size:100"`
	UserAgent string `gorm:"type:text"`
}
