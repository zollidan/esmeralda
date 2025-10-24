package mq

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// MQ represents a minimal RabbitMQ client that holds a single connection.
// Callers obtain channels from MQ and must close channels when done.
type MQ struct {
	conn *amqp.Connection
}

// failOnError logs the provided message and panics if err is non-nil.
// This is a helper used in this package for simplicity.
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// NewMQ dials the RabbitMQ server at the default URL and returns an MQ.
// It panics on failure (via failOnError). Caller should call CloseMQ when finished.
func NewMQ() *MQ {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")

	return &MQ{conn: conn}
}

// GetChannel opens and returns a new AMQP channel from the underlying connection.
// The caller is responsible for closing the returned channel.
func (mq *MQ) GetChannel() *amqp.Channel {
	ch, err := mq.conn.Channel()
	failOnError(err, "Failed to open a channel")
	return ch
}

// SendMessage declares the given queue (durable) and publishes body to it.
// It uses a short context timeout for the publish call and makes the message persistent.
// The function panics on error.
func (mq *MQ) SendMessage(ch *amqp.Channel, queueName string, body string) error {
	if ch == nil {
		log.Panic("SendMessage: channel is nil")
	}

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable -> changed to true for durability
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(body),
			DeliveryMode: amqp.Persistent, // make message persistent
		})
	if err != nil {
		return err
	}
	log.Printf(" [x] Sent %s\n", body)
	return nil
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