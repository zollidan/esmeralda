package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/zollidan/esmeralda-ru-api-fetcher/internal/queue"
	"github.com/zollidan/esmeralda-ru-api-fetcher/internal/server"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr: "redis:6379", 
	})

	producer := queue.NewProducer(rdb)
	handler := server.NewHandler(producer)

	r := gin.Default()
	server.SetupRoutes(r, handler)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}