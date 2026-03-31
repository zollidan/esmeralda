package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zollidan/esmeralda/internal/models"
	"github.com/zollidan/esmeralda/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

type AdminCredentials struct {
	Username string
	Password string
}

type UserRepository interface {
	AnyExists(ctx context.Context) (bool, error)
	Create(ctx context.Context, user *models.User) error
}

func CreateAdmin(username string, repo UserRepository) (*AdminCredentials, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := repo.AnyExists(ctx)
	if err != nil {
		return nil, "", err
	}

	if exists {
		return nil, "admin already exists, skip", nil
	}

	if username == "" {
		return nil, "", errors.New("username required")
	}

	password, err := utils.GeneratePassword(10)
	if err != nil {
		return nil, "", fmt.Errorf("error generating password: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("error hashing password: %w", err)
	}

	user := &models.User{
		Username:     username,
		PasswordHash: string(hash),
	}

	if err := repo.Create(ctx, user); err != nil {
		return nil, "", fmt.Errorf("error creating admin: %w", err)
	}

	return &AdminCredentials{
		Username: username,
		Password: password,
	}, "admin created", nil
}
