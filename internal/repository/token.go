package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zollidan/esmeralda/internal/models"
	"gorm.io/gorm"
)

type RefreshTokenRepository struct {
	*Repository[models.RefreshToken]
}

func NewRefreshTokenRepository(db *gorm.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{Repository: New[models.RefreshToken](db)}
}

func (r *RefreshTokenRepository) FindByToken(ctx context.Context, token string) (*models.RefreshToken, error) {
	var t models.RefreshToken
	err := r.DB().
		WithContext(ctx).
		Where("token = ?", token).
		First(&t).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("find refresh token: %w", err)
	}
	return &t, nil
}

func (r *RefreshTokenRepository) DeleteByToken(ctx context.Context, token string) error {
	if err := r.DB().
		WithContext(ctx).
		Where("token = ?", token).
		Delete(&models.RefreshToken{}).Error; err != nil {
		return fmt.Errorf("delete refresh token: %w", err)
	}
	return nil
}

func (r *RefreshTokenRepository) DeleteExpired(ctx context.Context) error {
	if err := r.DB().
		WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&models.RefreshToken{}).Error; err != nil {
		return fmt.Errorf("delete expired tokens: %w", err)
	}
	return nil
}
