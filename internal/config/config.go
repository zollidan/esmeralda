package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	RedisAddr  string
	ServerPort string
	DatabaseDSN string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return &Config{
		RedisAddr:   getEnv("REDIS_ADDR", "redis:6379"),
		ServerPort:  getEnv("SERVER_PORT", ":8080"),
		DatabaseDSN: getEnv("DATABASE_DSN", "host=localhost user=postgres password=postgres dbname=esmeralda port=5432 sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return fallback
}
