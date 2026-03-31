package server

import (
	"testing"
)

func TestNewHandler(t *testing.T) {
	h := NewHandler(nil, nil, nil, nil, nil, nil, nil, nil)
	if h == nil {
		t.Fatal("NewHandler returned nil")
	}
	if h.pending == nil {
		t.Error("pending map is nil")
	}
}
