package processor

import (
	"context"
	"log"
	"time"

	"github.com/zollidan/esmeralda/internal/api"
	"github.com/zollidan/esmeralda/internal/queue"
	"gorm.io/gorm"
)

type Processor struct {
	apiClient       api.MatchFetcher
	database        *gorm.DB
	resultsProducer *queue.Producer
	enrichProducer  *queue.Producer
}

func Init(apiClient api.MatchFetcher, database *gorm.DB, resultsProducer, enrichProducer *queue.Producer) *Processor {
	return &Processor{
		apiClient:       apiClient,
		database:        database,
		resultsProducer: resultsProducer,
		enrichProducer:  enrichProducer,
	}
}

func (p *Processor) ProcessParseTask(ctx context.Context, payload []byte) error {

	task, err := queue.Unmarshal[queue.ParseTask](payload)
	if err != nil {
		return err
	}

	log.Printf("received task: id=%s date=%s", task.ID, task.Date)

	date, err := time.Parse("2006-01-02", task.Date)
	if err != nil {
		return p.publishResult(ctx, task.ID, queue.StatusError, err.Error())
	}

	matches, totalMatches, err := p.apiClient.GetMatches(api.MatchesFilter{Date: date})
	if err != nil {
		return p.publishResult(ctx, task.ID, queue.StatusError, err.Error())
	}

	if err := ProcessMatches(p.apiClient, p.database, matches, totalMatches); err != nil {
		return p.publishResult(ctx, task.ID, queue.StatusError, err.Error())
	}

	return p.publishResult(ctx, task.ID, queue.StatusDone, "")
}

func (p *Processor) ProcessEnrichTask(ctx context.Context, payload []byte) error {return nil}

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
