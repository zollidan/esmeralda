package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config is the main configuration struct for the application. It contains all necessary configuration fields, grouped by their respective domains (e.g., SportAPIRU, Excel, Tech, TelegramBot).
type Config struct {
	SportAPIRU  SportAPIRU
	Excel       Excel
	Tech        Tech
	TelegramBot TelegramBot
	DatabaseURL string
	RedisAddr   string
	ServerPort  string
	Env         string
}

// SportAPIRU contains configuration related to the Sport API.
type SportAPIRU struct {
	BaseURL         string
	BaseFootballURL string
	Token           string
}

// Excel contains configuration related to Excel file handling.
type Excel struct {
	FilePath string
	FileName string
}

// Tech contains configuration related to technical aspects of the application.
type Tech struct {
	ParserPort string
	GamesLimit int
	Workers    int
}

// TelegramBot contains configuration related to the Telegram Bot integration.
type TelegramBot struct {
	Token      string
	Debug      string
	APIBaseURL string
	APIAuth    string
}

// Load reads the configuration from environment variables (and optionally from a .env file) and returns a Config struct populated with the values. It also performs basic validation to ensure required fields are set.
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	token := getEnv("SPORT_API_TOKEN", "")
	if token == "" {
		log.Fatal("SPORT_API_TOKEN is required in environment")
	}

	return &Config{
		SportAPIRU: SportAPIRU{
			BaseURL:         "https://api.api-sport.ru/v2",
			BaseFootballURL: "https://api.api-sport.ru/v2/football",
			Token:           token,
		},
		Excel: Excel{
			FilePath: getEnv("EXCEL_FILE_PATH", "./output/"),
			FileName: getEnv("EXCEL_FILE_NAME", "esmeralda-ru.xlsx"),
		},
		Tech: Tech{
			ParserPort: getEnv("PARSER_PORT", ":8081"),
			Workers:    getEnvInt("WORKERS", 20),
			GamesLimit: getEnvInt("GAMES_LIMIT", 20),
		},
		TelegramBot: TelegramBot{
			Token:      getEnv("TELEGRAM_BOT_TOKEN", ""),
			Debug:      getEnv("TELEGRAM_BOT_DEBUG", "true"),
			APIBaseURL: getEnv("API_BASE_URL", "http://localhost:8080"),
			APIAuth:    getEnv("API_TOKEN", ""),
		},
		DatabaseURL: getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		ServerPort:  getEnv("SERVER_PORT", ":8080"),
		Env:         getEnv("ENV", "local"),
	}
}

// InitConfig is a helper function that initializes the configuration by calling Load and returns the Config struct. This function can be used in the main application to easily get the configuration.
func InitConfig() Config {
	return *Load()
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	valStr := getEnv(key, "")
	if valStr == "" {
		return fallback
	}
	if v, err := strconv.Atoi(valStr); err == nil {
		return v
	}
	return fallback
}
