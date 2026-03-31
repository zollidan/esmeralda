package server

import (
	"encoding/json"
	"errors"
	"fmt"
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

//nolint:unused // used in swagger comments
type messageResponse struct {
	Message string `json:"message"`
}

//nolint:unused // used in swagger comments
type errorResponse struct {
	Error string `json:"error"`
}

// CreateTask godoc
// @Summary      Создать задачу парсинга
// @Description  Создает новую задачу на парсинг матчей за указанную дату
// @Tags         tasks
// @Accept       json
// @Produce      json
// @Param        request body      createTaskRequest true "Дата задачи"
// @Success      201     {object}  models.Task
// @Failure      400     {object}  errorResponse
// @Failure      500     {object}  errorResponse
// @Security     BearerAuth
// @Router       /tasks/ [post]
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

	if err := h.tasks.Create(c.Request.Context(), record); err != nil {
		_ = h.producer.Delete(c.Request.Context(), msgID)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось сохранить задачу в БД"})
		return
	}

	h.hub.Broadcast(record)
	c.JSON(http.StatusCreated, record)
}

// GetTasks godoc
// @Summary      Получить список задач
// @Description  Возвращает все задачи, отсортированные по времени создания
// @Tags         tasks
// @Produce      json
// @Success      200  {array}   models.Task
// @Failure      500  {object}  errorResponse
// @Security     BearerAuth
// @Router       /tasks/ [get]
func (h *Handler) GetTasks(c *gin.Context) {
	tasks, err := h.tasks.GetAllOrdered(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить задачи"})
		return
	}
	c.JSON(http.StatusOK, tasks)
}

// GetTask godoc
// @Summary      Получить задачу по ID
// @Description  Возвращает одну задачу по идентификатору
// @Tags         tasks
// @Produce      json
// @Param        id   path      string  true  "ID задачи"
// @Success      200  {object}  models.Task
// @Failure      400  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Security     BearerAuth
// @Router       /tasks/{id} [get]
func (h *Handler) GetTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id задачи обязателен"})
		return
	}

	task, err := h.tasks.FindByID(c.Request.Context(), taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "задача не найдена"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить задачу"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// DeleteTask godoc
// @Summary      Удалить задачу
// @Description  Удаляет задачу по ID
// @Tags         tasks
// @Produce      json
// @Param        id   path      string  true  "ID задачи"
// @Success      200  {object}  messageResponse
// @Failure      400  {object}  errorResponse
// @Failure      404  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Security     BearerAuth
// @Router       /tasks/{id} [delete]
func (h *Handler) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id задачи обязателен"})
		return
	}
	task, err := h.tasks.FindByID(c.Request.Context(), taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "задача не найдена"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить задачу"})
		return
	}

	err = h.tasks.Delete(c.Request.Context(), task.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось удалить задачу"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "задача успешно удалена"})
}

// StreamTasks godoc
// @Summary      Поток обновлений задач (SSE)
// @Description  Открывает Server-Sent Events поток с актуальным списком и изменениями задач
// @Tags         tasks
// @Produce      text/event-stream
// @Success      200  {string}  string  "SSE stream"
// @Security     BearerAuth
// @Router       /tasks/stream [get]
func (h *Handler) StreamTasks(c *gin.Context) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	client := h.hub.Subscribe()

	defer h.hub.Unsubscribe(client)

	tasks, err := h.tasks.GetAllOrdered(c.Request.Context())
	if err == nil {
		data, _ := json.Marshal(tasks)
		_, _ = fmt.Fprintf(c.Writer, "data: %s\n\n", data) //nolint:gosec // data is json.Marshal output
		c.Writer.Flush()
	}

	for {
		select {
		case msg, ok := <-client:
			if !ok {
				return
			}
			_, _ = fmt.Fprintf(c.Writer, "data: %s\n\n", msg) //nolint:gosec // data is json.Marshal output
			c.Writer.Flush()
		case <-c.Request.Context().Done():
			return
		}
	}
}
