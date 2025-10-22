package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string
	DBURL         string
}

func New() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️ Warning: .env file not found, using system environment variables")
	}

	return Config{
		ServerAddress: getEnv("SERVER_ADDRESS", "localhost:3333"),
		DBURL:         getEnv("DB_URL", "test.db"),
	}
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
