package handlers

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis"
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
	// database.Get[models.Task](w, r, h.DB)

	data, err := h.RedisClient.FindResults()
	if err != nil {
		if err == redis.Nil {
			utils.ResponseJSON(w, http.StatusNotFound, map[string]string{
				"error": err.Error(),
			})
			return
		}
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	utils.ResponseJSON(w, http.StatusOK, data)
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	var req schemas.CreateTaskRequest
	defer r.Body.Close()

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	parserAvailable, err := h.RedisClient.GetParsers()
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	if !slices.Contains(parserAvailable, req.ParserName) {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid parser name",
		})
		return
	}

	task := &models.Task{
		Name:        req.Name,
		Description: req.Description,
		ParserName:  req.ParserName,
	}

	err = database.Create(w, r, h.DB, task)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	streamID, err := h.RedisClient.SendTask(task.ID, req.ParserName, req.FileName)
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
