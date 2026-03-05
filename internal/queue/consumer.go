package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type MessageHandler func(ctx context.Context, payload []byte) error

type Consumer struct {
	rdb    *redis.Client
	stream string
}

func NewConsumer(rdb *redis.Client, stream string) *Consumer {
	return &Consumer{rdb: rdb, stream: stream}
}

func (c *Consumer) Consume(ctx context.Context, handler MessageHandler) error {
	lastID := "0"

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		streams, err := c.rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{c.stream, lastID},
			Count:   1,
			Block:   5*time.Second,
		}).Result()
		if err != nil {
			if err == context.Canceled {
				return nil
			}
			return fmt.Errorf("xread %s: %w", c.stream, err)
		}

		for _, s := range streams {
			for _, msg := range s.Messages {
				raw, ok := msg.Values["payload"].(string)
				if !ok {
					log.Printf("[consumer:%s] unexpected payload type: %T", c.stream, msg.Values["payload"])
					lastID = msg.ID
					continue
				}

				if err := handler(ctx, []byte(raw)); err != nil {
					log.Printf("[consumer:%s] handler error: %v", c.stream, err)
					continue
				}

				lastID = msg.ID
			}
		}

	}
}

func Unmarshal[T any](payload []byte) (T, error) {
	var v T
	if err := json.Unmarshal(payload, &v); err != nil {
		return v, fmt.Errorf("unmarshal: %w", err)
	}
	return v, nil
}
