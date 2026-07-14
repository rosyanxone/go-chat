package db

import (
	"context"
	"errors"
	"go-chat/internal/adapter/dto"
	"go-chat/internal/domain"
	"go-chat/internal/port"
	"strconv"
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

func (r *ChatRepository) GetChat(ctx context.Context, senderID string, targetID string) (*domain.Chat, error) {
	var chat domain.Chat

	subQuery := r.db.Table("chats").
		Select("chat_room_id").
		Where("user_id IN (?)", []string{senderID, targetID}).
		Group("chat_room_id").
		Having("COUNT(DISTINCT user_id) = ?", 2)

	err := r.db.WithContext(ctx).
		Where("user_id", senderID).
		Where("chat_room_id IN (?)", subQuery).
		First(&chat).
		Error

	if err == nil {
		return &chat, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	tx := r.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return nil, tx.Error
	}

	// Automatically rollback if the function exits before committing
	defer tx.Rollback()

	chatRoom := domain.ChatRoom{
		Type: domain.TypePrivate,
	}

	err = tx.Create(&chatRoom).Error

	if err != nil {
		return nil, err
	}

	parsedTargetID, err := strconv.ParseUint(targetID, 10, 64)

	if err != nil {
		return nil, err
	}

	chat = domain.Chat{
		ChatRoomID: chatRoom.ID,
		UserID:     uint(parsedTargetID),
	}

	err = tx.Create(&chat).Error

	if err != nil {
		return nil, err
	}

	tx.Commit()

	return &chat, nil
}

func (r *ChatRepository) GetTotalUnread(ctx context.Context, chatID string) (*uint64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&domain.ChatMessage{}).
		Where("chat_id = ?", chatID).
		Where("status = ?", domain.StatusSent).
		Count(&count).
		Error

	if err != nil {
		return nil, err
	}

	finalCount := uint64(count)

	return &finalCount, nil
}

func (r *ChatRepository) CreateNewMessage(ctx context.Context, chatMessage *domain.ChatMessage) error {
	tx := r.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return tx.Error
	}

	// Automatically rollback if the function exits before committing
	defer tx.Rollback()

	err := tx.Create(&chatMessage).Error

	if err != nil {
		return err
	}

	return tx.Commit().Error
}

func (r *ChatRepository) UpdateMessagesAsRead(ctx context.Context, chatRoomID string, userID string) error {
	tx := r.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return tx.Error
	}

	// Automatically rollback if the function exits before committing
	defer tx.Rollback()

	now := time.Now()

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
			ReadAt: &now,
		})

	return tx.Commit().Error
}
