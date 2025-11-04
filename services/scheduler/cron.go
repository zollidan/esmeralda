package scheduler

import (
	"fmt"
	"log"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
)

type Scheduler struct {
	Scheduler gocron.Scheduler
}

func NewScheduler() *Scheduler {
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("failed to create scheduler: %v", err)
		return nil
	}

	return &Scheduler{Scheduler: s}
}

func (s *Scheduler) NewJob(cronExpr string, taskID uint) (uuid.UUID, error) {
	// add a job to the scheduler
	j, err := s.Scheduler.NewJob(
		gocron.CronJob(
			cronExpr,
			false,
		),
		gocron.NewTask(
			func() {
				fmt.Printf("Hello from gocron %s %d\n", time.Now().Format(time.RFC3339), taskID)
			},
		),
		gocron.WithName(fmt.Sprintf("task-%d", taskID)),
		gocron.WithStartAt(
			gocron.WithStartImmediately(),
		),
	)

	if err != nil {
		return uuid.Nil, err
	}

	s.Start()

	return j.ID(), nil
}

func (s *Scheduler) Start() {
	s.Scheduler.Start()
}

func (s *Scheduler) Shutdown() {
	// when you're done, shut it down
	err := s.Scheduler.Shutdown()
	if err != nil {
		log.Printf("error shutting down scheduler: %v", err)
	}
}
