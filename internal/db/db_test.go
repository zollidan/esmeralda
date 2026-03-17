package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/zollidan/esmeralda/internal/config"
)

func TestNew(t *testing.T) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "postgres",
			"POSTGRES_PASSWORD": "postgres",
			"POSTGRES_DB":       "test_db",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	assert.NoError(t, err)

	endpoint, err := container.Endpoint(ctx, "")

	cfg := &config.Config{
		DatabaseURL: "postgres://postgres:postgres@" + endpoint + "/test_db?sslmode=disable",
	}

	db, err := New(cfg)
	if err != nil {
		t.Fatalf("db init failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("failed to get sql db: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("db ping failed: %v", err)
	}
}
