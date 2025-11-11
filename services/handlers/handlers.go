package handlers

import (
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/zollidan/esmeralda/config"
	"github.com/zollidan/esmeralda/mq"
	"github.com/zollidan/esmeralda/s3storage"
	"gorm.io/gorm"
)

// Parser - структура парсера (адаптируйте под ваши нужды)
type Parser struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    // добавьте другие поля
}

// ParserCache - in-memory кэш для парсеров
type ParserCache struct {
    mu      sync.RWMutex
    parsers []Parser
}

// NewParserCache - создание нового кэша
func NewParserCache() *ParserCache {
    return &ParserCache{
        parsers: make([]Parser, 0),
    }
}

// Add - добавление парсера в кэш
func (pc *ParserCache) Add(parser Parser) {
    pc.mu.Lock()
    defer pc.mu.Unlock()
    pc.parsers = append(pc.parsers, parser)
}

type Handlers struct {
	DB          *gorm.DB
	Cfg         *config.Config
	MQ          *mq.MQ
	S3Client    *s3storage.S3Storage
	parserCache *ParserCache
}

// Constructor
func New(db *gorm.DB, cfg *config.Config, mqClient *mq.MQ, s3Client *s3storage.S3Storage) *Handlers {

	// cache := NewParserCache()

	// mqClient.RefreshParsersRPC("parsers", cache)
	return &Handlers{
		DB:          db,
		Cfg:         cfg,
		MQ:          mqClient,
		S3Client:    s3Client,
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
