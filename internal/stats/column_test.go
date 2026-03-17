package stats

import (
	"testing"

	"github.com/zollidan/esmeralda/internal/models"
)

func TestHeaders(t *testing.T) {
	headers := Headers()

	if len(headers) != len(Columns) {
		t.Fatalf("headers len = %d, Columns len = %d", len(headers), len(Columns))
	}

	// First header should be "число"
	if headers[0] != "число" {
		t.Errorf("headers[0] = %v, want число", headers[0])
	}

	// Verify all headers are non-empty
	for i, h := range headers {
		if h == "" {
			t.Errorf("headers[%d] is empty", i)
		}
	}
}

func TestGameToSlice(t *testing.T) {
	g := &models.Game{
		Day:            15,
		Month:          3,
		Year:           2025,
		Time:           "18:00",
		HomeTeam:       "Team A",
		AwayTeam:       "Team B",
		League:         "League 1",
		H2HHomeMatches: 10,
		H2HHomeWin1:    5,
		H2HHomeDraw:    3,
		H2HHomeWin2:    2,
	}

	slice := GameToSlice(g)

	if len(slice) != len(Columns) {
		t.Fatalf("slice len = %d, want %d", len(slice), len(Columns))
	}

	// First 7 basic fields
	if slice[0] != 15 {
		t.Errorf("slice[0] (Day) = %v, want 15", slice[0])
	}
	if slice[1] != 3 {
		t.Errorf("slice[1] (Month) = %v, want 3", slice[1])
	}
	if slice[2] != 2025 {
		t.Errorf("slice[2] (Year) = %v, want 2025", slice[2])
	}
	if slice[3] != "18:00" {
		t.Errorf("slice[3] (Time) = %v, want '18:00'", slice[3])
	}
	if slice[4] != "Team A" {
		t.Errorf("slice[4] (HomeTeam) = %v, want 'Team A'", slice[4])
	}
	if slice[5] != "Team B" {
		t.Errorf("slice[5] (AwayTeam) = %v, want 'Team B'", slice[5])
	}
	if slice[6] != "League 1" {
		t.Errorf("slice[6] (League) = %v, want 'League 1'", slice[6])
	}

	// H2H home stats
	if slice[7] != 10 {
		t.Errorf("slice[7] (H2H home matches) = %v, want 10", slice[7])
	}
	if slice[8] != 5 {
		t.Errorf("slice[8] (H2H home win1) = %v, want 5", slice[8])
	}
}
