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
	"github.com/zollidan/esmeralda/mq"
	"github.com/zollidan/esmeralda/scheduler"
	"gorm.io/gorm"
)

func main() {

	// configuration setup
	cfg := config.New()

	// database setup
	db := database.InitDatabase(cfg)

	// cron/scheduler setup
	sched := scheduler.InitScheduler()
	defer sched.Shutdown()
	sched.Start()

	mq := mq.NewMQ()
	defer mq.CloseMQ()

	mqCh := mq.GetChannel()
	defer mqCh.Close()

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

				var file models.Files

				defer r.Body.Close()

				if err := json.NewDecoder(r.Body).Decode(&file); err != nil {
					http.Error(w, "invalid JSON", http.StatusBadRequest)
					return
				}

				err := gorm.G[models.Files](db).Create(context.Background(), &models.Files{Filename: file.Filename, FileURL: file.FileURL})
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to create file: %v", err), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				json.NewEncoder(w).Encode(file)
			})
		})
		r.Route("/tasks", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {

				tasks, err := gorm.G[models.Tasks](db).Find(context.Background())
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
				var task models.Tasks
				defer r.Body.Close()

				if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
					http.Error(w, "invalid JSON", http.StatusBadRequest)
					return
				}

				err := mq.SendMessage(mqCh, "test hello rabbitMQ", "Hello World RabbitMQ!!!")
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to send message: %v", err), http.StatusInternalServerError)
					return
				}

				jobID, err := sched.NewJob()
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to create job: %v", err), http.StatusInternalServerError)
					return
				}

				err = gorm.G[models.Tasks](db).Create(context.Background(), &models.Tasks{TaskID: task.TaskID, TaskName: task.TaskName, TaskComment: task.TaskComment})
				if err != nil {
					http.Error(w, fmt.Sprintf("failed to create task: %v", err), http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(map[string]string{
					"job_id": jobID.String(),
				})
			})
		})
	})

	fmt.Printf("Server is running on %s", cfg.ServerAddress)
	http.ListenAndServe(cfg.ServerAddress, r)
}
