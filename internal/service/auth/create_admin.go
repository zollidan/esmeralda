package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zollidan/esmeralda/internal/models"
	"github.com/zollidan/esmeralda/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

func CreateAdmin(username, password string, repo *repository.UserRepository) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exists, err := repo.AnyExists(ctx)
	if err != nil {
		return "", err
	}

	if exists {
		return "admin already exists, skip", nil
	}

	if username == "" || password == "" {
		return "", errors.New("username and password required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password: %w", err)
	}

	user := &models.User{
		Username:     username,
		PasswordHash: string(hash),
	}

	if err := repo.Create(ctx, user); err != nil {
		return "", fmt.Errorf("create admin: %w", err)
	}

	return "admin created", nil
}
