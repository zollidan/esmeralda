package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis"
	"github.com/zollidan/esmeralda/utils"
)

type Parser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handlers) ParsersRoutes(r chi.Router) {
	r.Get("/", h.GetParsers)
	r.Get("/{id}", h.GetParserByID)
}

func (h *Handlers) GetParserByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	data, err := h.RedisClient.Client.Get("parsers:" + id).Result()
	if err != nil {
		if err == redis.Nil {
			utils.ResponseJSON(w, http.StatusNotFound, map[string]string{
				"error": "Parser not found",
			})
			return
		}
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to get parser from Redis",
		})
		return
	}

	var parser Parser
	if err := json.Unmarshal([]byte(data), &parser); err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to unmarshal parser",
		})
		return
	}

	utils.ResponseJSON(w, http.StatusOK, parser)
}

func (h *Handlers) GetParsers(w http.ResponseWriter, r *http.Request) {
	data, err := h.RedisClient.Client.Get("parsers").Result()
	if err != nil {
		if err == redis.Nil {
			utils.ResponseJSON(w, http.StatusNotFound, map[string]string{
				"error": "Parsers not found",
			})
			return
		}
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to get parsers from Redis",
		})
		return
	}

	var parsers []Parser
	if err := json.Unmarshal([]byte(data), &parsers); err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to unmarshal parsers",
		})
		return
	}

	utils.ResponseJSON(w, http.StatusOK, parsers)
}
