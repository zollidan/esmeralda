package utils

import (
	"strings"
	"testing"
	"time"
)

func TestInputDate_Valid(t *testing.T) {
	input := "01.12.2025\n"

	r := strings.NewReader(input)

	date, err := InputDate(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)

	if !date.Equal(expected) {
		t.Fatalf("expected %v, got %v", expected, date)
	}
}

func TestInputDate_Empty(t *testing.T) {
	r := strings.NewReader("\n")

	date, err := InputDate(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if date.IsZero() {
		t.Fatal("expected current date, got zero")
	}
}

func TestInputDate_Invalid(t *testing.T) {
	r := strings.NewReader("2025-12-01\n")

	_, err := InputDate(r)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
