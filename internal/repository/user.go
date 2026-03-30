package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/zollidan/esmeralda/internal/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	*Repository[models.User]
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{Repository: New[models.User](db)}
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User

	err := r.DB().
		WithContext(ctx).
		Where("username = ?", username).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("find user by username: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) AnyExists(ctx context.Context) (bool, error) {
	var count int64

	if err := r.DB().
		WithContext(ctx).
		Model(&models.User{}).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("check users exist: %w", err)
	}

	return count > 0, nil
}