package server

import (
	"testing"

	"github.com/zollidan/esmeralda/internal/models"
	"github.com/zollidan/esmeralda/internal/stats"
)

func TestGameToSlice_BasicFields(t *testing.T) {
	g := models.Game{
		Day:            5,
		Month:          6,
		Year:           2025,
		Time:           "20:00",
		HomeTeam:       "Team1",
		AwayTeam:       "Team2",
		League:         "Liga",
		H2HHomeMatches: 10,
		H2HHomeWin1:    5,
		H2HHomeDraw:    3,
		H2HHomeWin2:    2,
	}

	slice := stats.GameToSlice(&g)

	if len(slice) != len(stats.Columns) {
		t.Errorf("slice len = %d, want %d", len(slice), len(stats.Columns))
	}

	// First 7 values should be basic info
	if slice[0] != 5 {
		t.Errorf("slice[0] (Day) = %v, want 5", slice[0])
	}
	if slice[1] != 6 {
		t.Errorf("slice[1] (Month) = %v, want 6", slice[1])
	}
	if slice[2] != 2025 {
		t.Errorf("slice[2] (Year) = %v, want 2025", slice[2])
	}
	if slice[3] != "20:00" {
		t.Errorf("slice[3] (Time) = %v, want '20:00'", slice[3])
	}
	if slice[4] != "Team1" {
		t.Errorf("slice[4] (HomeTeam) = %v, want 'Team1'", slice[4])
	}
	if slice[5] != "Team2" {
		t.Errorf("slice[5] (AwayTeam) = %v, want 'Team2'", slice[5])
	}
	if slice[6] != "Liga" {
		t.Errorf("slice[6] (League) = %v, want 'Liga'", slice[6])
	}
}

func TestGameToSlice_H2H15Stats(t *testing.T) {
	g := models.Game{
		Day:            1,
		Month:          1,
		Year:           2025,
		H2HHomeMatches: 10,
		H2HHomeWin1:    5,
		H2HHomeDraw:    3,
		H2HHomeWin2:    2,
	}

	slice := stats.GameToSlice(&g)

	// H2H15 stats start at column 7
	if slice[7] != 10 {
		t.Errorf("slice[7] (H2H15 Matches) = %v, want 10", slice[7])
	}
	if slice[8] != 5 {
		t.Errorf("slice[8] (H2H15 Win1) = %v, want 5", slice[8])
	}
	if slice[9] != 3 {
		t.Errorf("slice[9] (H2H15 Draw) = %v, want 3", slice[9])
	}
	if slice[10] != 2 {
		t.Errorf("slice[10] (H2H15 Win2) = %v, want 2", slice[10])
	}
}
