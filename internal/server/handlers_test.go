package server

import (
	"sync"
	"testing"

	"github.com/zollidan/esmeralda/internal/queue"
)

func TestSubscribeProgress(t *testing.T) {
	h := &Handler{
		progressSub: make(map[string]map[chan queue.TaskProgress]struct{}),
	}

	ch := h.subscribeProgress("task-1")
	if ch == nil {
		t.Fatal("subscribeProgress returned nil")
	}

	if len(h.progressSub["task-1"]) != 1 {
		t.Errorf("subscribers = %d, want 1", len(h.progressSub["task-1"]))
	}
}

func TestSubscribeProgress_MultipleSubs(t *testing.T) {
	h := &Handler{
		progressSub: make(map[string]map[chan queue.TaskProgress]struct{}),
	}

	ch1 := h.subscribeProgress("task-1")
	ch2 := h.subscribeProgress("task-1")
	_ = h.subscribeProgress("task-2")

	if len(h.progressSub["task-1"]) != 2 {
		t.Errorf("task-1 subscribers = %d, want 2", len(h.progressSub["task-1"]))
	}
	if len(h.progressSub["task-2"]) != 1 {
		t.Errorf("task-2 subscribers = %d, want 1", len(h.progressSub["task-2"]))
	}
	_ = ch1
	_ = ch2
}

func TestUnsubscribeProgress(t *testing.T) {
	h := &Handler{
		progressSub: make(map[string]map[chan queue.TaskProgress]struct{}),
	}

	ch := h.subscribeProgress("task-1")
	h.unsubscribeProgress("task-1", ch)

	if _, exists := h.progressSub["task-1"]; exists {
		t.Error("task-1 entry should be cleaned up")
	}
}

func TestUnsubscribeProgress_PartialCleanup(t *testing.T) {
	h := &Handler{
		progressSub: make(map[string]map[chan queue.TaskProgress]struct{}),
	}

	ch1 := h.subscribeProgress("task-1")
	_ = h.subscribeProgress("task-1")

	h.unsubscribeProgress("task-1", ch1)

	if len(h.progressSub["task-1"]) != 1 {
		t.Errorf("subscribers = %d, want 1", len(h.progressSub["task-1"]))
	}
}

func TestBroadcastProgress(t *testing.T) {
	h := &Handler{
		progressSub: make(map[string]map[chan queue.TaskProgress]struct{}),
	}

	ch := h.subscribeProgress("task-1")

	p := queue.TaskProgress{
		TaskID:       "task-1",
		Status:       queue.StatusProcessing,
		TotalMatches: 20,
		CurrentMatch: 5,
	}

	h.broadcastProgress(p)

	select {
	case received := <-ch:
		if received.TaskID != "task-1" {
			t.Errorf("TaskID = %q, want %q", received.TaskID, "task-1")
		}
		if received.CurrentMatch != 5 {
			t.Errorf("CurrentMatch = %d, want 5", received.CurrentMatch)
		}
	default:
		t.Fatal("expected to receive progress on channel")
	}
}

func TestBroadcastProgress_NoSubscribers(t *testing.T) {
	h := &Handler{
		progressSub: make(map[string]map[chan queue.TaskProgress]struct{}),
	}

	// Should not panic with no subscribers
	p := queue.TaskProgress{TaskID: "task-1"}
	h.broadcastProgress(p)
}

func TestBroadcastProgress_MultipleSubscribers(t *testing.T) {
	h := &Handler{
		progressSub: make(map[string]map[chan queue.TaskProgress]struct{}),
	}

	ch1 := h.subscribeProgress("task-1")
	ch2 := h.subscribeProgress("task-1")

	p := queue.TaskProgress{
		TaskID:       "task-1",
		CurrentMatch: 10,
	}

	h.broadcastProgress(p)

	// Both channels should receive the progress
	for i, ch := range []chan queue.TaskProgress{ch1, ch2} {
		select {
		case received := <-ch:
			if received.CurrentMatch != 10 {
				t.Errorf("ch%d: CurrentMatch = %d, want 10", i+1, received.CurrentMatch)
			}
		default:
			t.Errorf("ch%d: expected to receive progress", i+1)
		}
	}
}

func TestBroadcastProgress_FullChannel(t *testing.T) {
	h := &Handler{
		progressSub: make(map[string]map[chan queue.TaskProgress]struct{}),
	}

	ch := h.subscribeProgress("task-1")

	// Fill the channel buffer (64)
	for i := 0; i < 64; i++ {
		h.broadcastProgress(queue.TaskProgress{TaskID: "task-1", CurrentMatch: i})
	}

	// Sending one more should not block (non-blocking select in broadcastProgress)
	h.broadcastProgress(queue.TaskProgress{TaskID: "task-1", CurrentMatch: 64})

	// Drain and verify
	received := <-ch
	if received.CurrentMatch != 0 {
		t.Errorf("CurrentMatch = %d, want 0", received.CurrentMatch)
	}
}

func TestBroadcastProgress_ConcurrentSafe(t *testing.T) {
	h := &Handler{
		progressSub: make(map[string]map[chan queue.TaskProgress]struct{}),
	}

	var wg sync.WaitGroup

	// Concurrently subscribe, broadcast, and unsubscribe
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch := h.subscribeProgress("task-1")
			h.broadcastProgress(queue.TaskProgress{TaskID: "task-1"})
			h.unsubscribeProgress("task-1", ch)
		}()
	}

	wg.Wait()
}

func TestNewHandler(t *testing.T) {
	h := NewHandler(nil, nil, nil, nil)
	if h == nil {
		t.Fatal("NewHandler returned nil")
	}
	if h.pending == nil {
		t.Error("pending map is nil")
	}
	if h.progressSub == nil {
		t.Error("progressSub map is nil")
	}
}
