package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zollidan/esmeralda/internal/config"
	"github.com/zollidan/esmeralda/internal/queue"
	"github.com/zollidan/esmeralda/internal/server"
)

func main() {
	cfg := config.Load()

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	producer := queue.NewProducer(rdb)
	handler := server.NewHandler(producer)

	r := gin.Default()
	server.SetupRoutes(r, handler)

	if err := r.Run(cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}