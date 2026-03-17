package processor

import (
	"context"
	"log"
	"time"

	"github.com/zollidan/esmeralda/internal/api"
	"github.com/zollidan/esmeralda/internal/queue"
	"github.com/zollidan/esmeralda/internal/repository"
)

type Processor struct {
	apiClient        api.MatchFetcher
	games            *repository.GameRepository
	resultsProducer  *queue.Producer
	progressProducer *queue.Producer
	workers          int
}

func Init(apiClient api.MatchFetcher, games *repository.GameRepository, resultsProducer *queue.Producer, progressProducer *queue.Producer, workers int) *Processor {
	return &Processor{
		apiClient:        apiClient,
		games:            games,
		resultsProducer:  resultsProducer,
		progressProducer: progressProducer,
		workers:          workers,
	}
}

func (p *Processor) ProcessParseTask(ctx context.Context, payload []byte) error {
	task, err := queue.Unmarshal[queue.ParseTask](payload)
	if err != nil {
		return err
	}

	log.Printf("received task: id=%s date=%s", task.ID, task.Date)

	p.publishResult(ctx, task.ID, queue.StatusProcessing, "")

	date, err := time.Parse("2006-01-02", task.Date)
	if err != nil {
		return p.publishResult(ctx, task.ID, queue.StatusError, err.Error())
	}

	matches, totalMatches, err := p.apiClient.GetMatches(api.MatchesFilter{Date: date})
	if err != nil {
		return p.publishResult(ctx, task.ID, queue.StatusError, err.Error())
	}

	if err := p.ProcessMatches(ctx, p.apiClient, p.games, matches, totalMatches, task.ID); err != nil {
		return p.publishResult(ctx, task.ID, queue.StatusError, err.Error())
	}

	return p.publishResult(ctx, task.ID, queue.StatusDone, "")
}

func (p *Processor) publishResult(ctx context.Context, taskID string, status queue.Status, errMsg string) error {
	result := queue.TaskResult{
		TaskID: taskID,
		Status: status,
		Error:  errMsg,
	}
	_, err := p.resultsProducer.Publish(ctx, result)
	if err != nil {
		log.Printf("publish result for task %s: %v", taskID, err)
	}
	return err
}

func (p *Processor) publishProgress(ctx context.Context, taskID string, status queue.Status, totalMatches int, currentMatch int) error {
	progress := queue.TaskProgress{
		TaskID:       taskID,
		Status:       status,
		TotalMatches: totalMatches,
		CurrentMatch: currentMatch,
	}

	_, err := p.progressProducer.Publish(ctx, progress)
	if err != nil {
		log.Printf("publish progress for task %s: %v", taskID, err)
	}
	return err
}
