package utils

import "testing"

func TestGeneratePassword(t *testing.T) {
	tests := []struct {
		name        string
		length      int
		expectError bool
	}{
		{"valid length 12", 12, false},
		{"negative length", -1, true},
		{"zero length", 0, true},
		{"boundary length 40", 40, false},
		{"boundary length 41", 41, true},
		{"large length 1000", 1000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := GeneratePassword(tt.length)

			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected error for length %d, got nil", tt.length)
				}
				return
			}

			if err != nil {
				t.Fatalf("Expected no error for length %d, got %v", tt.length, err)
			}
			if len(password) != tt.length {
				t.Fatalf("Expected password length %d, got %d", tt.length, len(password))
			}
		})
	}
}
