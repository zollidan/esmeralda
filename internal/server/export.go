package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"github.com/zollidan/esmeralda/internal/stats"
)

func (h *Handler) ExportGames(c *gin.Context) {
	dateStartStr := c.Query("date_start")
	dateEndStr := c.Query("date_end")

	if dateStartStr == "" || dateEndStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "параметры date_start и date_end обязательны (формат: YYYY-MM-DD)"})
		return
	}

	dateStart, err := time.Parse("2006-01-02", dateStartStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "формат date_start: YYYY-MM-DD"})
		return
	}

	dateEnd, err := time.Parse("2006-01-02", dateEndStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "формат date_end: YYYY-MM-DD"})
		return
	}

	if dateEnd.Before(dateStart) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date_end не может быть раньше date_start"})
		return
	}

	games, err := h.games.FindByDateRange(
		c.Request.Context(),
		dateStart.Year(), int(dateStart.Month()), dateStart.Day(),
		dateEnd.Year(), int(dateEnd.Month()), dateEnd.Day(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось получить данные из БД"})
		return
	}

	if len(games) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "нет данных за указанный период"})
		return
	}

	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	sheet := "Games"
	if err := f.SetSheetName("Sheet1", sheet); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось подготовить лист Excel"})
		return
	}

	headers := stats.Headers()
	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		if err := f.SetCellValue(sheet, cell, header); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось записать заголовки Excel"})
			return
		}
	}

	for rowIdx, game := range games {
		values := stats.GameToSlice(&game)
		for col, val := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, rowIdx+2)
			if err := f.SetCellValue(sheet, cell, val); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось записать данные Excel"})
				return
			}
		}
	}

	filename := fmt.Sprintf("esmeralda_%s_%s.xlsx", dateStartStr, dateEndStr)
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

	if err := f.Write(c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось сформировать Excel файл"})
		return
	}
}
