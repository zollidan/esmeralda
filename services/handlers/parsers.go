package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zollidan/esmeralda/utils"
)

type Parser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handlers) ParsersRoutes(r chi.Router) {
	r.Get("/", h.GetParsers)
}

func (h *Handlers) GetParsers(w http.ResponseWriter, r *http.Request) {
	data, err := h.RedisClient.GetParsers()
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Failed to get parsers from Redis",
		})
		return
	}

	utils.ResponseJSON(w, http.StatusOK, data)

}
