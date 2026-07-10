package port

import (
	"context"
	"go-chat/internal/adapter/dto"
)

type ChatRepository interface {
	GetRooms(ctx context.Context, userID string) ([]dto.ChatRoomRow, error)
	GetMessages(ctx context.Context, chatRoomID string, offset int) ([]dto.ChatMessagesRow, error)
	UpdateMessagesAsRead(ctx context.Context, chatRoomID string, userID string) error
}
