package models

import (
	"time"

	"gorm.io/gorm"
)

type Files struct {
	gorm.Model
	Filename string `json:"filename"`
	FileURL  string `json:"file_url"`
}

// Parser представляет тип парсера
type Parser struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex;not null"`
	Description string
	IsActive    bool `gorm:"default:true"`
}

// Task представляет задачу для выполнения парсера
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
	JobID      string     `gorm:"size:100;index"` // ID job'а в scheduler
	LastRun    *time.Time
	NextRun    *time.Time
	RunCount   int    `gorm:"default:0"`
	FailCount  int    `gorm:"default:0"`
	LastStatus string `gorm:"size:50"` // success, failed, running, pending
	LastError  string `gorm:"type:text"`
	
	// История запусков
	Runs []TaskRun `gorm:"foreignKey:TaskID"`
}

// TaskRun история выполнения задач
type TaskRun struct {
	gorm.Model
	TaskID uint `gorm:"not null;index:idx_task_runs"`
	Task   Task `gorm:"foreignKey:TaskID"`
	
	Status     string     `gorm:"size:50;not null;index"` // running, success, failed, cancelled
	StartedAt  time.Time  `gorm:"not null;index:idx_task_runs"`
	FinishedAt *time.Time
	Duration   int64 `gorm:"default:0"` // миллисекунды
	
	ErrorMessage string `gorm:"type:text"`
	
	// Метрики
	ItemsParsed  int `gorm:"default:0"`
	ItemsSkipped int `gorm:"default:0"`
	ItemsFailed  int `gorm:"default:0"`
}