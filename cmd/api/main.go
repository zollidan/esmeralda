package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zollidan/esmeralda/internal/config"
	"github.com/zollidan/esmeralda/internal/db"
	"github.com/zollidan/esmeralda/internal/queue"
	"github.com/zollidan/esmeralda/internal/repository"
	"github.com/zollidan/esmeralda/internal/server"
	"github.com/zollidan/esmeralda/internal/sse"

	"context"
	"os"
	"os/signal"
	"syscall"
)

// @title           Esmeralda API
// @version         1.0
// @description     API для управления задачами парсинга
// @host            localhost:8080
// @BasePath        /api
func main() {
	cfg := config.Load()

	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("connect db: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	parseProducer := queue.NewProducer(rdb, queue.StreamParse)
	taskRepo := repository.NewTaskRepository(database)
	gameRepo := repository.NewGameRepository(database)

	hub := sse.NewHub()

	handler := server.NewHandler(parseProducer, rdb, taskRepo, gameRepo, hub)
	handler.StartConsumers(ctx)

	r := gin.Default()
	server.SetupRoutes(r, handler)

	srv := &http.Server{
		Addr:              cfg.ServerPort,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
	}

	log.Println("server stopped")
}
