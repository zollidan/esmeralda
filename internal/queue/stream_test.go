package queue

import (
	"encoding/json"
	"testing"
	"time"
)

func TestUnmarshal_ParseTask(t *testing.T) {
	task := ParseTask{
		ID:        "abc-123",
		Date:      "2025-01-15",
		CreatedAt: time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got, err := Unmarshal[ParseTask](data)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.ID != task.ID {
		t.Errorf("ID = %q, want %q", got.ID, task.ID)
	}
	if got.Date != task.Date {
		t.Errorf("Date = %q, want %q", got.Date, task.Date)
	}
}

func TestUnmarshal_TaskResult(t *testing.T) {
	result := TaskResult{
		TaskID: "task-1",
		Status: StatusDone,
		Error:  "",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got, err := Unmarshal[TaskResult](data)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.TaskID != result.TaskID {
		t.Errorf("TaskID = %q, want %q", got.TaskID, result.TaskID)
	}
	if got.Status != result.Status {
		t.Errorf("Status = %q, want %q", got.Status, result.Status)
	}
}

func TestUnmarshal_TaskProgress(t *testing.T) {
	progress := TaskProgress{
		TaskID:       "task-2",
		Status:       StatusProcessing,
		TotalMatches: 20,
		CurrentMatch: 5,
	}

	data, err := json.Marshal(progress)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got, err := Unmarshal[TaskProgress](data)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.TotalMatches != 20 {
		t.Errorf("TotalMatches = %d, want 20", got.TotalMatches)
	}
	if got.CurrentMatch != 5 {
		t.Errorf("CurrentMatch = %d, want 5", got.CurrentMatch)
	}
}

func TestUnmarshal_InvalidJSON(t *testing.T) {
	_, err := Unmarshal[ParseTask]([]byte("not json"))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestUnmarshal_TaskResultWithError(t *testing.T) {
	result := TaskResult{
		TaskID: "task-err",
		Status: StatusError,
		Error:  "something went wrong",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got, err := Unmarshal[TaskResult](data)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.Error != "something went wrong" {
		t.Errorf("Error = %q, want %q", got.Error, "something went wrong")
	}
}

func TestUnmarshal_TaskResultOmitEmptyError(t *testing.T) {
	result := TaskResult{
		TaskID: "task-ok",
		Status: StatusDone,
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	// Verify that "error" field is omitted when empty
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("unmarshal raw: %v", err)
	}
	if _, exists := raw["error"]; exists {
		t.Error("expected 'error' field to be omitted when empty")
	}
}

func TestMatchDataTask_JSON(t *testing.T) {
	task := MatchDataTask{
		ID:        "mdt-1",
		DateStart: "2025-01-01",
		DateEnd:   "2025-01-31",
		CreatedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	got, err := Unmarshal[MatchDataTask](data)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if got.DateStart != "2025-01-01" {
		t.Errorf("DateStart = %q, want %q", got.DateStart, "2025-01-01")
	}
	if got.DateEnd != "2025-01-31" {
		t.Errorf("DateEnd = %q, want %q", got.DateEnd, "2025-01-31")
	}
}

func TestStatusConstants(t *testing.T) {
	tests := []struct {
		status Status
		want   string
	}{
		{StatusPending, "pending"},
		{StatusProcessing, "processing"},
		{StatusDone, "done"},
		{StatusError, "error"},
		{StatusFailed, "failed"},
		{StatusCancelled, "cancelled"},
	}

	for _, tt := range tests {
		if string(tt.status) != tt.want {
			t.Errorf("Status = %q, want %q", tt.status, tt.want)
		}
	}
}

func TestStreamConstants(t *testing.T) {
	if StreamParse != "tasks:parse" {
		t.Errorf("StreamParse = %q, want %q", StreamParse, "tasks:parse")
	}
	if StreamResults != "tasks:results" {
		t.Errorf("StreamResults = %q, want %q", StreamResults, "tasks:results")
	}
	if StreamProgress != "tasks:progress" {
		t.Errorf("StreamProgress = %q, want %q", StreamProgress, "tasks:progress")
	}
}
