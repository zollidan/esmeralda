package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"github.com/zollidan/esmeralda/internal/models"
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
	defer f.Close()

	sheet := "Games"
	f.SetSheetName("Sheet1", sheet)

	headers := stats.Headers()
	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheet, cell, header)
	}

	for rowIdx, game := range games {
		row := gameToRow(game)
		values := row.ToSlice()
		for col, val := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, rowIdx+2)
			f.SetCellValue(sheet, cell, val)
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

func gameToRow(g models.Game) stats.Row {
	return stats.Row{
		Day:      g.Day,
		Month:    time.Month(g.Month),
		Year:     g.Year,
		Time:     g.Time,
		HomeTeam: g.HomeTeam,
		AwayTeam: g.AwayTeam,
		League:   g.League,

		H2H15:     stats.TotalStats(g.H2H15),
		H2H25:     stats.TotalStats(g.H2H25),
		Home25:    stats.TotalStats(g.Home25),
		Away15:    stats.TotalStats(g.Away15),
		Away25:    stats.TotalStats(g.Away25),
		AwayK2_25: stats.TotalStats(g.AwayK2_25),

		GoalsH2H:  stats.GoalStats(g.GoalsH2H),
		GoalsHome: stats.GoalStats(g.GoalsHome),
		GoalsAway: stats.GoalStats(g.GoalsAway),
	}
}
