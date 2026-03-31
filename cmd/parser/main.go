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
	log.Printf("config loaded: redis=%s workers=%d games_limit=%d parser_port=%s",
		cfg.RedisAddr, cfg.Tech.Workers, cfg.Tech.GamesLimit, cfg.Tech.ParserPort)

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("redis ping failed: %v", err)
	}
	log.Printf("redis connected: %s", cfg.RedisAddr)

	client := api.NewClient(cfg.SportAPIRU.BaseURL, cfg.SportAPIRU.Token)
	log.Printf("sport api client created: base_url=%s", cfg.SportAPIRU.BaseURL)

	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	log.Println("database connected")

	gameRepo := repository.NewGameRepository(database)

	parseConsumer := queue.NewConsumer(rdb, queue.StreamParse, "parse_group", "parse_consumer")
	resultsProducer := queue.NewProducer(rdb, queue.StreamResults)
	progressProducer := queue.NewProducer(rdb, queue.StreamProgress)
	log.Printf("queues initialized: parse=%s results=%s progress=%s",
		queue.StreamParse, queue.StreamResults, queue.StreamProgress)

	consumeProcess := processor.Init(client, gameRepo, resultsProducer, progressProducer, cfg.Tech.Workers, cfg.Tech.GamesLimit)

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

		log.Printf("health server listening on %s", cfg.Tech.ParserPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("health server error: %v", err)
		}
	}()

	log.Println("processor started, waiting for tasks...")
	err = parseConsumer.Consume(ctx, consumeProcess.ProcessParseTask)

	if err != nil && !errors.Is(err, context.Canceled) {
		log.Printf("consumer stopped with error: %v", err)
	} else {
		log.Println("processor stopped gracefully")
	}
}
