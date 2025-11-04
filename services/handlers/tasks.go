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

func (h *Handlers) TasksRoutes(r chi.Router) {
	r.Get("/", h.GetTasks)
	r.Post("/", h.CreateTask)
}

func (h *Handlers) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := gorm.G[models.Task](h.DB).Preload("Parser", nil).Find(context.Background())
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error reading tasks.",
		})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req schemas.CreateTaskRequest
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Find parser to ensure it exists
	_, err := gorm.G[models.Parser](h.DB).Where("id = ?", req.ParserID).First(context.Background())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			http.Error(w, "parser not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("failed to find parser: %v", err), http.StatusInternalServerError)
		}
		return
	}

	task := &models.Task{
		Name:        req.Name,
		Description: req.Description,
		ParserID:    req.ParserID,
	}

	err = h.TaskManager.StartTask(task.ID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to start task: %v", err), http.StatusInternalServerError)
		return
	}

	err = gorm.G[models.Task](h.DB).Create(context.Background(), task)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create task in database: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"task_id": task.ID,
	})

}