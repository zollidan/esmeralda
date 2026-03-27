package sse

import (
	"encoding/json"
	"sync"
)

type Client chan []byte

type TaskProgress struct {
	TaskID  string  `json:"task_id"`
	Percent float64 `json:"percent"`
	Message string  `json:"message"`
}

type Hub struct {
	mu      sync.RWMutex
	clients map[Client]struct{}
}

func NewHub() *Hub {
	return &Hub{clients: make(map[Client]struct{})}
}

func (h *Hub) Subscribe() Client {
	ch := make(Client, 8)
	h.mu.Lock()
	h.clients[ch] = struct{}{}
	h.mu.Unlock()
	return ch
}

func (h *Hub) Unsubscribe(ch Client) {
	h.mu.Lock()
	delete(h.clients, ch)
	close(ch)
	h.mu.Unlock()
}

func (h *Hub) Broadcast(payload any) {
	data, _ := json.Marshal(payload)
	h.mu.RLock()
	defer h.mu.RUnlock()
	for ch := range h.clients {
		select {
		case ch <- data:
		default:
		}
	}
}
