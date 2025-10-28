package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerAddress string
	DBURL         string
	S3URL 		  string
	BucketName    string
	S3AccessKey   string
	S3SecretKey   string
}

func New() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("\n⚠️ Warning: .env file not found, using system environment variables")
	}

	return Config{
		ServerAddress: getEnv("SERVER_ADDRESS", "localhost:3333"),
		DBURL:         getEnv("DB_URL", "test.db"),
		S3URL:         getEnv("S3URL", ""),
		BucketName:    getEnv("BUCKET_NAME", "my-bucket"),
		S3AccessKey:   getEnv("S3_ACCESS_KEY", ""),
		S3SecretKey:   getEnv("S3_SECRET_KEY", ""),
	}
}

func getEnv(key, defaultVal string) string {
	if val, exists := os.LookupEnv(key); exists {
		return val
	}
	return defaultVal
}
