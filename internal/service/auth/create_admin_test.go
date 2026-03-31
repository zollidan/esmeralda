package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/zollidan/esmeralda/internal/models"
)

type mockUserRepo struct {
	anyExists    bool
	anyExistsErr error
	createErr    error
}

func (m *mockUserRepo) AnyExists(_ context.Context) (bool, error) {
	return m.anyExists, m.anyExistsErr
}

func (m *mockUserRepo) Create(_ context.Context, _ *models.User) error {
	return m.createErr
}

func TestCreateAdmin(t *testing.T) {
	tests := []struct {
		name          string
		username      string
		repo          *mockUserRepo
		expectError   bool
		expectMessage string
		expectCreds   bool
	}{
		{
			name:          "admin already exists",
			username:      "admin",
			repo:          &mockUserRepo{anyExists: true},
			expectError:   false,
			expectMessage: "admin already exists, skip",
			expectCreds:   false,
		},
		{
			name:        "empty username",
			username:    "",
			repo:        &mockUserRepo{anyExists: false},
			expectError: true,
		},
		{
			name:        "repo.AnyExists error",
			username:    "admin",
			repo:        &mockUserRepo{anyExistsErr: errors.New("db error")},
			expectError: true,
		},
		{
			name:        "repo.Create error",
			username:    "admin",
			repo:        &mockUserRepo{createErr: errors.New("insert failed")},
			expectError: true,
		},
		{
			name:          "success",
			username:      "admin",
			repo:          &mockUserRepo{},
			expectError:   false,
			expectMessage: "admin created",
			expectCreds:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creds, msg, err := CreateAdmin(tt.username, tt.repo)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if msg != tt.expectMessage {
				t.Fatalf("expected message %q, got %q", tt.expectMessage, msg)
			}
			if tt.expectCreds && creds == nil {
				t.Fatal("expected credentials, got nil")
			}
			if !tt.expectCreds && creds != nil {
				t.Fatal("expected no credentials, got non-nil")
			}
		})
	}
}
