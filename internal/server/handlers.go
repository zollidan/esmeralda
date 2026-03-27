package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zollidan/esmeralda/internal/queue"
	"github.com/zollidan/esmeralda/internal/repository"
	"github.com/zollidan/esmeralda/internal/sse"
)

type Handler struct {
	producer *queue.Producer
	rdb      *redis.Client
	tasks    *repository.TaskRepository
	games    *repository.GameRepository
	pending  map[string]chan *queue.MatchDataResult
	hub      *sse.Hub
}

func NewHandler(producer *queue.Producer, rdb *redis.Client, tasks *repository.TaskRepository, games *repository.GameRepository, hub *sse.Hub) *Handler {
	return &Handler{
		producer: producer,
		rdb:      rdb,
		tasks:    tasks,
		games:    games,
		hub:      hub,
		pending:  make(map[string]chan *queue.MatchDataResult),
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

			if err := h.tasks.UpdateStatus(ctx, result.TaskID, string(result.Status)); err != nil {
				return err
			}

			updatedTask, err := h.tasks.FindByID(ctx, result.TaskID)
			if err == nil {
				h.hub.Broadcast(updatedTask)
			}

			return nil
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

			percent := 0
			if progress.TotalMatches > 0 {
				percent = progress.CurrentMatch * 100 / progress.TotalMatches
			}

			h.hub.Broadcast(sse.TaskProgress{
				TaskID:  progress.TaskID,
				Percent: float64(percent),
				Message: fmt.Sprintf("Матч %d из %d: ", progress.CurrentMatch, progress.TotalMatches),
			})

			return nil
		})

		if err != nil && !errors.Is(err, context.Canceled) {
			log.Printf("progress consumer error: %v", err)
		}
	}()
}
