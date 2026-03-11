package server

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	defer ws.Close()

	ch := h.subscribeProgress(taskID)
	defer h.unsubscribeProgress(taskID, ch)

	// Close channel when client disconnects
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			if _, _, err := ws.ReadMessage(); err != nil {
				return
			}
		}
	}()

	for {
		select {
		case p, ok := <-ch:
			if !ok {
				return
			}
			if err := ws.WriteJSON(p); err != nil {
				log.Printf("websocket write: %v", err)
				return
			}
		case <-done:
			return
		}
	}
}
