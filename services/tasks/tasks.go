package tasks

import (
	"context"
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
	"github.com/zollidan/esmeralda/mq"
	"github.com/zollidan/esmeralda/scheduler"
	"github.com/zollidan/esmeralda/utils"
)

type TaskManager interface {
	StartTask(taskID uint) error
	StartTaskWithScheduler(taskID uint, cronExpr string)
	GetResult() error
	CancelTask(taskID uint) error
	Shutdown()
}

type Manager struct {
	TaskQueue   *amqp091.Queue
	ResultQueue *amqp091.Queue
	Scheduler   *scheduler.Scheduler
	mq          *mq.MQ
	ch          *amqp091.Channel
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewManager() *Manager {

	rabbitMQ := mq.NewMQ()
	sched := scheduler.NewScheduler()

	ch := rabbitMQ.GetChannel()

	taskQueue, err := ch.QueueDeclare(
		"test_queue", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")

	resultQueue, err := ch.QueueDeclare(
		"test_result_queue", // name
		false,               // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	utils.FailOnError(err, "Failed to declare a result queue")

	ctx, cancel := context.WithCancel(context.Background())

	return &Manager{
		TaskQueue:   &taskQueue,
		ResultQueue: &resultQueue,
		Scheduler:   sched,
		mq:          rabbitMQ,
		ch:          ch,
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (tm *Manager) StartTask(taskID uint, parserID uint) error {
	// Send task ID to queue for processing
	taskMessage := fmt.Sprintf(`{"task_id": %d, "parser_id": %d}`, taskID, parserID)
	err := tm.mq.SendMessage(tm.ch, tm.TaskQueue.Name, taskMessage)
	utils.FailOnError(err, "Error starting task.")
	return err
}

func (tm *Manager) StartTaskWithScheduler(taskID uint, cronExpr string) {
	// Implementation for starting a task with the scheduler can be added here
	_, err := tm.Scheduler.NewJob(cronExpr, taskID)
	utils.FailOnError(err, "Error creating cron task")
	taskMessage := fmt.Sprintf(`{"task_id": %d}`, taskID)
	err = tm.mq.SendMessage(tm.ch, tm.TaskQueue.Name, taskMessage)
	utils.FailOnError(err, "Error starting task.")
}

// GetResult получает сообщения из очереди результатов и выводит в консоль
func (tm *Manager) GetResult() error {
	msgs, err := tm.ch.Consume(
		tm.ResultQueue.Name, // queue
		"",                  // consumer
		false,               // auto-ack - выключен для надёжности
		false,               // exclusive
		false,               // no-local
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		return fmt.Errorf("failed to register a consumer: %w", err)
	}

	log.Printf("Starting result consumer for queue: %s", tm.ResultQueue.Name)

	go tm.processMessages(msgs)
	return nil
}

// processMessages обрабатывает входящие сообщения
func (tm *Manager) processMessages(msgs <-chan amqp091.Delivery) {
	for {
		select {
		case <-tm.ctx.Done():
			log.Println("Stopping result consumer...")
			return
		case msg, ok := <-msgs:
			if !ok {
				log.Println("Message channel closed")
				return
			}

			// Обработка сообщения
			log.Printf("Received result: %s", string(msg.Body))

			// Acknowledge сообщение после обработки
			if err := msg.Ack(false); err != nil {
				log.Printf("Failed to acknowledge message: %v", err)
			}
		}
	}
}

func (tm *Manager) CancelTask(taskID uint) error {
	// Send task ID to queue for processing
	taskMessage := fmt.Sprintf(`{"task_id": %d, "action": "cancel"}`, taskID)
	err := tm.mq.SendMessage(tm.ch, tm.TaskQueue.Name, taskMessage)
	utils.FailOnError(err, "Error canceling task.")
	return err
}

func (tm *Manager) Shutdown() {
	log.Println("Shutting down Task Manager...")

	// Cancel context to stop consumers
	if tm.cancel != nil {
		tm.cancel()
	}

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

	log.Println("Task Manager shut down successfully")
}
