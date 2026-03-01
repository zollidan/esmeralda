package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const TasksQueue = "tasks:parse"

type ParseTask struct {
	ID        string    `json:"id"`
	Date      string    `json:"date"` // "2006-01-02"
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"` // pending/processing/done/error
}

type Producer struct {
	rdb *redis.Client
}

func NewProducer(rdb *redis.Client) *Producer {
	return &Producer{rdb: rdb}
}

func (p *Producer) Enqueue(ctx context.Context, date string) (*ParseTask, error) {
	task := &ParseTask{
		ID:        uuid.New().String(),
		Date:      date,
		CreatedAt: time.Now(),
		Status:    "pending",
	}

	data, err := json.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("marshal task: %w", err)
	}

	if err := p.rdb.LPush(ctx, TasksQueue, data).Err(); err != nil {
		return nil, fmt.Errorf("lpush: %w", err)
	}

	return task, nil
}