package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Files struct {
	gorm.Model
	Filename string `json:"filename"`
	FileURL  string `json:"file_url"`
}

type Tasks struct {
	gorm.Model
	TaskID      uuid.UUID
	TaskName    string
	TaskComment string
}
