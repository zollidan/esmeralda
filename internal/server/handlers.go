package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/zollidan/esmeralda/internal/config"
	"github.com/zollidan/esmeralda/internal/queue"
	"github.com/zollidan/esmeralda/internal/repository"
	"github.com/zollidan/esmeralda/internal/sse"
)

type Handler struct {
	producer      *queue.Producer
	cfg           *config.Config
	rdb           *redis.Client
	tasks         *repository.TaskRepository
	games         *repository.GameRepository
	users         *repository.UserRepository
	refreshTokens *repository.RefreshTokenRepository
	pending       map[string]chan *queue.MatchDataResult
	hub           *sse.Hub
}

func NewHandler(producer *queue.Producer, cfg *config.Config, rdb *redis.Client, tasks *repository.TaskRepository, games *repository.GameRepository, users *repository.UserRepository, refreshTokens *repository.RefreshTokenRepository, hub *sse.Hub) *Handler {
	return &Handler{
		producer:      producer,
		cfg:           cfg,
		rdb:           rdb,
		tasks:         tasks,
		games:         games,
		users:         users,
		refreshTokens: refreshTokens,
		hub:           hub,
		pending:       make(map[string]chan *queue.MatchDataResult),
	}
}

// Health godoc
// @Summary      Проверка доступности сервиса
// @Description  Возвращает статус API
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

type SportRUHealthResponse struct {
    Status  string `json:"status"`          
    Latency int64  `json:"latency_ms"`      
    Error   string `json:"error,omitempty"` 
}

// @Summary Health check внешнего SportRU API
// @Tags sportru
// @Produce json
// @Success 200 {object} SportRUHealthResponse
// @Router /api/sportru/health [get]
func (h *Handler) SportRUAPIHealth(c *gin.Context) {
    apiKey := h.cfg.SportAPIRU.Token
    wsURL := "wss://ws.api.api-sport.ru?apiKey=" + apiKey

    start := time.Now()

    dialer := websocket.Dialer{
        HandshakeTimeout: 5 * time.Second,
    }

    conn, _, err := dialer.Dial(wsURL, nil)
    if err != nil {
        c.JSON(http.StatusOK, SportRUHealthResponse{
            Status:  "offline",
            Latency: time.Since(start).Milliseconds(),
            Error:   err.Error(),
        })
        return
    }
    defer conn.Close()

    // Сервер сразу шлёт {"type":"connected"} — этого достаточно
    conn.SetReadDeadline(time.Now().Add(5 * time.Second))

    var msg struct {
        Type string `json:"type"`
    }
    if err := conn.ReadJSON(&msg); err != nil {
        c.JSON(http.StatusOK, SportRUHealthResponse{
            Status:  "offline",
            Latency: time.Since(start).Milliseconds(),
            Error:   "no connected message: " + err.Error(),
        })
        return
    }

    if msg.Type != "connected" {
        c.JSON(http.StatusOK, SportRUHealthResponse{
            Status:  "offline",
            Latency: time.Since(start).Milliseconds(),
            Error:   "unexpected: " + msg.Type,
        })
        return
    }

    c.JSON(http.StatusOK, SportRUHealthResponse{
        Status:  "online",
        Latency: time.Since(start).Milliseconds(),
    })
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
