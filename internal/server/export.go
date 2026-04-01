package server

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"github.com/zollidan/esmeralda/internal/models"
	"github.com/zollidan/esmeralda/internal/stats"
)

// ExportGames godoc
// @Summary      Экспортировать матчи
// @Description  Формирует файл по матчам за период date_start - date_end. Формат: xlsx (по умолчанию) или csv.
// @Tags         export
// @Produce      application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Produce      text/csv
// @Param        date_start  query     string  true   "Дата начала периода (YYYY-MM-DD)"
// @Param        date_end    query     string  true   "Дата конца периода (YYYY-MM-DD)"
// @Param        format      query     string  false  "Формат файла: xlsx или csv (по умолчанию xlsx)"
// @Success      200         {file}    file
// @Failure      400         {object}  errorResponse
// @Failure      404         {object}  errorResponse
// @Failure      500         {object}  errorResponse
// @Security     BearerAuth
// @Router       /export/ [get]
func (h *Handler) ExportGames(c *gin.Context) {
	dateStartStr := c.Query("date_start")
	dateEndStr := c.Query("date_end")
	format := c.DefaultQuery("format", "xlsx")

	if dateStartStr == "" || dateEndStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "параметры date_start и date_end обязательны (формат: YYYY-MM-DD)"})
		return
	}

	if format != "xlsx" && format != "csv" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "параметр format должен быть xlsx или csv"})
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

	switch format {
	case "csv":
		filename := fmt.Sprintf("esmeralda_%s_%s.csv", dateStartStr, dateEndStr)
		c.Header("Content-Type", "text/csv; charset=utf-8")
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		if err := writeCSV(c.Writer, games); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось сформировать CSV файл"})
		}
	default:
		filename := fmt.Sprintf("esmeralda_%s_%s.xlsx", dateStartStr, dateEndStr)
		c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		if err := writeXLSX(c.Writer, games); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось сформировать Excel файл"})
		}
	}
}

func writeXLSX(w http.ResponseWriter, games []models.Game) error {
	f := excelize.NewFile()
	defer func() {
		_ = f.Close()
	}()

	sheet := "Games"
	if err := f.SetSheetName("Sheet1", sheet); err != nil {
		return err
	}

	headers := stats.Headers()
	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		if err := f.SetCellValue(sheet, cell, header); err != nil {
			return err
		}
	}

	for rowIdx, game := range games {
		values := stats.GameToSlice(&game)
		for col, val := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, rowIdx+2)
			if err := f.SetCellValue(sheet, cell, val); err != nil {
				return err
			}
		}
	}

	return f.Write(w)
}

func writeCSV(w http.ResponseWriter, games []models.Game) error {
	cw := csv.NewWriter(w)

	if err := cw.Write(stats.Headers()); err != nil {
		return err
	}

	for _, game := range games {
		values := stats.GameToSlice(&game)
		row := make([]string, len(values))
		for i, v := range values {
			row[i] = fmt.Sprintf("%v", v)
		}
		if err := cw.Write(row); err != nil {
			return err
		}
	}

	cw.Flush()
	return cw.Error()
}
