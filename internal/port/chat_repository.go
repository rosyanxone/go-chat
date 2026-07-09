package port

import (
	"context"
	"go-chat/internal/adapter/dto"
)

type ChatRepository interface {
	GetRooms(ctx context.Context, userID string) ([]dto.ChatRoomRow, error)
}
