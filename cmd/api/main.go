package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zollidan/esmeralda-ru-api-fetcher/internal/config"
	"github.com/zollidan/esmeralda-ru-api-fetcher/internal/db"
	"github.com/zollidan/esmeralda-ru-api-fetcher/internal/queue"
	"github.com/zollidan/esmeralda-ru-api-fetcher/internal/server"
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

	producer := queue.NewProducer(rdb)
	handler := server.NewHandler(producer, database)

	r := gin.Default()
	server.SetupRoutes(r, handler)

	if err := r.Run(cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}