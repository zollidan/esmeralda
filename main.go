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
	"github.com/zollidan/esmeralda/models"
	"github.com/zollidan/esmeralda/schemas"
	"github.com/zollidan/esmeralda/tasks"
	"gorm.io/gorm"
)

func main() {

	// configuration setup
	cfg := config.New()

	// database setup
	db := database.InitDatabase(cfg)

	// task manager setup
	taskManager := tasks.NewManager()
	defer taskManager.Shutdown()

	// HTTP router setup
	r := chi.NewRouter()

	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Logger)

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/files", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {

				files, err := gorm.G[models.Files](db).Find(context.Background())
				if err != nil {
					json.NewEncoder(w).Encode(map[string]string{
						"error": "Error reading files.",
					})
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(files)
			})
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				var req schemas.CreateFileRequest

				defer r.Body.Close()

				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, "invalid JSON", http.StatusBadRequest)
					return
				}

				err := gorm.G[models.Files](db).Create(context.Background(), &models.Files{Filename: req.Filename, FileURL: req.FileURL})
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to create file: %v", err), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(req)
			})
		})
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
		r.Route("/parsers", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				parsers, err := gorm.G[models.Parser](db).Find(context.Background())
				if err != nil {
					json.NewEncoder(w).Encode(map[string]string{
						"error": "Error reading parsers.",
					})
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(parsers)
			})
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				var req schemas.CreateParserRequest

				defer r.Body.Close()

				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					http.Error(w, "invalid JSON", http.StatusBadRequest)
					return
				}

				err := gorm.G[models.Parser](db).Create(context.Background(), &models.Parser{Name: req.ParserName, Description: req.ParserDescription})
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to create parser: %v", err), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(req)
			})
		})
	})

	fmt.Printf("Server is running on http://%s\n", cfg.ServerAddress)
	http.ListenAndServe(cfg.ServerAddress, r)
}
