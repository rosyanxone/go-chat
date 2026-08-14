package websocket

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type Event struct {
	Event   string      `json:"event"`
	Channel string      `json:"channel"`
	Data    interface{} `json:"data"`
}

type Broadcaster struct {
	rdb *redis.Client
}

func NewBroadcaster(rdb *redis.Client) *Broadcaster {
	return &Broadcaster{
		rdb: rdb,
	}
}

func (b *Broadcaster) Broadcast(
	ctx context.Context,
	channel string,
	event string,
	data interface{},
) error {
	payload := Event{
		Event:   event,
		Channel: channel,
		Data:    data,
	}

	jsonData, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	return b.rdb.Publish(ctx, channel, jsonData).Err()
}
