package rediska

import (
	"log"

	"github.com/go-redis/redis"
)

type RedisClient struct {
	Client *redis.Client
}

var redisClient *RedisClient

func Init(redisURL string) *RedisClient{
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal("Error connecting to Redis.")
	}

	return &RedisClient{
		Client: redis.NewClient(opt),
	}
}
