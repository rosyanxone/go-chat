package app

import (
	"context"
	"go-chat/internal/adapter/dto"
	"go-chat/internal/domain"
	"go-chat/internal/port"
	"strconv"
	"time"
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
