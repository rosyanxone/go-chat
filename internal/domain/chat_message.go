package domain

import "time"

type ChatMessage struct {
	ID uint `gorm:"primaryKey"`

	ChatID  uint
	Message string
	Url     string

	ReadAt time.Time
	SendAt time.Time

	Status       ChatMessageStatus       `gorm:"type:enum('read', 'sent');default:'sent'"`
	ActionStatus ChatMessageActionStatus `gorm:"type:enum('done', 'pending')"`
	UniqueCode   string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChatMessageStatus string
type ChatMessageActionStatus string

const (
	StatusRead          ChatMessageStatus       = "read"
	StatusSent          ChatMessageStatus       = "sent"
	ActionStatusDone    ChatMessageActionStatus = "done"
	ActionStatusPending ChatMessageActionStatus = "pending"
)
