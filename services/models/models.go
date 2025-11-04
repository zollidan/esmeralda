package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

type Files struct {
	gorm.Model
	Filename string `json:"filename"`
	FileURL  string `json:"file_url"`
	S3Key    string `json:"s3_key"`    // S3 object key for easier deletion
	FileSize int64  `json:"file_size"` // File size in bytes
}

// Parser представляет тип парсера
type Parser struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
	IsActive    bool `gorm:"default:true"`
}

// Task представляет задачу для выполнения парсера
// TODO: add task status tracking
type Task struct {
	gorm.Model
	Name        string `gorm:"not null"`
	Description string

	// Связь с парсером
	ParserID uint   `gorm:"not null;index"`
	Parser   Parser `gorm:"foreignKey:ParserID"`

	// Конфигурация запуска
	CronExpression string `gorm:"size:100"` // "*/5 * * * *" или пусто для разовой задачи
	IsRecurring    bool   `gorm:"default:false"`
	IsActive       bool   `gorm:"default:true;index"`

	// Метаданные выполнения
	JobID      string `gorm:"size:100;index"`            // ID job'а в scheduler
	TaskStatus string `gorm:"size:50;default:'pending'"` // pending, running, completed, failed
}
