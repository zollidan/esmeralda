package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis"
	"github.com/zollidan/esmeralda/config"
	"github.com/zollidan/esmeralda/mq"
	"github.com/zollidan/esmeralda/s3storage"
	"gorm.io/gorm"
)
type Handlers struct {
	DB          *gorm.DB
	Cfg         *config.Config
	MQ          *mq.MQ
	S3Client    *s3storage.S3Storage
	RedisClient *redis.Client
}

// Constructor
func New(db *gorm.DB, cfg *config.Config, mqClient *mq.MQ, s3Client *s3storage.S3Storage, redisClient *redis.Client) *Handlers {

	// cache := NewParserCache()

	// mqClient.RefreshParsersRPC("parsers", cache)
	return &Handlers{
		DB:          db,
		Cfg:         cfg,
		MQ:          mqClient,
		S3Client:    s3Client,
		RedisClient: redisClient,
		// parserCache: cache,
	}
}

// Метод для регистрации всех роутов
func (h *Handlers) RegisterRoutes(r chi.Router) {
	r.Route("/auth", h.AuthRoutes)
	r.Route("/files", h.FilesRoutes)
	r.Route("/tasks", h.TasksRoutes)
	r.Route("/parsers", h.ParsersRoutes)
}
