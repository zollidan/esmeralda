package server

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/zollidan/esmeralda/internal/models"
	"github.com/zollidan/esmeralda/internal/queue"
	"gorm.io/gorm"
)

type Handler struct {
	producer *queue.Producer
	db       *gorm.DB
	rdb      *redis.Client
	pending  map[string]chan *queue.MatchDataResult
	mu       sync.Mutex
}

func NewHandler(producer *queue.Producer, rdb *redis.Client, db *gorm.DB) *Handler {
	return &Handler{
		producer: producer,
		db:       db,
		rdb:      rdb,
		pending:  make(map[string]chan *queue.MatchDataResult),
	}
}

func (h *Handler) StartConsumers(ctx context.Context) {
	go func() {
		c := queue.NewConsumer(h.rdb, queue.StreamResults, "results_group", "results_consumer")
		err := c.Consume(ctx, func(ctx context.Context, payload []byte) error {
			result, err := queue.Unmarshal[queue.TaskResult](payload)
			if err != nil {
				return err
			}
			return h.db.WithContext(ctx).
				Model(&models.Task{}).
				Where("id = ?", result.TaskID).
				Update("status", string(result.Status)).Error
		})
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("results consumer error: %v", err)
		}
	}()
}
