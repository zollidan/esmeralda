package schemas

import "github.com/google/uuid"

// CreateFileRequest represents the JSON schema for creating a file
type CreateFileRequest struct {
	Filename string `json:"filename" validate:"required"`
	FileURL  string `json:"file_url" validate:"required"`
}

// CreateTaskRequest represents the JSON schema for creating a task
type CreateTaskRequest struct {
	TaskID      uuid.UUID `json:"task_id"`
	TaskName    string    `json:"task_name" validate:"required"`
	TaskComment string    `json:"task_comment"`
}
