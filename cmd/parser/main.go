package main

import (
	"context"
	"errors"
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
	"github.com/zollidan/esmeralda/internal/repository"
)

func main() {
	cfg := config.Load()

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})

	client := api.NewClient(cfg.SportAPIRU.BaseURL, cfg.SportAPIRU.Token)
	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}

	gameRepo := repository.NewGameRepository(database)

	parseConsumer := queue.NewConsumer(rdb, queue.StreamParse, "parse_group", "parse_consumer")
	resultsProducer := queue.NewProducer(rdb, queue.StreamResults)
	progressProducer := queue.NewProducer(rdb, queue.StreamProgress)

	consumeProcess := processor.Init(client, gameRepo, resultsProducer, progressProducer)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// go func() {
	// 	mux := http.NewServeMux()
	// 	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusOK)
	// 		w.Write([]byte(`{"status":"ok"}`))
	// 	})
	// 	if err := http.ListenAndServe(":8081", mux); err != nil {
	// 		log.Printf("health server: %v", err)
	// 	}
	// }()

	err = parseConsumer.Consume(ctx, consumeProcess.ProcessParseTask)

	if err != nil && !errors.Is(err, context.Canceled) {
		log.Fatal(err)
	}
}
