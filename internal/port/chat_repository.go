package port

import (
	"context"
	"go-chat/internal/adapter/dto"
	"go-chat/internal/domain"
)

type ChatRepository interface {
	GetRooms(ctx context.Context, userID string, offset int) ([]dto.ChatRoomRow, error)
	GetMessages(ctx context.Context, chatRoomID string, offset int) ([]dto.ChatMessagesRow, error)
	GetChat(ctx context.Context, senderID string, targetID string) (*domain.Chat, error)
	GetTotalUnread(ctx context.Context, chatID string) (*uint64, error)
	GetMemberInfoByChatRoomId(ctx context.Context, userID string, chatRoomID string) (*domain.User, error)
	CreateNewMessage(ctx context.Context, chatMessage *domain.ChatMessage) error
	UpdateMessagesAsRead(ctx context.Context, chatRoomID string, userID string) error
}
