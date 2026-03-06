package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zollidan/esmeralda/internal/queue"
)

func (h *Handler) GetExcel(c *gin.Context) {
	dateStart := c.Query("date_start")
	dateEnd := c.Query("date_end")

	if dateStart == "" || dateEnd == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "параметры date_start и date_end обязательны"})
		return
	}

	if _, err := time.Parse("2006-01-02", dateStart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "формат date_start: YYYY-MM-DD"})
		return
	}

	if _, err := time.Parse("2006-01-02", dateEnd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "формат date_end: YYYY-MM-DD"})
		return
	}

	task := queue.MatchDataTask{
		ID:        uuid.New().String(),
		DateStart: dateStart,
		DateEnd:   dateEnd,
		CreatedAt: time.Now().UTC(),
	}

	ch := make(chan *queue.MatchDataResult, 1)
	h.mu.Lock()
	h.pending[task.ID] = ch
	h.mu.Unlock()
	defer func() {
		h.mu.Lock()
		delete(h.pending, task.ID)
		h.mu.Unlock()
	}()

	if _, err := h.enrichProducer.Publish(c.Request.Context(), task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось отправить запрос"})
		return
	}

	select {
	case result := <-ch:
		c.JSON(http.StatusOK, result)
	case <-time.After(10 * time.Second):
		c.JSON(http.StatusGatewayTimeout, gin.H{"error": "таймаут ожидания данных"})
	case <-c.Request.Context().Done():
	}
}
