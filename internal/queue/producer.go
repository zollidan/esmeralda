package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusDone       Status = "done"
	StatusError      Status = "error"
	StatusFailed     Status = "failed"
	StatusCancelled   Status = "cancelled"

)

const (
	StreamName = "tasks:parse"
	GroupName  = "parse-workers" 
)

type ParseTask struct {
	ID        string    `json:"id"`
	Date      string    `json:"date"` // "2006-01-02"
	CreatedAt time.Time `json:"created_at"`
	Status    Status    `json:"status"`
}

type Producer struct {
	rdb *redis.Client
}

func NewProducer(rdb *redis.Client) *Producer {
	return &Producer{rdb: rdb}
}

func (p *Producer) Enqueue(ctx context.Context, date string) (*ParseTask, string, error) {
	task := &ParseTask{
		ID:        uuid.New().String(),
		Date:      date,
		CreatedAt: time.Now().UTC(), 
		Status:    "pending",
	}

	data, err := json.Marshal(task)
	if err != nil {
		return nil, "", fmt.Errorf("marshal task: %w", err)
	}

	res, err := p.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: StreamName,
		Values: map[string]interface{}{
			"payload": string(data),
		},
	}).Result()
	if err != nil {
		return nil, "", fmt.Errorf("xadd failed: %w", err)
	}

	return task, res, nil
}