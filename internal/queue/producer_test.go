package queue

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProducer_Publish(t *testing.T) {
	rdb := setupRedis(t)
	ctx := context.Background()

	p := NewProducer(rdb, "pub-stream")

	type Event struct {
		Name string `json:"name"`
	}

	msgID, err := p.Publish(ctx, Event{Name: "hello"})
	require.NoError(t, err)
	assert.NotEmpty(t, msgID)

	msgs, err := rdb.XRange(ctx, "pub-stream", "-", "+").Result()
	require.NoError(t, err)
	require.Len(t, msgs, 1)

	raw, ok := msgs[0].Values["payload"].(string)
	require.True(t, ok)

	var evt Event
	require.NoError(t, json.Unmarshal([]byte(raw), &evt))
	assert.Equal(t, "hello", evt.Name)
}

func TestProducer_Publish_ReturnsMessageID(t *testing.T) {
	rdb := setupRedis(t)
	ctx := context.Background()

	p := NewProducer(rdb, "id-stream")

	id1, err := p.Publish(ctx, map[string]string{"n": "1"})
	require.NoError(t, err)

	id2, err := p.Publish(ctx, map[string]string{"n": "2"})
	require.NoError(t, err)

	assert.NotEqual(t, id1, id2, "each message must get a unique ID")

	msgs, err := rdb.XRange(ctx, "id-stream", "-", "+").Result()
	require.NoError(t, err)
	assert.Equal(t, id1, msgs[0].ID)
	assert.Equal(t, id2, msgs[1].ID)
}

func TestProducer_Publish_UnmarshalablePayload(t *testing.T) {
	rdb := setupRedis(t)
	ctx := context.Background()

	p := NewProducer(rdb, "bad-stream")

	_, err := p.Publish(ctx, make(chan int))
	assert.ErrorContains(t, err, "marshal payload")
}

func TestProducer_Delete(t *testing.T) {
	rdb := setupRedis(t)
	ctx := context.Background()

	p := NewProducer(rdb, "del-stream")

	msgID, err := p.Publish(ctx, map[string]string{"x": "1"})
	require.NoError(t, err)

	require.NoError(t, p.Delete(ctx, msgID))

	msgs, err := rdb.XRange(ctx, "del-stream", "-", "+").Result()
	require.NoError(t, err)
	assert.Empty(t, msgs, "stream must be empty after delete")
}

func TestProducer_Delete_NonExistentID(t *testing.T) {
	rdb := setupRedis(t)
	ctx := context.Background()

	p := NewProducer(rdb, "ghost-stream")

	err := p.Delete(ctx, "0-0")
	assert.NoError(t, err)
}
