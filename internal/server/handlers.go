package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zollidan/esmeralda/internal/queue"
)

type Handler struct {
	producer *queue.Producer
}

func NewHandler(producer *queue.Producer) *Handler {
	return &Handler{producer: producer}
}

type createTaskRequest struct {
	Date string `json:"date" binding:"required"`
}

func (h *Handler) CreateTask(c *gin.Context) {
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "поле date обязательно"})
		return
	}

	// валидация формата даты
	if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "формат даты: YYYY-MM-DD"})
		return
	}

	task, err := h.producer.Enqueue(c.Request.Context(), req.Date)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось добавить задачу"})
		return
	}

	c.JSON(http.StatusCreated, task)
}