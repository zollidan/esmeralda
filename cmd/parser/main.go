package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
	"github.com/zollidan/esmeralda/internal/api"
	"github.com/zollidan/esmeralda/internal/config"
	"github.com/zollidan/esmeralda/internal/db"
	"github.com/zollidan/esmeralda/internal/queue"
)

// в бд сохраеяются значения на все матчи нулями, сделать фикс

func main() {

	cfg := config.InitConfig()

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	client := api.NewClient(cfg.SportAPIRU.BaseURL, cfg.SportAPIRU.Token)
	database, err := db.New(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}

	consumer := queue.NewConsumer(rdb, client, database)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := consumer.Consume(ctx); err != nil && err != context.Canceled {
		log.Fatal(err)
	}
}

