package models

import "time"

type Task struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	Date      string    `gorm:"not null" json:"date"`
	Status    string    `gorm:"not null;default:pending" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
