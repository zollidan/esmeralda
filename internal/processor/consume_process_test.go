package processor

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/zollidan/esmeralda/internal/queue"
)

func TestProcessParseTask_TaskUnmarshal(t *testing.T) {
	task := queue.ParseTask{
		ID:        "abc-123",
		Date:      "2025-01-15",
		CreatedAt: time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	decoded, err := queue.Unmarshal[queue.ParseTask](data)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.ID != task.ID {
		t.Errorf("ID = %q, want %q", decoded.ID, task.ID)
	}
	if decoded.Date != task.Date {
		t.Errorf("Date = %q, want %q", decoded.Date, task.Date)
	}
}

func TestProcessParseTask_InvalidPayload(t *testing.T) {
	_, err := queue.Unmarshal[queue.ParseTask]([]byte("invalid"))
	if err == nil {
		t.Fatal("expected error for invalid payload")
	}
}

func TestProcessParseTask_DateParsing(t *testing.T) {
	tests := []struct {
		name    string
		date    string
		wantErr bool
	}{
		{"valid", "2025-01-15", false},
		{"invalid format", "15-01-2025", true},
		{"invalid chars", "not-a-date", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := time.Parse("2006-01-02", tt.date)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse(%q) error = %v, wantErr %v", tt.date, err, tt.wantErr)
			}
		})
	}
}

func TestTaskResult_ErrorStatus(t *testing.T) {
	result := queue.TaskResult{
		TaskID: "task-1",
		Status: queue.StatusError,
		Error:  "failed to fetch matches",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	decoded, err := queue.Unmarshal[queue.TaskResult](data)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.Status != queue.StatusError {
		t.Errorf("Status = %q, want %q", decoded.Status, queue.StatusError)
	}
	if decoded.Error != "failed to fetch matches" {
		t.Errorf("Error = %q, want %q", decoded.Error, "failed to fetch matches")
	}
}

func TestTaskResult_DoneStatus(t *testing.T) {
	result := queue.TaskResult{
		TaskID: "task-1",
		Status: queue.StatusDone,
	}

	data, _ := json.Marshal(result)
	decoded, _ := queue.Unmarshal[queue.TaskResult](data)

	if decoded.Status != queue.StatusDone {
		t.Errorf("Status = %q, want %q", decoded.Status, queue.StatusDone)
	}
}
