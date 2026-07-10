package dto

import (
	"go-chat/internal/domain"
	"time"
)

// ChatRoomRow catches the exact columns from custom SQL query
type ChatRoomRow struct {
	ID            uint
	Title         string
	LastMessage   string
	LastSendAt    *time.Time
	OtherUserName string
	TotalUnread   int
}

type ChatRoomResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	LastMessage string `json:"last_message"`
	LastSendAt  string `json:"last_send_at"`
	TotalUnread int    `json:"total_unread"`
}

type ChatMessagesRow struct {
	domain.ChatMessage
	UserID uint
}

type ChatMessageResponse struct {
	ID           uint    `json:"id"`
	Message      string  `json:"message"`
	Url          *string `json:"url"`
	Status       string  `json:"status"`
	ActionStatus *string `json:"action_status"`
	IsUser       bool    `json:"is_user"`
	IsRead       bool    `json:"is_read"`
	Date         string  `json:"date"`
	Time         string  `json:"time"`
}
