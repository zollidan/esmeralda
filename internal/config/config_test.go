package config

import (
	"net/url"
	"os"
	"testing"
)

func TestGetEnv_ReturnsValue(t *testing.T) {
	t.Setenv("TEST_KEY", "hello")
	if got := getEnv("TEST_KEY", "default"); got != "hello" {
		t.Errorf("expected 'hello', got %q", got)
	}
}

func TestGetEnv_ReturnsFallbackWhenMissing(t *testing.T) {
	err := os.Unsetenv("TEST_KEY_MISSING")

	if err != nil {
		t.Fatalf("failed to unset TEST_KEY_MISSING: %v", err)
	}

	if got := getEnv("TEST_KEY_MISSING", "default"); got != "default" {
		t.Errorf("expected 'default', got %q", got)
	}
}

func TestGetEnv_ReturnsFallbackWhenEmpty(t *testing.T) {
	t.Setenv("TEST_KEY_EMPTY", "")
	if got := getEnv("TEST_KEY_EMPTY", "fallback"); got != "fallback" {
		t.Errorf("expected 'fallback', got %q", got)
	}
}

func TestGetEnvInt_ReturnsValue(t *testing.T) {
	t.Setenv("TEST_INT", "42")
	if got := getEnvInt("TEST_INT", 10); got != 42 {
		t.Errorf("expected 42, got %d", got)
	}
}

func TestGetEnvInt_ReturnsFallbackWhenMissing(t *testing.T) {
	err := os.Unsetenv("TEST_INT_MISSING")
	if err != nil {
		t.Fatalf("failed to unset TEST_INT_MISSING: %v", err)
	}
	if got := getEnvInt("TEST_INT_MISSING", 99); got != 99 {
		t.Errorf("expected 99, got %d", got)
	}
}

func TestGetEnvInt_ReturnsFallbackOnInvalidValue(t *testing.T) {
	t.Setenv("TEST_INT_BAD", "notanumber")
	if got := getEnvInt("TEST_INT_BAD", 5); got != 5 {
		t.Errorf("expected 5, got %d", got)
	}
}

func TestLoad_DefaultValues(t *testing.T) {
	t.Setenv("SPORT_API_TOKEN", "test-token")
	err := os.Unsetenv("EXCEL_FILE_PATH")
	if err != nil {
		t.Fatalf("failed to unset EXCEL_FILE_PATH: %v", err)
	}
	err = os.Unsetenv("EXCEL_FILE_NAME")
	if err != nil {
		t.Fatalf("failed to unset EXCEL_FILE_NAME: %v", err)
	}
	err = os.Unsetenv("WORKERS")
	if err != nil {
		t.Fatalf("failed to unset WORKERS: %v", err)
	}
	err = os.Unsetenv("REDIS_ADDR")
	if err != nil {
		t.Fatalf("failed to unset REDIS_ADDR: %v", err)
	}
	err = os.Unsetenv("SERVER_PORT")
	if err != nil {
		t.Fatalf("failed to unset SERVER_PORT: %v", err)
	}
	err = os.Unsetenv("DATABASE_DSN")
	if err != nil {
		t.Fatalf("failed to unset DATABASE_DSN: %v", err)
	}

	cfg := Load()

	if cfg.SportAPIRU.Token != "test-token" {
		t.Errorf("expected token 'test-token', got %q", cfg.SportAPIRU.Token)
	}
	if cfg.SportAPIRU.BaseURL != "https://api.api-sport.ru/v2" {
		t.Errorf("unexpected BaseURL: %q", cfg.SportAPIRU.BaseURL)
	}
	if cfg.SportAPIRU.BaseFootballURL != "https://api.api-sport.ru/v2/football" {
		t.Errorf("unexpected BaseFootballURL: %q", cfg.SportAPIRU.BaseFootballURL)
	}
	if cfg.Tech.Workers != 20 {
		t.Errorf("expected default Workers 20, got %d", cfg.Tech.Workers)
	}
	if cfg.RedisAddr != "localhost:6379" {
		t.Errorf("expected default RedisAddr 'localhost:6379', got %q", cfg.RedisAddr)
	}
	if cfg.ServerPort != ":8080" {
		t.Errorf("expected default ServerPort ':8080', got %q", cfg.ServerPort)
	}
	parsedDBURL, err := url.Parse(cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("failed to parse DatabaseURL: %v", err)
	}
	if parsedDBURL.Scheme != "postgres" {
		t.Errorf("expected scheme 'postgres', got %q", parsedDBURL.Scheme)
	}
	if parsedDBURL.Host != "localhost:5432" {
		t.Errorf("expected host 'localhost:5432', got %q", parsedDBURL.Host)
	}
	if parsedDBURL.Path != "/postgres" {
		t.Errorf("expected path '/postgres', got %q", parsedDBURL.Path)
	}
	if parsedDBURL.Query().Get("sslmode") != "disable" {
		t.Errorf("expected sslmode 'disable', got %q", parsedDBURL.Query().Get("sslmode"))
	}
}

func TestLoad_CustomValues(t *testing.T) {
	t.Setenv("SPORT_API_TOKEN", "my-token")
	t.Setenv("EXCEL_FILE_PATH", "/tmp/output/")
	t.Setenv("EXCEL_FILE_NAME", "custom.xlsx")
	t.Setenv("WORKERS", "5")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("SERVER_PORT", ":9090")
	t.Setenv("DATABASE_URL", "postgres://test:test@localhost:5432/test?sslmode=disable")

	cfg := Load()

	if cfg.SportAPIRU.Token != "my-token" {
		t.Errorf("expected 'my-token', got %q", cfg.SportAPIRU.Token)
	}
	if cfg.Tech.Workers != 5 {
		t.Errorf("expected Workers 5, got %d", cfg.Tech.Workers)
	}
	if cfg.RedisAddr != "localhost:6379" {
		t.Errorf("expected 'localhost:6379', got %q", cfg.RedisAddr)
	}
	if cfg.ServerPort != ":9090" {
		t.Errorf("expected ':9090', got %q", cfg.ServerPort)
	}
	parsedDBURL, err := url.Parse(cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("failed to parse DatabaseURL: %v", err)
	}
	if parsedDBURL.Scheme != "postgres" {
		t.Errorf("expected scheme 'postgres', got %q", parsedDBURL.Scheme)
	}
	if parsedDBURL.Host != "localhost:5432" {
		t.Errorf("expected host 'localhost:5432', got %q", parsedDBURL.Host)
	}
	if parsedDBURL.Path != "/test" {
		t.Errorf("expected path '/test', got %q", parsedDBURL.Path)
	}
	if parsedDBURL.Query().Get("sslmode") != "disable" {
		t.Errorf("expected sslmode 'disable', got %q", parsedDBURL.Query().Get("sslmode"))
	}
}

func TestInitConfig_ReturnsConfig(t *testing.T) {
	t.Setenv("SPORT_API_TOKEN", "init-token")

	cfg := InitConfig()
	if cfg.SportAPIRU.Token != "init-token" {
		t.Errorf("expected 'init-token', got %q", cfg.SportAPIRU.Token)
	}
}
