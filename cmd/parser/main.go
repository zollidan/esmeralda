package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/zollidan/esmeralda/internal/api"
	"github.com/zollidan/esmeralda/internal/config"
	"github.com/zollidan/esmeralda/internal/db"
	"github.com/zollidan/esmeralda/internal/queue"
)

func main() {
	cfg := config.InitConfig()

	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddr,
	})

	client := api.NewClient(cfg.SportAPIRU.BaseURL, cfg.SportAPIRU.Token)
	_, err := db.New(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}

	parseConsumer := queue.NewConsumer(rdb, queue.StreamParse)
	resultsProducer := queue.NewProducer(rdb, queue.StreamResults)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	err = parseConsumer.Consume(ctx, func(ctx context.Context, payload []byte) error {
		task, err := queue.Unmarshal[queue.ParseTask](payload)
		if err != nil {
			return err
		}

		log.Printf("received task: id=%s date=%s", task.ID, task.Date)

		date, err := time.Parse("2006-01-02", task.Date)
		if err != nil {
			return publishResult(ctx, resultsProducer, task.ID, queue.StatusError, err.Error())
		}

		_, totalMatches, err := client.GetMatches(api.MatchesFilter{Date: date})
		if err != nil {
			return publishResult(ctx, resultsProducer, task.ID, queue.StatusError, err.Error())
		}

		log.Printf("found %d matches for date %s", totalMatches, task.Date)

		// if err := processor.ProcessMatches(client, database, matches, totalMatches); err != nil {
		// 	return publishResult(ctx, resultsProducer, task.ID, queue.StatusError, err.Error())
		// }

		return publishResult(ctx, resultsProducer, task.ID, queue.StatusDone, "")
	})

	if err != nil && err != context.Canceled {
		log.Fatal(err)
	}
}

func publishResult(ctx context.Context, p *queue.Producer, taskID string, status queue.Status, errMsg string) error {
	result := queue.TaskResult{
		TaskID: taskID,
		Status: status,
		Error:  errMsg,
	}
	_, err := p.Publish(ctx, result)
	if err != nil {
		log.Printf("publish result for task %s: %v", taskID, err)
	}
	return err
}
