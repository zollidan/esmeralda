package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	consumeProcess := processor.Init(client, gameRepo, resultsProducer, progressProducer, cfg.Tech.Workers)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if _, err := w.Write([]byte(`{"status":"ok"}`)); err != nil {
				log.Printf("failed to write health response: %v", err)
			}
		})

		srv := &http.Server{
			Addr:              cfg.Tech.ParserPort,
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		}

		if err := srv.ListenAndServe(); err != nil {
			log.Printf("health server: %v", err)
		}
	}()

	log.Println("Processor started")
	err = parseConsumer.Consume(ctx, consumeProcess.ProcessParseTask)

	if err != nil && !errors.Is(err, context.Canceled) {
		log.Print(err)
	}
}
