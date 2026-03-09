package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Repository is a generic base repository providing common CRUD operations.
type Repository[T any] struct {
	db *gorm.DB
}

func New[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

func (r *Repository[T]) DB() *gorm.DB {
	return r.db
}

func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return fmt.Errorf("repository create: %w", err)
	}
	return nil
}

func (r *Repository[T]) CreateInBatches(ctx context.Context, entities []T, batchSize int) error {
	if err := r.db.WithContext(ctx).CreateInBatches(entities, batchSize).Error; err != nil {
		return fmt.Errorf("repository create in batches: %w", err)
	}
	return nil
}

func (r *Repository[T]) FindAll(ctx context.Context) ([]T, error) {
	var entities []T
	if err := r.db.WithContext(ctx).Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("repository find all: %w", err)
	}
	return entities, nil
}

func (r *Repository[T]) FindByID(ctx context.Context, id any) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("repository find by id: %w", err)
	}
	return &entity, nil
}

func (r *Repository[T]) Update(ctx context.Context, entity *T) error {
	if err := r.db.WithContext(ctx).Save(entity).Error; err != nil {
		return fmt.Errorf("repository update: %w", err)
	}
	return nil
}

func (r *Repository[T]) Delete(ctx context.Context, id any) error {
	var entity T
	if err := r.db.WithContext(ctx).Delete(&entity, "id = ?", id).Error; err != nil {
		return fmt.Errorf("repository delete: %w", err)
	}
	return nil
}
