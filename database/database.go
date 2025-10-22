package database

import (
	"github.com/zollidan/esmeralda/config"
	"github.com/zollidan/esmeralda/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitDatabase(cfg config.Config) *gorm.DB {

	db, err := gorm.Open(sqlite.Open(cfg.DBURL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Files{}, &models.Tasks{})

	return db
}
