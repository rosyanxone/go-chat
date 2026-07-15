package app

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"go-chat/internal/adapter/dto"
	"go-chat/internal/domain"
	"go-chat/internal/port"
	"go-chat/internal/shared/convert"
)

type ChatService struct {
	repo port.ChatRepository
}

func NewChatService(r port.ChatRepository) *ChatService {
	return &ChatService{repo: r}
}

func (s *ChatService) GetRooms(ctx context.Context, userID string) ([]dto.ChatRoomResponse, error) {
	// limit := 10
	// offset := (page - 1) * limit

	rows, err := s.repo.GetRooms(ctx, userID)

	if err != nil {
		return nil, err
	}

	var responses []dto.ChatRoomResponse

	// Loop through and map the data (Equivalent to Laravel's map function)
	for _, row := range rows {

		// Fallback Title
		title := row.Title
		if title == "" {
			title = row.OtherUserName
		}

		// Truncate Message (Equivalent to Str::limit)
		lastMessage := row.LastMessage
		if len(lastMessage) > 50 {
			lastMessage = lastMessage[:47] + "..."
		}

		// Map to final struct
		responses = append(responses, dto.ChatRoomResponse{
			ID:          row.ID,
			Title:       title,
			LastMessage: lastMessage,
			LastSendAt:  formatLastSendAt(row.LastSendAt),
			TotalUnread: row.TotalUnread,
		})
	}

	return responses, nil
}

func formatLastSendAt(t *time.Time) string {
	if t == nil {
		return ""
	}

	now := time.Now()

	// isToday()
	if t.Year() == now.Year() && t.YearDay() == now.YearDay() {
		return t.Format("15:04") // H:i format
	}

	// isCurrentWeek()
	y1, w1 := t.ISOWeek()
	y2, w2 := now.ISOWeek()
	if y1 == y2 && w1 == w2 {
		return t.Format("Monday") // l format
	}

	// Else
	return t.Format("02/01/2006") // d/m/Y format
}

func (s *ChatService) GetAndReadMessages(ctx context.Context, chatRoomID string, userID string, page int) ([]dto.ChatMessageResponse, error) {
	err := s.repo.UpdateMessagesAsRead(ctx, chatRoomID, userID)

	if err != nil {
		return nil, err
	}

	limit := 25
	offset := (page - 1) * limit
	rows, err := s.repo.GetMessages(ctx, chatRoomID, offset)

	if err != nil {
		return nil, err
	}

	var responses []dto.ChatMessageResponse

	for _, row := range rows {
		isUser := strconv.FormatUint(uint64(row.UserID), 10) == userID
		isRead := row.Status == domain.StatusRead
		date := row.SendAt.Format("2 Jan 2006") // d M Y format
		time := row.SendAt.Format("15:04")      // H:i format

		responses = append(responses, dto.ChatMessageResponse{
			ID:           row.ID,
			Message:      row.Message,
			Url:          row.Url,
			Status:       string(row.Status),
			ActionStatus: (*string)(row.ActionStatus),
			IsUser:       isUser,
			IsRead:       isRead,
			Date:         date,
			Time:         time,
		})
	}

	return responses, nil
}

func (s *ChatService) GetChat(ctx context.Context, senderID uint, targetID uint) (*domain.Chat, error) {
	strUserSenderID := convert.UintToString(senderID)
	strUserTargetID := convert.UintToString(targetID)

	return s.repo.GetChat(ctx, strUserSenderID, strUserTargetID)
}

func (s *ChatService) GetTotalUnread(ctx context.Context, chatID string) (*uint64, error) {
	return s.repo.GetTotalUnread(ctx, chatID)
}

func (s *ChatService) GetRoomChatUrl(ctx context.Context, chatID uint, chatRoomID uint) (*string, error) {
	strChatID := convert.UintToString(chatID)

	unreadTotal, err := s.repo.GetTotalUnread(ctx, strChatID)

	if err != nil {
		return nil, err
	}

	strUnreadTotal := convert.UintToString(uint(*unreadTotal))
	strChatRoomID := convert.UintToString(chatRoomID)

	appDomain := os.Getenv("CHAT_DOMAIN")

	url := fmt.Sprintf("%s/chats/%s?unread=%s", appDomain, strChatRoomID, strUnreadTotal)

	return &url, nil
}

func (s *ChatService) CreateNewMessage(ctx context.Context, cmd dto.CreateMessageCommand) error {
	var actionStatus *string
	var uniqueCode *string

	if cmd.UniqueCode != nil {
		uniqueCode = convert.NullIfEmpty(*cmd.UniqueCode)
	}

	if uniqueCode != nil {
		actionStatus = convert.StringPtr(string(domain.ActionStatusPending))
	} else {
		actionStatus = nil
	}

	chatMessage := domain.ChatMessage{
		ChatID:       cmd.ChatID,
		Message:      cmd.Message,
		Url:          cmd.Url,
		UniqueCode:   cmd.UniqueCode,
		ActionStatus: (*domain.ChatMessageActionStatus)(actionStatus),
	}

	return s.repo.CreateNewMessage(ctx, &chatMessage)
}
