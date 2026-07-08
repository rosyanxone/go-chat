package domain

import "time"

type Chat struct {
	ID         uint `gorm:"primaryKey"`
	ChatRoomID uint
	UserID     uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
