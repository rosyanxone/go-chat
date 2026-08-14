package app

import (
	"context"

	"go-chat/internal/adapter/websocket"
)

type BroadcastService struct {
	broadcaster *websocket.Broadcaster
}

func NewBroadcastService(
	broadcaster *websocket.Broadcaster,
) *BroadcastService {
	return &BroadcastService{
		broadcaster: broadcaster,
	}
}

func (s *BroadcastService) Send(
	ctx context.Context,
	channel string,
	event string,
	data interface{},
) error {
	return s.broadcaster.Broadcast(
		ctx,
		channel,
		event,
		data,
	)
}
