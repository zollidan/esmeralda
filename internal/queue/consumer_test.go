package queue

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupRedis(t *testing.T) *redis.Client {
	t.Helper()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "redis:latest",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	endpoint, err := container.Endpoint(ctx, "")
	require.NoError(t, err)

	rdb := redis.NewClient(&redis.Options{Addr: endpoint})
	require.NoError(t, rdb.Ping(ctx).Err())

	return rdb
}

func publishMessage(t *testing.T, rdb *redis.Client, stream string, payload any) {
	t.Helper()
	raw, err := json.Marshal(payload)
	require.NoError(t, err)

	err = rdb.XAdd(context.Background(), &redis.XAddArgs{
		Stream: stream,
		Values: map[string]any{"payload": string(raw)},
	}).Err()
	require.NoError(t, err)
}

func TestConsumer_Consume_HandlerCalled(t *testing.T) {
	rdb := setupRedis(t)

	const stream, group, consumer = "test-stream", "test-group", "worker-1"

	type Event struct {
		Name string `json:"name"`
	}

	publishMessage(t, rdb, stream, Event{Name: "hello"})

	received := make(chan []byte, 1)
	handler := func(ctx context.Context, payload []byte) error {
		received <- payload
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c := NewConsumer(rdb, stream, group, consumer)

	go func() {
		_ = c.Consume(ctx, handler)
	}()

	select {
	case payload := <-received:
		var evt Event
		require.NoError(t, json.Unmarshal(payload, &evt))
		assert.Equal(t, "hello", evt.Name)
	case <-ctx.Done():
		t.Fatal("timeout: handler was not called")
	}
}

func TestConsumer_Consume_AckAndDelete(t *testing.T) {
	rdb := setupRedis(t)

	const stream, group, consumer = "ack-stream", "ack-group", "worker-1"

	publishMessage(t, rdb, stream, map[string]string{"key": "value"})

	done := make(chan struct{})
	handler := func(ctx context.Context, payload []byte) error {
		close(done)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c := NewConsumer(rdb, stream, group, consumer)
	go func() { _ = c.Consume(ctx, handler) }()

	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("timeout: message not processed")
	}

	// Даём время на XAck + XDel
	time.Sleep(100 * time.Millisecond)

	// PEL должен быть пуст
	pending, err := rdb.XPending(context.Background(), stream, group).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(0), pending.Count, "PEL should be empty after ack")

	// Стрим должен быть пуст после XDel
	msgs, err := rdb.XRange(context.Background(), stream, "-", "+").Result()
	require.NoError(t, err)
	assert.Empty(t, msgs, "stream should be empty after xdel")
}

func TestConsumer_Consume_HandlerError_StillAcks(t *testing.T) {
	rdb := setupRedis(t)

	const stream, group, consumer = "err-stream", "err-group", "worker-1"

	publishMessage(t, rdb, stream, map[string]string{"x": "1"})

	done := make(chan struct{})
	handler := func(ctx context.Context, payload []byte) error {
		defer close(done)
		return assert.AnError // хендлер падает
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c := NewConsumer(rdb, stream, group, consumer)
	go func() { _ = c.Consume(ctx, handler) }()

	select {
	case <-done:
	case <-ctx.Done():
		t.Fatal("timeout")
	}

	time.Sleep(100 * time.Millisecond)

	pending, err := rdb.XPending(context.Background(), stream, group).Result()
	require.NoError(t, err)
	assert.Equal(t, int64(0), pending.Count, "message must be acked even on handler error")
}

func TestUnmarshal(t *testing.T) {
	type Event struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	raw, _ := json.Marshal(Event{ID: 42, Name: "test"})

	evt, err := Unmarshal[Event](raw)
	require.NoError(t, err)
	assert.Equal(t, 42, evt.ID)
	assert.Equal(t, "test", evt.Name)
}
