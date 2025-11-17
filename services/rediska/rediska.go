package rediska

import (
	"log"

	"github.com/go-redis/redis"
)

type RedisClient struct {
	Client *redis.Client
}

func Init(redisURL string) *RedisClient {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal("Error connecting to Redis.")
	}

	return &RedisClient{
		Client: redis.NewClient(opt),
	}
}

func (c *RedisClient) SendTask(taskID uint, parserName string, fileName string) (string, error) {
	result, err := c.Client.XAdd(&redis.XAddArgs{
		Stream: "parser_tasks",
		Values: map[string]interface{}{"task": taskID, "parser": parserName, "fileName": fileName},
	}).Result()

	return result, err
}

func (c *RedisClient) GetParsers() ([]string, error) {

	result, err := c.Client.SMembers("available_parsers").Result()
	if err != nil {
		return []string{}, err
	}

	return result, nil
}
