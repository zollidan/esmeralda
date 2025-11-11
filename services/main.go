package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zollidan/esmeralda/config"
	"github.com/zollidan/esmeralda/database"
	handlers "github.com/zollidan/esmeralda/handlers"
	"github.com/zollidan/esmeralda/mq"
	"github.com/zollidan/esmeralda/rediska"
	"github.com/zollidan/esmeralda/s3storage"
)

func main() {
	// Configuration setup
	cfg := config.New()

	// Database setup
	db := database.InitDatabase(cfg)

	// S3Storage setup
	s3Client := s3storage.New(cfg)

	// RabbitMQ setup
	mqClient := mq.NewMQ()
	defer mqClient.CloseMQ()

	// Redis setup
	redisClient := rediska.Init(cfg.RedisURL)

	// HTTP router setup
	r := chi.NewRouter()

	// Register handlers
	h := handlers.New(db, cfg, mqClient, s3Client, redisClient.Client)

	// Middleware
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Logger)

	r.Route("/api/v1", func(r chi.Router) {
		h.RegisterRoutes(r)
	})

	fmt.Println()
	fmt.Println("Esemeralda - aaf-bet.ru big data and parsing service")
	fmt.Println("==============================================")
	fmt.Printf("  Server is running on: \033[1;36mhttp://%s\033[0m\n", cfg.ServerAddress)
	fmt.Printf("  Mode: \033[1;33m%s\033[0m\n", cfg.AppMode)
	fmt.Println("  Press CTRL+C to stop the server")
	fmt.Println("==============================================")
	fmt.Println()

	log.Printf("Server started on http://%s | Mode: %s\n", cfg.ServerAddress, cfg.AppMode)

	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}