package db

import (
	"fmt"

	"github.com/zollidan/esmeralda-ru-api-fetcher/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func New(dsn string) (*gorm.DB, error) {
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	if err := database.AutoMigrate(&models.Task{}); err != nil {
		return nil, fmt.Errorf("auto migrate: %w", err)
	}

	return database, nil
}
