package domain

import "time"

type ChatRoom struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	Type      ChatRoomType `gorm:"type:enum('private', 'group');default:'private'"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChatRoomType string

const (
	TypePrivate ChatRoomType = "private"
	TypeGroup   ChatRoomType = "group"
)
