package database

import (
	"context"
	"net/http"

	"github.com/zollidan/esmeralda/config"
	"github.com/zollidan/esmeralda/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Model interface {
	models.Files | models.Task
}

func Get[T Model](w http.ResponseWriter, r *http.Request, db *gorm.DB) (*[]T, error) {
	items, err := gorm.G[T](db).Find(context.Background())
	if err != nil {
		return nil, err
	}
	return &items, nil
}

func Create[T Model](w http.ResponseWriter, r *http.Request, db *gorm.DB, item *T) error {
	err := gorm.G[T](db).Create(context.Background(), item)
	if err != nil {
		return err
	}
	return nil
}

func GetBy[T Model](db *gorm.DB, id string, searchField string) (*T, error) {
	result, err := gorm.G[T](db).Where(searchField+" = ?", id).First(context.Background())
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func Delete[T Model](db *gorm.DB, id string) error {
	_, err := gorm.G[T](db).Where("id = ?", id).Delete(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func Init(cfg *config.Config) *gorm.DB {

	db, err := gorm.Open(postgres.Open(cfg.DBURL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Files{}, &models.Task{})

	return db
}
