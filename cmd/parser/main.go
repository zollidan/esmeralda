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
	"github.com/zollidan/esmeralda/internal/processor"
	"github.com/zollidan/esmeralda/internal/queue"
)

func main() {
	cfg := config.InitConfig()

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})

	client := api.NewClient(cfg.SportAPIRU.BaseURL, cfg.SportAPIRU.Token)
	database, err := db.New(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}

	parseConsumer := queue.NewConsumer(rdb, queue.StreamParse, "parse_group", "parse_consumer")
	resultsProducer := queue.NewProducer(rdb, queue.StreamResults)

	processor := processor.Init(client, database, resultsProducer)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	err = parseConsumer.Consume(ctx, processor.ProcessParseTask)

	if err != nil && err != context.Canceled {
		log.Fatal(err)
	}
}
