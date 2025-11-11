package mq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/zollidan/esmeralda/utils"
)

// MQ represents a minimal RabbitMQ client that holds a single connection.
// Callers obtain channels from MQ and must close channels when done.
type MQ struct {
	conn *amqp.Connection
}

// NewMQ dials the RabbitMQ server at the default URL and returns an MQ.
// It panics on failure (via failOnError). Caller should call CloseMQ when finished.
func NewMQ() *MQ {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")

	return &MQ{conn: conn}
}

// GetChannel opens and returns a new AMQP channel from the underlying connection.
// The caller is responsible for closing the returned channel.
func (mq *MQ) GetChannel() *amqp.Channel {
	ch, err := mq.conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	return ch
}

// DeclareQueue declares the given queue as durable.
// The function returns an error if the declaration fails.
func (mq *MQ) DeclareQueue(ch *amqp.Channel, queueName string) (amqp.Queue, error) {
	if ch == nil {
		log.Panic("DeclareQueue: channel is nil")
	}

	queue, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	return queue, err
}

// SendMessage publishes body to the given queue.
// It uses a short context timeout for the publish call and makes the message persistent.
// The function returns an error if the publish fails.
func (mq *MQ) SendMessage(ch *amqp.Channel, queueName string, body string) error {
	if ch == nil {
		log.Panic("SendMessage: channel is nil")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := ch.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(body),
			DeliveryMode: amqp.Persistent, // make message persistent
		})
	if err != nil {
		return err
	}
	return nil
}

func (mq *MQ) RefreshParsersRPC(queueName string) (string, error) {
	ch := mq.GetChannel()
	defer ch.Close()

	// Создаём временную очередь для ответа
	replyQueue, err := ch.QueueDeclare(
		"",    // name (пустое = случайное имя)
		false, // durable
		true,  // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return "", fmt.Errorf("failed to declare reply queue: %w", err)
	}

	// Слушаем ответы
	msgs, err := ch.Consume(
		replyQueue.Name, // queue
		"",              // consumer
		true,            // auto-ack
		false,           // exclusive
		false,           // no-local
		false,           // no-wait
		nil,             // args
	)
	if err != nil {
		return "", fmt.Errorf("failed to register consumer: %w", err)
	}

	// Генерируем correlation ID
	corrId := generateUUID()

	// Отправляем запрос
	refreshMessage := map[string]string{
		"action": "refresh_parsers",
	}

	body, err := json.Marshal(refreshMessage)
	if err != nil {
		return "", fmt.Errorf("failed to marshal message: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,
			ReplyTo:       replyQueue.Name,
			Body:          body,
		})
	if err != nil {
		return "", fmt.Errorf("failed to publish message: %w", err)
	}

	// Ждём ответ
	for d := range msgs {
		if corrId == d.CorrelationId {
			return string(d.Body), nil
		}
	}

	return "", fmt.Errorf("timeout waiting for response")
}

func generateUUID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// CloseMQ closes the underlying RabbitMQ connection and logs any error.
// It is safe to call CloseMQ multiple times.
func (mq *MQ) CloseMQ() {
	if mq.conn == nil {
		return
	}
	if err := mq.conn.Close(); err != nil {
		log.Printf("Failed to close RabbitMQ connection: %v", err)
	}
}
