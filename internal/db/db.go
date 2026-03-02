package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/zollidan/esmeralda/internal/models"
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
