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

	"context"
	"os"
	"os/signal"
	"syscall"
)

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

	handler := server.NewHandler(parseProducer, rdb, taskRepo, gameRepo)
	handler.StartConsumers(ctx)

	r := gin.Default()
	server.SetupRoutes(r, handler)

	srv := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: r,
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
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server stopped")
}
