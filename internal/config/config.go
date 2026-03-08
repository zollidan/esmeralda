package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	SportAPIRU  SportAPIRU
	Excel       Excel
	Tech        Tech
	GRPC        GRPC
	RedisAddr   string
	ServerPort  string
	DatabaseDSN string
	Env         string
}

type SportAPIRU struct {
	BaseURL         string
	BaseFootballURL string
	Token           string
}

type Excel struct {
	FilePath string
	FileName string
}

type Tech struct {
	GamesLimit int
	Workers    int
}

type GRPC struct {
	Port    int
	Timeout time.Duration
}

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
			Workers: getEnvInt("WORKERS", 20),
		},
		GRPC: GRPC{
			Port:    getEnvInt("GRPC_PORT", 44044),
			Timeout: time.Duration(getEnvInt("GRPC_TIMEOUT", 10)),
		},
		RedisAddr:   getEnv("REDIS_ADDR", "redis:6379"),
		ServerPort:  getEnv("SERVER_PORT", ":8080"),
		DatabaseDSN: getEnv("DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=esmeralda port=5432 sslmode=disable"),
		Env:         getEnv("ENV", "local"),
	}
}

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
