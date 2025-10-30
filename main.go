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
	"github.com/zollidan/esmeralda/s3storage"
	"github.com/zollidan/esmeralda/tasks"
)

func main() {

	// Configuration setup
	cfg := config.New()

	// Database setup
	db := database.InitDatabase(cfg)

	// S3Storage setup
	s3Client := s3storage.New(&cfg)

	// Task manager setup
	taskManager := tasks.NewManager()
	defer taskManager.Shutdown()

	// HTTP router setup
	r := chi.NewRouter()

	// Register handlers
	h := handlers.New(db, &cfg, taskManager, s3Client)

	// Middleware
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Logger)

	r.Route("/api/v1", func(r chi.Router) {
		h.RegisterRoutes(r)
	})

	fmt.Printf("Server is running on http://%s\n", cfg.ServerAddress)
	log.Printf("Server started on http://%s\nMode: %s\n", cfg.ServerAddress, cfg.AppMode)
	if err := http.ListenAndServe(cfg.ServerAddress, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
