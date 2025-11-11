package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)


func (h *Handlers) ParsersRoutes(r chi.Router) {
	r.Get("/", h.GetParsers)
	r.Post("/refresh", h.RefreshParsers)
}

// GetAll - получение всех парсеров
func (pc *ParserCache) GetAll() []Parser {
    pc.mu.RLock()
    defer pc.mu.RUnlock()
    // возвращаем копию, чтобы избежать race conditions
    result := make([]Parser, len(pc.parsers))
    copy(result, pc.parsers)
    return result
}

func (h *Handlers) RefreshParsers(w http.ResponseWriter, r *http.Request) {
    // TODO: реализация refresh
    w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handlers) GetParsers(w http.ResponseWriter, r *http.Request) {
    parsers := h.parserCache.GetAll()
    
    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(parsers); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}
