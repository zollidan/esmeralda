package queue

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Publisher interface {
    Publish(ctx context.Context, payload any) (string, error)
    Delete(ctx context.Context, msgID string) error
}

type Producer struct {
	rdb    *redis.Client
	stream string
}

func NewProducer(rdb *redis.Client, stream string) *Producer {
	return &Producer{rdb: rdb, stream: stream}
}

func (p *Producer) Publish(ctx context.Context, payload any) (string, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal payload: %w", err)
	}

	msgID, err := p.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: p.stream,
		Values: map[string]any{
			"payload": string(data),
		},
	}).Result()
	if err != nil {
		return "", fmt.Errorf("xadd to %s: %w", p.stream, err)
	}

	return msgID, nil
}

func (p *Producer) Delete(ctx context.Context, msgID string) error {
	if err := p.rdb.XDel(ctx, p.stream, msgID).Err(); err != nil {
		return fmt.Errorf("xdel from %s: %w", p.stream, err)
	}
	return nil
}
