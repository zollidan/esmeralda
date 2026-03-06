package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type MessageHandler func(ctx context.Context, payload []byte) error

type Consumer struct {
    rdb      *redis.Client
    stream   string
    group    string
    consumer string
}

func NewConsumer(rdb *redis.Client, stream, group, consumer string) *Consumer {
    return &Consumer{rdb: rdb, stream: stream, group: group, consumer: consumer}
}

func (c *Consumer) ensureGroup(ctx context.Context) error {
    err := c.rdb.XGroupCreateMkStream(ctx, c.stream, c.group, "0").Err()
    if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
        return fmt.Errorf("xgroup create: %w", err)
    }
    return nil
}

func (c *Consumer) Consume(ctx context.Context, handler MessageHandler) error {
    if err := c.ensureGroup(ctx); err != nil {
        return err
    }

    readID := "0"

    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
        }

        streams, err := c.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
            Group:    c.group,
            Consumer: c.consumer,
            Streams:  []string{c.stream, readID},
            Count:    1,
            Block:    5 * time.Second,
        }).Result()

        if err != nil {
            if errors.Is(err, context.Canceled) {
                return nil
            }
            if errors.Is(err, redis.Nil) {
                if readID == "0" {
                    readID = ">"
                }
                continue
            }
            return fmt.Errorf("xreadgroup %s: %w", c.stream, err)
        }

        if readID == "0" && (len(streams) == 0 || len(streams[0].Messages) == 0) {
            readID = ">"
            continue
        }

        for _, s := range streams {
            for _, msg := range s.Messages {
                raw, ok := msg.Values["payload"].(string)
                if !ok {
                    log.Printf("[consumer:%s] unexpected payload type for msg %s", c.stream, msg.ID)
                    c.ack(ctx, msg.ID)
                    continue
                }

				handlerErr := handler(ctx, []byte(raw))
				c.ack(ctx, msg.ID)
				if err := c.rdb.XDel(ctx, c.stream, msg.ID).Err(); err != nil {
					log.Printf("[consumer:%s] xdel msg %s: %v", c.stream, msg.ID, err)
				}
				if handlerErr != nil {
					log.Printf("[consumer:%s] handler error msg %s: %v", c.stream, msg.ID, handlerErr)
				}
            }
        }
    }
}

func (c *Consumer) ack(ctx context.Context, msgID string) {
    if err := c.rdb.XAck(ctx, c.stream, c.group, msgID).Err(); err != nil {
        log.Printf("[consumer:%s] xack msg %s: %v", c.stream, msgID, err)
    }
}

func Unmarshal[T any](payload []byte) (T, error) {
    var v T
    if err := json.Unmarshal(payload, &v); err != nil {
        return v, fmt.Errorf("unmarshal: %w", err)
    }
    return v, nil
}