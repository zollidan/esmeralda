package repository

import (
	"context"
	"fmt"

	"github.com/zollidan/esmeralda/internal/models"
	"gorm.io/gorm"
)

type TaskRepository struct {
	*Repository[models.Task]
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{Repository: New[models.Task](db)}
}

func (r *TaskRepository) GetAllOrdered(ctx context.Context) ([]models.Task, error) {
	var tasks []models.Task
	if err := r.DB().WithContext(ctx).Order("created_at desc").Find(&tasks).Error; err != nil {
		return nil, fmt.Errorf("get tasks ordered: %w", err)
	}
	return tasks, nil
}

func (r *TaskRepository) UpdateStatus(ctx context.Context, taskID, status string) error {
	if err := r.DB().WithContext(ctx).
		Model(&models.Task{}).
		Where("id = ?", taskID).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("update task status: %w", err)
	}
	return nil
}
