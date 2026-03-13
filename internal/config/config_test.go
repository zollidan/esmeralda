package config

import (
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
	os.Unsetenv("TEST_KEY_MISSING")
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
	os.Unsetenv("TEST_INT_MISSING")
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
	os.Unsetenv("EXCEL_FILE_PATH")
	os.Unsetenv("EXCEL_FILE_NAME")
	os.Unsetenv("WORKERS")
	os.Unsetenv("REDIS_ADDR")
	os.Unsetenv("SERVER_PORT")
	os.Unsetenv("DATABASE_DSN")

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
	if cfg.Excel.FilePath != "./output/" {
		t.Errorf("expected default FilePath './output/', got %q", cfg.Excel.FilePath)
	}
	if cfg.Excel.FileName != "esmeralda-ru.xlsx" {
		t.Errorf("expected default FileName 'esmeralda-ru.xlsx', got %q", cfg.Excel.FileName)
	}
	if cfg.Tech.Workers != 20 {
		t.Errorf("expected default Workers 20, got %d", cfg.Tech.Workers)
	}
	if cfg.RedisAddr != "redis:6379" {
		t.Errorf("expected default RedisAddr 'redis:6379', got %q", cfg.RedisAddr)
	}
	if cfg.ServerPort != ":8080" {
		t.Errorf("expected default ServerPort ':8080', got %q", cfg.ServerPort)
	}
	if cfg.DatabaseDSN != "" {
		t.Errorf("expected empty default DatabaseDSN, got %q", cfg.DatabaseDSN)
	}
}

func TestLoad_CustomValues(t *testing.T) {
	t.Setenv("SPORT_API_TOKEN", "my-token")
	t.Setenv("EXCEL_FILE_PATH", "/tmp/output/")
	t.Setenv("EXCEL_FILE_NAME", "custom.xlsx")
	t.Setenv("WORKERS", "5")
	t.Setenv("REDIS_ADDR", "localhost:6379")
	t.Setenv("SERVER_PORT", ":9090")
	t.Setenv("DATABASE_DSN", "host=localhost user=test password=test dbname=test port=5432 sslmode=disable")

	cfg := Load()

	if cfg.SportAPIRU.Token != "my-token" {
		t.Errorf("expected 'my-token', got %q", cfg.SportAPIRU.Token)
	}
	if cfg.Excel.FilePath != "/tmp/output/" {
		t.Errorf("expected '/tmp/output/', got %q", cfg.Excel.FilePath)
	}
	if cfg.Excel.FileName != "custom.xlsx" {
		t.Errorf("expected 'custom.xlsx', got %q", cfg.Excel.FileName)
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
	if cfg.DatabaseDSN != "host=localhost user=test password=test dbname=test port=5432 sslmode=disable" {
		t.Errorf("expected custom DatabaseDSN, got %q", cfg.DatabaseDSN)
	}
}

func TestInitConfig_ReturnsConfig(t *testing.T) {
	t.Setenv("SPORT_API_TOKEN", "init-token")

	cfg := InitConfig()
	if cfg.SportAPIRU.Token != "init-token" {
		t.Errorf("expected 'init-token', got %q", cfg.SportAPIRU.Token)
	}
}
