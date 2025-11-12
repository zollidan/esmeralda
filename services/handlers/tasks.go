package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/zollidan/esmeralda/database"
	"github.com/zollidan/esmeralda/models"
	"github.com/zollidan/esmeralda/schemas"
	"github.com/zollidan/esmeralda/utils"
)

func (h *Handlers) TasksRoutes(r chi.Router) {
	r.Get("/", h.GetTasks)
	r.Post("/", h.CreateTask)
}

func (h *Handlers) GetTasks(w http.ResponseWriter, r *http.Request) {
	database.Get[models.Task](w, r, h.DB)
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req schemas.CreateTaskRequest
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	task := &models.Task{
		Name:        req.Name,
		Description: req.Description,
	}

	err := database.Create(w, r, h.DB, task)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	streamID, err := h.RedisClient.SendTask(task.ID)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, map[string]interface{}{
		"task_id":   task.ID,
		"stream_id": streamID,
	})

}
