package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/zollidan/esmeralda/config"
	"github.com/zollidan/esmeralda/database"
	handlers "github.com/zollidan/esmeralda/hadlers"
	"github.com/zollidan/esmeralda/models"
	"github.com/zollidan/esmeralda/schemas"
	"github.com/zollidan/esmeralda/tasks"
	"gorm.io/gorm"
)

func main() {

	// Configuration setup
	cfg := config.New()

	// Database setup
	db := database.InitDatabase(cfg)

	// Task manager setup
	taskManager := tasks.NewManager()
	defer taskManager.Shutdown()

	// HTTP router setup
	r := chi.NewRouter()

	// Register handlers
	h := handlers.New(db, &cfg, taskManager)

	// Middleware
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Logger)

	r.Route("/api/v1", func(r chi.Router) {

		h.RegisterRoutes(r)

		r.Route("/tasks", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				tasks, err := gorm.G[models.Task](db).Preload("Parser", nil).Find(context.Background())
				if err != nil {
					json.NewEncoder(w).Encode(map[string]string{
						"error": "Error reading tasks.",
					})
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(tasks)
			})
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				var req schemas.CreateTaskRequest
				defer r.Body.Close()

				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, "invalid JSON", http.StatusBadRequest)
					return
				}

				// Find parser to ensure it exists
				_, err := gorm.G[models.Parser](db).Where("id = ?", req.ParserID).First(context.Background())
				if err != nil {
					if err == gorm.ErrRecordNotFound {
						http.Error(w, "parser not found", http.StatusNotFound)
					} else {
						http.Error(w, fmt.Sprintf("failed to find parser: %v", err), http.StatusInternalServerError)
					}
					return
				}

				task := &models.Task{
					Name:        req.Name,
					Description: req.Description,
					ParserID:    req.ParserID,
				}

				err = taskManager.StartTask(task.ID)
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to start task: %v", err), http.StatusInternalServerError)
					return
				}

				err = gorm.G[models.Task](db).Create(context.Background(), task)
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to create task in database: %v", err), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"task_id": task.ID,
				})
			})
		})
	})

	fmt.Printf("Server is running on http://%s\n", cfg.ServerAddress)
	http.ListenAndServe(cfg.ServerAddress, r)
}
