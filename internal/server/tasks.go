package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zollidan/esmeralda/internal/models"
	"github.com/zollidan/esmeralda/internal/queue"
	"gorm.io/gorm"
)

type createTaskRequest struct {
	Date string `json:"date" binding:"required"`
}

func (h *Handler) CreateTask(c *gin.Context) {
	var req createTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "поле date обязательно"})
		return
	}

	if _, err := time.Parse("2006-01-02", req.Date); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "формат даты: YYYY-MM-DD"})
		return
	}

	task := queue.ParseTask{
		ID:        uuid.New().String(),
		Date:      req.Date,
		CreatedAt: time.Now().UTC(),
	}

	msgID, err := h.producer.Publish(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось добавить задачу в очередь"})
		return
	}

	record := &models.Task{
		ID:        task.ID,
		Date:      task.Date,
		Status:    string(queue.StatusPending),
		CreatedAt: task.CreatedAt,
	}

	if err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
		return tx.Create(record).Error
	}); err != nil {
		_ = h.producer.Delete(c.Request.Context(), msgID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось сохранить задачу в БД"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *Handler) GetTasks(c *gin.Context) {
	var tasks []models.Task
	if err := h.db.Order("created_at desc").Find(&tasks).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить задачи"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}
