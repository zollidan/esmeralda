package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zollidan/esmeralda/internal/models"
	"github.com/zollidan/esmeralda/internal/queue"

	"gorm.io/gorm"
)

type Handler struct {
	producer *queue.Producer
	db       *gorm.DB
}

func NewHandler(producer *queue.Producer, db *gorm.DB) *Handler {
	return &Handler{producer: producer, db: db}
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

    if _, err := time.Parse("2006-01-02", req.Date); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "формат даты: YYYY-MM-DD"})
        return
    }

    task, streamID, err := h.producer.Enqueue(c.Request.Context(), req.Date)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось добавить задачу в очередь"})
        return
    }

    record := &models.Task{
        ID:        task.ID,
        Date:      task.Date,
        Status:    string(task.Status),
        CreatedAt: task.CreatedAt,
    }

    if err := h.db.WithContext(c.Request.Context()).Transaction(func(tx *gorm.DB) error {
        return tx.Create(record).Error
    }); err != nil {
        _ = h.producer.Delete(c.Request.Context(), streamID)
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