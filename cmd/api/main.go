package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zollidan/esmeralda/internal/config"
	"github.com/zollidan/esmeralda/internal/db"
	"github.com/zollidan/esmeralda/internal/queue"
	"github.com/zollidan/esmeralda/internal/server"

	"context"
	"os"
	"os/signal"
	"syscall"
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

	parseProducer := queue.NewProducer(rdb, queue.StreamParse)

	handler := server.NewHandler(parseProducer, rdb, database)
	handler.StartConsumers(ctx)

	r := gin.Default()
	server.SetupRoutes(r, handler)

	if err := r.Run(cfg.ServerPort); err != nil {
		log.Fatal(err)
	}
}
