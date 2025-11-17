package database

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/zollidan/esmeralda/config"
	"github.com/zollidan/esmeralda/models"
	"github.com/zollidan/esmeralda/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Model interface {
	models.Files | models.Task
}

func Get[T Model](w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	items, err := gorm.G[T](db).Find(context.Background())
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Error reading data.",
		})
		return
	}
	utils.ResponseJSON(w, http.StatusOK, items)
}

func Create[T Model](w http.ResponseWriter, r *http.Request, db *gorm.DB, item *T) error {
	err := gorm.G[T](db).Create(context.Background(), item)
	if err != nil {
		return err
	}
	return nil
}

func GetByID[T Model](w http.ResponseWriter, r *http.Request, db *gorm.DB, id string) {
	result, err := gorm.G[models.Files](db).Where("id = ?", id).First(context.Background())
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	utils.ResponseJSON(w, http.StatusOK, result)
}

func Delete[T Model](w http.ResponseWriter, r *http.Request, db *gorm.DB, id string) {
	_, err := gorm.G[T](db).Where("id = ?", id).Delete(context.Background())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "Error deleting data.",
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Item deleted successfully.",
	})
}

func Init(cfg *config.Config) *gorm.DB {

	db, err := gorm.Open(postgres.Open(cfg.DBURL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&models.Files{}, &models.Task{})

	return db
}
