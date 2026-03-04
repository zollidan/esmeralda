package queue

import "time"

const (
	StreamParse   = "tasks:parse"
	StreamResults = "tasks:results"
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
