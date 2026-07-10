package db

import (
	"context"
	"go-chat/internal/adapter/dto"
	"go-chat/internal/domain"
	"go-chat/internal/port"
	"time"

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

func (r *ChatRepository) GetMessages(ctx context.Context, chatRoomID string, offset int) ([]dto.ChatMessagesRow, error) {
	var rows []dto.ChatMessagesRow

	query := r.db.WithContext(ctx).
		Table("chat_messages").
		Select(
			"chat_messages.id",
			"chat_messages.chat_id",
			"chat_messages.message",
			"chat_messages.url",
			"chat_messages.status",
			"chat_messages.action_status",
			"chat_messages.send_at",
			"chats.user_id",
		).
		Joins("JOIN chats ON chat_messages.chat_id = chats.id").
		Where("chats.chat_room_id = ?", chatRoomID).
		Order("chat_messages.send_at DESC").
		Limit(25).
		Offset(offset).
		Find(&rows)

	if query.Error != nil {
		return nil, query.Error
	}

	return rows, nil
}

func (r *ChatRepository) UpdateMessagesAsRead(ctx context.Context, chatRoomID string, userID string) error {
	tx := r.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return tx.Error
	}

	// Automatically rollback if the function exits before committing
	defer tx.Rollback()

	tx.Table("chat_messages").
		Select("status", "read_at").
		// Must using subquery, since the update use struct domain
		Where("chat_id IN (?)", tx.Table("chats").
			Select("id").
			Where("chat_room_id = ?", chatRoomID).
			Where("user_id != ?", userID),
		).
		Updates(domain.ChatMessage{
			Status: domain.StatusRead,
			ReadAt: time.Now(),
		})

	return tx.Commit().Error
}
