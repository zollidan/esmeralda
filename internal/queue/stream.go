package queue

import (
	"time"

	"github.com/zollidan/esmeralda/internal/models"
)

const (
	StreamParse         = "tasks:parse"
	StreamResults       = "tasks:results"
	StreamEnrich        = "tasks:enrich"
	StreamEnrichResults = "tasks:enrich_results"
)

type Status string

const (
	StatusPending    Status = "pending"
	StatusProcessing Status = "processing"
	StatusDone       Status = "done"
	StatusError      Status = "error"
	StatusFailed     Status = "failed"
	StatusCancelled  Status = "cancelled"
)

type MatchDataTask struct {
	ID        string    `json:"id"`
	DateStart string    `json:"date_start"` // "2006-01-02"
	DateEnd   string    `json:"date_end"`   // "2006-01-02"
	CreatedAt time.Time `json:"created_at"`
}

type MatchDataResult struct {
	TaskID    string        `json:"task_id"`
	DateStart string        `json:"date_start"` // "2006-01-02"
	DateEnd   string        `json:"date_end"`   // "2006-01-02"
	Matches   []models.Game `json:"matches"`
}

type ParseTask struct {
	ID        string    `json:"id"`
	Date      string    `json:"date"` // "2006-01-02"
	CreatedAt time.Time `json:"created_at"`
}

type TaskResult struct {
	TaskID string `json:"task_id"`
	Status Status `json:"status"`
	Error  string `json:"error,omitempty"`
}
