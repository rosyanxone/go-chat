package db

import (
	"context"
	"go-chat/internal/adapter/dto"
	"go-chat/internal/port"

	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) port.ChatRepository {
	return &ChatRepository{db}
}

func (r *ChatRepository) GetRooms(ctx context.Context, userID string) ([]dto.ChatRoomRow, error) {
	var rows []dto.ChatRoomRow

	selectQuery := `
		chat_rooms.id,
		chat_rooms.title,

		(SELECT message FROM chat_messages
		JOIN chats ON chats.id = chat_messages.chat_id
		WHERE chats.chat_room_id = chat_rooms.id
		ORDER BY chat_messages.send_at DESC LIMIT 1) as last_message,

		(SELECT send_at FROM chat_messages
		JOIN chats ON chats.id = chat_messages.chat_id
		WHERE chats.chat_room_id = chat_rooms.id
		ORDER BY chat_messages.send_at DESC LIMIT 1) as last_send_at,

		(SELECT users.name FROM chats
		JOIN users ON users.id = chats.user_id
		WHERE chats.chat_room_id = chat_rooms.id AND users.id != ?
		LIMIT 1) as other_user_name,

		(SELECT COUNT(chat_messages.id) FROM chat_messages
		JOIN chats ON chats.id = chat_messages.chat_id
		WHERE chats.chat_room_id = chat_rooms.id AND chats.user_id != ?
		AND chat_messages.status = 'sent') as total_unread
	`

	err := r.db.WithContext(ctx).
		Table("chat_rooms").
		Select(selectQuery, userID, userID).
		Where("EXISTS (SELECT 1 FROM chats WHERE chats.chat_room_id = chat_rooms.id AND chats.user_id = ?)", userID).
		Where("EXISTS (SELECT 1 FROM chat_messages JOIN chats ON chats.id = chat_messages.chat_id WHERE chats.chat_room_id = chat_rooms.id)").
		Order("last_send_at DESC").
		Find(&rows).
		Error

	return rows, err
}
