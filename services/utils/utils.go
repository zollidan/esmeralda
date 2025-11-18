package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/zollidan/esmeralda/database"
	"github.com/zollidan/esmeralda/models"
	"github.com/zollidan/esmeralda/rediska"
	"gorm.io/gorm"
)

type Task struct {
	TaskID string `json:"task_id"`
	Status string `json:"status"`
	Result string `json:"result"`
}

func ResponseJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

// failOnError logs the provided message and panics if err is non-nil.
// This is a helper used in this package for simplicity.
func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Under construction
func StartFetchingResults(redisClient *rediska.RedisClient, db *gorm.DB, taskID uint) error {
	data, err := redisClient.FindResults()
	if err == redis.Nil {
		return err
	} else if err != nil {
		return err
	}

	for _, item := range data {
		var task Task
		err := json.Unmarshal([]byte(item), &task)
		if err != nil {
			continue
		}

		result, err := database.GetBy[models.Task](db, "id", task.TaskID)
		if err != nil {
			return err
		}

		if result.TaskStatus != task.Status && task.Status == "success" {
			continue
		}

		redisClient.Client.Del("result:" + task.TaskID)
	}

	return nil
}
