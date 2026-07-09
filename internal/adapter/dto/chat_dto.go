package dto

import "time"

// ChatRoomRow catches the exact columns from custom SQL query
type ChatRoomRow struct {
	ID            uint
	Title         string
	LastMessage   string
	LastSendAt    *time.Time
	OtherUserName string
	TotalUnread   int
}

// ChatRoomResponse final JSON output for the frontend
type ChatRoomResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	LastMessage string `json:"last_message"`
	LastSendAt  string `json:"last_send_at"`
	TotalUnread int    `json:"total_unread"`
}
