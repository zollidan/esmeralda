package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/zollidan/esmeralda/config"
	"github.com/zollidan/esmeralda/tasks"
	"gorm.io/gorm"
)

type Handlers struct {
	DB          *gorm.DB
	Cfg         *config.Config
	TaskManager *tasks.Manager
}

// Constructor
func New(db *gorm.DB, cfg *config.Config, tm *tasks.Manager) *Handlers {
	return &Handlers{
		DB:          db,
		Cfg:         cfg,
		TaskManager: tm,
	}
}

// Метод для регистрации всех роутов
func (h *Handlers) RegisterRoutes(r chi.Router) {
	r.Route("/files", h.FilesRoutes)
	// r.Route("/tasks", h.TasksRoutes)
	r.Route("/parsers", h.ParsersRoutes)
	r.Route("/auth", h.AuthRoutes)
}
