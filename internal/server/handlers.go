package server

import (
	"context"
	"errors"
	"log"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/zollidan/esmeralda/internal/queue"
	"github.com/zollidan/esmeralda/internal/repository"
)

type Handler struct {
	producer *queue.Producer
	rdb      *redis.Client
	tasks    *repository.TaskRepository
	games    *repository.GameRepository
	pending  map[string]chan *queue.MatchDataResult
	mu       sync.Mutex
}

func NewHandler(producer *queue.Producer, rdb *redis.Client, tasks *repository.TaskRepository, games *repository.GameRepository) *Handler {
	return &Handler{
		producer: producer,
		rdb:      rdb,
		tasks:    tasks,
		games:    games,
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
			return h.tasks.UpdateStatus(ctx, result.TaskID, string(result.Status))
		})
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("results consumer error: %v", err)
		}
	}()
}
