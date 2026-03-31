package utils

import (
	"strings"
	"testing"
	"time"
)

func TestInputDate(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		checkDate   func(t *testing.T, date time.Time)
	}{
		{
			name:        "valid date",
			input:       "01.12.2025\n",
			expectError: false,
			checkDate: func(t *testing.T, date time.Time) {
				expected := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
				if !date.Equal(expected) {
					t.Fatalf("expected %v, got %v", expected, date)
				}
			},
		},
		{
			name:        "empty input returns current date",
			input:       "\n",
			expectError: false,
			checkDate: func(t *testing.T, date time.Time) {
				if date.IsZero() {
					t.Fatal("expected current date, got zero")
				}
			},
		},
		{
			name:        "invalid format",
			input:       "2025-12-01\n",
			expectError: true,
			checkDate:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			date, err := InputDate(r)

			if tt.expectError {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.checkDate != nil {
				tt.checkDate(t, date)
			}
		})
	}
}
