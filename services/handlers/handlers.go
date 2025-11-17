package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/zollidan/esmeralda/config"
	"github.com/zollidan/esmeralda/rediska"
	"github.com/zollidan/esmeralda/s3storage"
	"gorm.io/gorm"
)

type Handlers struct {
	DB          *gorm.DB
	Cfg         *config.Config
	S3Client    *s3storage.S3Storage
	RedisClient *rediska.RedisClient
}

// Constructor
func New(db *gorm.DB, cfg *config.Config, s3Client *s3storage.S3Storage, redisClient *rediska.RedisClient) *Handlers {
	return &Handlers{
		DB:          db,
		Cfg:         cfg,
		S3Client:    s3Client,
		RedisClient: redisClient,
	}
}

// Метод для регистрации всех роутов
func (h *Handlers) RegisterRoutes(r chi.Router) {
	r.Route("/files", h.FilesRoutes)
	r.Route("/tasks", h.TasksRoutes)
	r.Route("/parsers", h.ParsersRoutes)
}
