package db

import (
	"fmt"

	"github.com/zollidan/esmeralda/internal/config"
	"github.com/zollidan/esmeralda/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New(cfg *config.Config) (*gorm.DB, error) {
	database, err := gorm.Open(sqlite.Open(cfg.Database.SQLiteURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := database.AutoMigrate(&models.Task{}, &models.Game{}); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	return database, nil
}
