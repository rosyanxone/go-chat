package db

import (
	"context"
	"errors"
	"go-chat/internal/adapter/dto"
	"go-chat/internal/domain"
	"go-chat/internal/port"
	"go-chat/internal/shared/convert"
	"time"

	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) port.ChatRepository {
	return &ChatRepository{db}
}

func (r *ChatRepository) GetRooms(ctx context.Context, userID string, offset int) ([]dto.ChatRoomRow, error) {
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
		Limit(15).
		Offset(offset).
		Find(&rows).
		Error

	return rows, err
}

func (r *ChatRepository) GetMessages(ctx context.Context, chatRoomID string, offset int) ([]dto.ChatMessagesRow, error) {
	var rows []dto.ChatMessagesRow

	err := r.db.WithContext(ctx).
		Table("chat_messages AS cm").
		Select(
			"cm.id",
			"cm.chat_id",
			"cm.message",
			"cm.url",
			"cm.status",
			"cm.action_status",
			"cm.send_at",
			"c.user_id",
		).
		Joins("JOIN chats AS c ON c.id = cm.chat_id").
		Where("c.chat_room_id = ?", chatRoomID).
		Order("cm.send_at DESC").
		Order("cm.id DESC").
		Limit(25).
		Offset(offset).
		Find(&rows).
		Error

	if err != nil {
		return nil, err
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

	usersID := make([]uint, 0, 2)

	for _, id := range []string{senderID, targetID} {
		parsedID := convert.StringToInt(id)

		usersID = append(usersID, uint(parsedID))
	}

	chats := []domain.Chat{
		{
			ChatRoomID: chatRoom.ID,
			UserID:     usersID[0],
		},
		{
			ChatRoomID: chatRoom.ID,
			UserID:     usersID[1],
		},
	}

	if err := tx.Create(&chats).Error; err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return &chats[0], nil // index 0 is the sender chat
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

func (r *ChatRepository) GetMemberInfoByChatRoomId(ctx context.Context, userID string, chatRoomID string) (*domain.User, error) {
	var user domain.User

	err := r.db.WithContext(ctx).
		Table("chats").
		Select("users.id", "users.name", "users.phone_number").
		Joins("JOIN users ON chats.user_id = users.id").
		Where("chats.chat_room_id = ?", chatRoomID).
		Where("users.id != ?", userID).
		First(&user).
		Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *ChatRepository) CreateNewMessage(ctx context.Context, chatMessage *domain.ChatMessage) error {
	tx := r.db.WithContext(ctx).Begin()

	if tx.Error != nil {
		return tx.Error
	}

	// Automatically rollback if the function exits before committing
	defer tx.Rollback()

	err := tx.Create(chatMessage).Error

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

	chatsId := tx.Table("chats").
		Select("id").
		Where("chat_room_id = ?", chatRoomID).
		Where("user_id != ?", userID)

	err := tx.Model(&domain.ChatMessage{}).
		Where("chat_id IN (?)", chatsId).
		Where("read_at IS NULL").
		Updates(domain.ChatMessage{
			Status: domain.StatusRead,
			ReadAt: &now,
		}).
		Error

	if err != nil {
		return err
	}

	return tx.Commit().Error
}
