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

	if err := database.AutoMigrate(&models.Task{}, &models.Game{}, &models.User{}, &models.RefreshToken{}); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	return database, nil
}
