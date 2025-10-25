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

func InitScheduler() *Scheduler {
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("failed to create scheduler: %v", err)
		return nil
	}

	return &Scheduler{Scheduler: s}
}

func (s *Scheduler) NewJob() (uuid.UUID, error) {
	// add a job to the scheduler
	j, err := s.Scheduler.NewJob(
		gocron.DurationJob(
			20*time.Second,
		),
		gocron.NewTask(
			func(a string, b int) {
				fmt.Println("Hello from gocron")
			},
			"hello",
			1,
		),
		gocron.WithStartAt(gocron.WithStartImmediately()),
	)

	if err != nil {
		return uuid.Nil, err
	}

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
