package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zollidan/esmeralda/internal/config"
	"github.com/zollidan/esmeralda/internal/db"
	"github.com/zollidan/esmeralda/internal/models"
	"github.com/zollidan/esmeralda/internal/queue"
	"github.com/zollidan/esmeralda/internal/server"
)

func main() {
	cfg := config.Load()

	database, err := db.New(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	resultsConsumer := queue.NewConsumer(rdb, queue.StreamResults, "results_group", "results_consumer")
	go func() {
		err := resultsConsumer.Consume(ctx, func(ctx context.Context, payload []byte) error {
			result, err := queue.Unmarshal[queue.TaskResult](payload)
			if err != nil {
				return err
			}
			return database.WithContext(ctx).
				Model(&models.Task{}).
				Where("id = ?", result.TaskID).
				Update("status", string(result.Status)).Error
		})
		if err != nil && err != context.Canceled {
			log.Printf("results consumer error: %v", err)
		}
	}()

	producer := queue.NewProducer(rdb, queue.StreamParse)
	handler := server.NewHandler(producer, database)

	r := gin.Default()
	server.SetupRoutes(r, handler)

	if err := r.Run(cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}
