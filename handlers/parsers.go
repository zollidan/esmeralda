package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zollidan/esmeralda/models"
	"github.com/zollidan/esmeralda/schemas"
	"gorm.io/gorm"
)

func (h *Handlers) ParsersRoutes(r chi.Router) {
	r.Get("/", h.GetParsers)
	r.Post("/", h.CreateParser)
}

func (h *Handlers) GetParsers(w http.ResponseWriter, r *http.Request) {
	parsers, err := gorm.G[models.Parser](h.DB).Find(context.Background())
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error reading parsers.",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(parsers)
}

func (h *Handlers) CreateParser(w http.ResponseWriter, r *http.Request) {
	var req schemas.CreateParserRequest

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	err := gorm.G[models.Parser](h.DB).Create(context.Background(), &models.Parser{Name: req.ParserName, Description: req.ParserDescription})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create parser: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}
