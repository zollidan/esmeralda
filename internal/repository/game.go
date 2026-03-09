package repository

import (
	"context"
	"fmt"

	"github.com/zollidan/esmeralda/internal/models"
	"gorm.io/gorm"
)

type GameRepository struct {
	*Repository[models.Game]
}

func NewGameRepository(db *gorm.DB) *GameRepository {
	return &GameRepository{Repository: New[models.Game](db)}
}

// FindByDateRange returns games within [startDay/startMonth/startYear .. endDay/endMonth/endYear].
func (r *GameRepository) FindByDateRange(ctx context.Context, startYear, startMonth, startDay, endYear, endMonth, endDay int) ([]models.Game, error) {
	var games []models.Game

	// Build a comparable integer date: year*10000 + month*100 + day
	err := r.DB().WithContext(ctx).
		Where("(year * 10000 + month * 100 + day) >= ? AND (year * 10000 + month * 100 + day) <= ?",
			startYear*10000+startMonth*100+startDay,
			endYear*10000+endMonth*100+endDay,
		).
		Order("year, month, day, time").
		Find(&games).Error
	if err != nil {
		return nil, fmt.Errorf("find games by date range: %w", err)
	}
	return games, nil
}

func (r *GameRepository) FindByTaskID(ctx context.Context, taskID string) ([]models.Game, error) {
	var games []models.Game
	if err := r.DB().WithContext(ctx).Where("task_id = ?", taskID).Find(&games).Error; err != nil {
		return nil, fmt.Errorf("find games by task id: %w", err)
	}
	return games, nil
}
