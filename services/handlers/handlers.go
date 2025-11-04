package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/zollidan/esmeralda/config"
	"github.com/zollidan/esmeralda/s3storage"
	"github.com/zollidan/esmeralda/tasks"
	"gorm.io/gorm"
)

type Handlers struct {
	DB          *gorm.DB
	Cfg         *config.Config
	TaskManager *tasks.Manager
	S3Client    *s3storage.S3Storage
}

// Constructor
func New(db *gorm.DB, cfg *config.Config, tm *tasks.Manager, s3Client *s3storage.S3Storage) *Handlers {
	return &Handlers{
		DB:          db,
		Cfg:         cfg,
		TaskManager: tm,
		S3Client:    s3Client,
	}
}

// Метод для регистрации всех роутов
func (h *Handlers) RegisterRoutes(r chi.Router) {
	r.Route("/auth", h.AuthRoutes)
	r.Route("/files", h.FilesRoutes)
	r.Route("/tasks", h.TasksRoutes)
	r.Route("/parsers", h.ParsersRoutes)
}
