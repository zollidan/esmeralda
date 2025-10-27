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

func (h *Handlers) FilesRoutes(r chi.Router) {
	r.Get("/", h.GetFiles)
	r.Post("/", h.CreateFile)
}

func (h *Handlers) GetFiles(w http.ResponseWriter, r *http.Request) {
	files, err := gorm.G[models.Files](h.DB).Find(context.Background())
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error reading files.",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func (h *Handlers) CreateFile(w http.ResponseWriter, r *http.Request) {
	var req schemas.CreateFileRequest

	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	err := gorm.G[models.Files](h.DB).Create(context.Background(), &models.Files{Filename: req.Filename, FileURL: req.FileURL})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create file: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)
}
