package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zollidan/esmeralda/internal/queue"
	"github.com/zollidan/esmeralda/internal/repository"
)

type Handler struct {
	producer    *queue.Producer
	rdb         *redis.Client
	tasks       *repository.TaskRepository
	games       *repository.GameRepository
	pending     map[string]chan *queue.MatchDataResult
	progressSub map[string]map[chan queue.TaskProgress]struct{}
	mu          sync.Mutex
}

func NewHandler(producer *queue.Producer, rdb *redis.Client, tasks *repository.TaskRepository, games *repository.GameRepository) *Handler {
	return &Handler{
		producer:    producer,
		rdb:         rdb,
		tasks:       tasks,
		games:       games,
		pending:     make(map[string]chan *queue.MatchDataResult),
		progressSub: make(map[string]map[chan queue.TaskProgress]struct{}),
	}
}

func (h *Handler) subscribeProgress(taskID string) chan queue.TaskProgress {
	h.mu.Lock()
	defer h.mu.Unlock()
	ch := make(chan queue.TaskProgress, 64)
	if h.progressSub[taskID] == nil {
		h.progressSub[taskID] = make(map[chan queue.TaskProgress]struct{})
	}
	h.progressSub[taskID][ch] = struct{}{}
	return ch
}

func (h *Handler) unsubscribeProgress(taskID string, ch chan queue.TaskProgress) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.progressSub[taskID], ch)
	if len(h.progressSub[taskID]) == 0 {
		delete(h.progressSub, taskID)
	}
	close(ch)
}

func (h *Handler) broadcastProgress(p queue.TaskProgress) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for ch := range h.progressSub[p.TaskID] {
		select {
		case ch <- p:
		default:
		}
	}
}

func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
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

	go func() {
		c := queue.NewConsumer(h.rdb, queue.StreamProgress, "progress_group", "progress_consumer")
		err := c.Consume(ctx, func(ctx context.Context, payload []byte) error {
			progress, err := queue.Unmarshal[queue.TaskProgress](payload)
			if err != nil {
				return err
			}
			h.broadcastProgress(progress)
			return nil
		})
		if err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("progress consumer error: %v", err)
		}
	}()
}
