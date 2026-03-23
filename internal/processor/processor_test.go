package processor

import "testing"

func TestSetGamesLimit(t *testing.T) {
	tests := []struct {
		name       string
		matches    int
		gamesLimit int
		want       int
	}{
		{"no limit", 10, 0, 10},
		{"limit less than matches", 10, 5, 5},
		{"limit greater than matches", 10, 20, 10},
		{"limit equals matches", 10, 10, 10},
		{"zero matches", 0, 5, 0},
		{"both zero", 0, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := setGamesLimit(tt.matches, tt.gamesLimit)
			if got != tt.want {
				t.Errorf("setGamesLimit(%d, %d) = %d, want %d", tt.matches, tt.gamesLimit, got, tt.want)
			}
		})
	}
}
