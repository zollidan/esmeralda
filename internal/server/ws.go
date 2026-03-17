package server

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/zollidan/esmeralda/internal/queue"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) WSHandler(c *gin.Context) {
	taskID := c.Query("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task_id is required"})
		return
	}

	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket upgrade: %v", err)
		return
	}
	defer func() {
		_ = ws.Close()
	}()

	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()

	// горутина для детекта дисконнекта клиента
	go func() {
		for {
			if _, _, err := ws.ReadMessage(); err != nil {
				cancel()
				return
			}
		}
	}()

	lastID := "0-0"

	for {
		streams, err := h.rdb.XRead(ctx, &redis.XReadArgs{
			Streams: []string{queue.StreamProgress, lastID},
			Count:   100,
			Block:   2 * time.Second,
		}).Result()

		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			if errors.Is(err, redis.Nil) {
				continue
			}
			log.Printf("xread error: %v", err)
			return
		}

		for _, stream := range streams {
			for _, msg := range stream.Messages {
				lastID = msg.ID

				payload, ok := msg.Values["payload"].(string)
				if !ok {
					continue
				}

				progress, err := queue.Unmarshal[queue.TaskProgress]([]byte(payload))
				if err != nil {
					continue
				}

				if progress.TaskID != taskID {
					continue
				}

				if err := ws.WriteJSON(progress); err != nil {
					log.Printf("websocket write: %v", err)
					return
				}

				if isTerminal(string(progress.Status)) {
					return
				}
			}
		}
	}
}

func isTerminal(status string) bool {
	switch status {
	case "done", "error", "failed", "cancelled":
		return true
	}
	return false
}
