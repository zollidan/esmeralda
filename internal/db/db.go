package db

import (
	"fmt"

	"github.com/zollidan/esmeralda/internal/config"
	"github.com/zollidan/esmeralda/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(cfg *config.Config) (*gorm.DB, error) {
	database, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	// Drop FK constraint if it exists (games.task_id no longer references tasks)
	// удалить после миграции, когда FK будет удален из модели Game
	database.Exec("ALTER TABLE games DROP CONSTRAINT IF EXISTS fk_games_task")

	if err := database.AutoMigrate(&models.Task{}, &models.Game{}); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	return database, nil
}
