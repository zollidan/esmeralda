package tasks

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/zollidan/esmeralda/mq"
	"github.com/zollidan/esmeralda/scheduler"
	"github.com/zollidan/esmeralda/utils"
)

type Manager struct {
	QueueMQ   *amqp091.Queue
	Scheduler *scheduler.Scheduler
	mq        *mq.MQ
	ch        *amqp091.Channel
}

func NewManager() *Manager {

	rabbitMQ := mq.NewMQ()
	sched := scheduler.NewScheduler()

	ch := rabbitMQ.GetChannel()

	q, err := ch.QueueDeclare(
		"test_queue", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)

	utils.FailOnError(err, "Failed to declare a queue")

	return &Manager{
		QueueMQ:   &q,
		Scheduler: sched,
		mq:        rabbitMQ,
		ch:        ch,
	}
}

func (tm *Manager) StartTask(taskID uint) error {
	// Send task ID to queue for processing
	taskMessage := fmt.Sprintf(`{"task_id": %d}`, taskID)
	err := tm.mq.SendMessage(tm.ch, tm.QueueMQ.Name, taskMessage)
	utils.FailOnError(err, "Error starting task.")
	return err
}

func (tm *Manager) StartTaskWithScheduler(taskID uint, cronExpr string) {
	// Implementation for starting a task with the scheduler can be added here
	_, err := tm.Scheduler.NewJob(cronExpr, taskID)
	utils.FailOnError(err, "Error creating cron task")
	taskMessage := fmt.Sprintf(`{"task_id": %d}`, taskID)
	err = tm.mq.SendMessage(tm.ch, tm.QueueMQ.Name, taskMessage)
	utils.FailOnError(err, "Error starting task.")
}

func (tm *Manager) Shutdown() {
	// Close channel first
	if tm.ch != nil {
		if err := tm.ch.Close(); err != nil {
			log.Printf("Failed to close RabbitMQ channel: %v", err)
		}
	}

	// Then close MQ connection
	if tm.mq != nil {
		tm.mq.CloseMQ()
	}

	// Finally shutdown scheduler
	if tm.Scheduler != nil {
		tm.Scheduler.Shutdown()
	}
}
