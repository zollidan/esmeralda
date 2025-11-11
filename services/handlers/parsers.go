package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis"
)

const parsersRedisKey = "parsers"

type Parser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handlers) ParsersRoutes(r chi.Router) {
	r.Get("/", h.GetParsers)
	r.Post("/refresh", h.RefreshParsers)
}

func (h *Handlers) RefreshParsers(w http.ResponseWriter, r *http.Request) {
	// make rabbit request to piko service to refresh parsers
}

func (h *Handlers) GetParsers(w http.ResponseWriter, r *http.Request) {
	data, err := h.RedisClient.Get(parsersRedisKey).Result()
	if err != nil {
		if err == redis.Nil {
			http.Error(w, "Parsers not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Failed to get parsers from Redis", http.StatusInternalServerError)
		return
	}

	var parsers []Parser
	if err := json.Unmarshal([]byte(data), &parsers); err != nil {
		http.Error(w, "Failed to unmarshal parsers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(parsers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
