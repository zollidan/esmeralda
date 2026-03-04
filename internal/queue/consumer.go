package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zollidan/esmeralda/internal/api"
	"github.com/zollidan/esmeralda/internal/processor"
	"gorm.io/gorm"
)

type Consumer struct {
	rdb *redis.Client
	client *api.Client
	db *gorm.DB
}

func NewConsumer(rdb *redis.Client, client *api.Client, db *gorm.DB) *Consumer {
	return &Consumer{rdb: rdb, client: client, db: db}
}

func (c *Consumer) DeleteTask(ctx context.Context, id string) error {
	_, err := c.rdb.XDel(ctx, StreamName, id).Result()
	return err
}

func (c *Consumer) Consume(ctx context.Context) error {
	lastID := "0"

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		streams, err := c.rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{StreamName, lastID},
			Count:   1,
			Block:   0,
		}).Result()
		if err != nil {
			if err == context.Canceled {
				return nil
			}
			return fmt.Errorf("xread: %w", err)
		}

		for _, s := range streams {
			for _, msg := range s.Messages {
				raw, ok := msg.Values["payload"].(string)
				if !ok {
					log.Printf("unexpected payload type: %#v", msg.Values["payload"])
					continue
				}

				var task ParseTask
				if err := json.Unmarshal([]byte(raw), &task); err != nil {
					log.Printf("unmarshal task: %v", err)
					continue
				}

				fmt.Printf("received task: %+v\n", task)

				date, err := time.Parse("02.01.2006", task.Date)
				if err != nil {
					log.Printf("parse date: %v", err)
					continue
				}

				matches, totalMatches, err := c.client.GetMatches(api.MatchesFilter{
					Date: date,
				})
				if err != nil {
					log.Printf("get matches: %v", err)
					continue
				}

				log.Printf("found %d matches for date %s", totalMatches, date.Format("02.01.2006"))
				
				err = processor.ProcessMatches(c.client, c.db, matches, totalMatches)
				if err != nil {
					log.Printf("process matches: %v", err)
					continue
				}

				lastID = msg.ID
			}

			time.Sleep(15 * time.Second)
		}
	}
}