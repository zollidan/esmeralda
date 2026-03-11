package server

import (
	"testing"
	"time"

	"github.com/zollidan/esmeralda/internal/models"
	"github.com/zollidan/esmeralda/internal/stats"
)

func TestGameToRow(t *testing.T) {
	g := models.Game{
		Day:      15,
		Month:    3,
		Year:     2025,
		Time:     "18:30",
		HomeTeam: "FC Home",
		AwayTeam: "FC Away",
		League:   "Premier League",
		H2H15:    models.TotalStats{Matches: 10, Win1: 5, Draw: 3, Win2: 2},
		H2H25:    models.TotalStats{Matches: 20, Win1: 10, Draw: 5, Win2: 5},
		Home25:   models.TotalStats{Matches: 25, Win1: 15, Draw: 5, Win2: 5},
		Away15:   models.TotalStats{Matches: 15, Win1: 7, Draw: 4, Win2: 4},
		Away25:   models.TotalStats{Matches: 25, Win1: 12, Draw: 6, Win2: 7},
		AwayK2_25: models.TotalStats{Matches: 25, Win1: 8, Draw: 8, Win2: 9},
		GoalsH2H: models.GoalStats{
			Over25Matches: 6, Over25Total: 24,
			Over3Matches: 3, Over3Total: 15,
			Over5Matches: 1, Over5Total: 6,
		},
		GoalsHome: models.GoalStats{
			Over25Matches: 10, Over25Total: 40,
			Over3Matches: 5, Over3Total: 25,
			Over5Matches: 2, Over5Total: 12,
		},
		GoalsAway: models.GoalStats{
			Over25Matches: 8, Over25Total: 32,
			Over3Matches: 4, Over3Total: 20,
			Over5Matches: 1, Over5Total: 7,
		},
	}

	row := gameToRow(g)

	if row.Day != 15 {
		t.Errorf("Day = %d, want 15", row.Day)
	}
	if row.Month != time.March {
		t.Errorf("Month = %v, want March", row.Month)
	}
	if row.Year != 2025 {
		t.Errorf("Year = %d, want 2025", row.Year)
	}
	if row.HomeTeam != "FC Home" {
		t.Errorf("HomeTeam = %q, want %q", row.HomeTeam, "FC Home")
	}
	if row.AwayTeam != "FC Away" {
		t.Errorf("AwayTeam = %q, want %q", row.AwayTeam, "FC Away")
	}
	if row.League != "Premier League" {
		t.Errorf("League = %q, want %q", row.League, "Premier League")
	}
	if row.H2H15.Matches != 10 {
		t.Errorf("H2H15.Matches = %d, want 10", row.H2H15.Matches)
	}
	if row.H2H25.Win1 != 10 {
		t.Errorf("H2H25.Win1 = %d, want 10", row.H2H25.Win1)
	}
	if row.GoalsH2H.Over25Matches != 6 {
		t.Errorf("GoalsH2H.Over25Matches = %d, want 6", row.GoalsH2H.Over25Matches)
	}
}

func TestGameToRow_ZeroValues(t *testing.T) {
	g := models.Game{
		Day:      1,
		Month:    1,
		Year:     2025,
		HomeTeam: "A",
		AwayTeam: "B",
	}

	row := gameToRow(g)

	if row.H2H15.Matches != 0 {
		t.Errorf("H2H15.Matches = %d, want 0", row.H2H15.Matches)
	}
	if row.GoalsH2H.Over25Matches != 0 {
		t.Errorf("GoalsH2H.Over25Matches = %d, want 0", row.GoalsH2H.Over25Matches)
	}
}

func TestGameToRow_TypeConversion(t *testing.T) {
	g := models.Game{
		Month: 12,
		H2H15: models.TotalStats{Matches: 5, Win1: 2, Draw: 1, Win2: 2},
	}

	row := gameToRow(g)

	// models.TotalStats -> stats.TotalStats should preserve values
	if row.H2H15 != (stats.TotalStats{Matches: 5, Win1: 2, Draw: 1, Win2: 2}) {
		t.Errorf("H2H15 conversion failed: %+v", row.H2H15)
	}

	// Month conversion: int -> time.Month
	if row.Month != time.December {
		t.Errorf("Month = %v, want December", row.Month)
	}
}

func TestGameToRow_ToSlice(t *testing.T) {
	g := models.Game{
		Day:      5,
		Month:    6,
		Year:     2025,
		Time:     "20:00",
		HomeTeam: "Team1",
		AwayTeam: "Team2",
		League:   "Liga",
	}

	row := gameToRow(g)
	slice := row.ToSlice()

	if len(slice) != len(stats.Columns) {
		t.Errorf("slice len = %d, want %d", len(slice), len(stats.Columns))
	}

	// First 7 values should be basic info
	if slice[0] != 5 {
		t.Errorf("slice[0] (Day) = %v, want 5", slice[0])
	}
	if slice[4] != "Team1" {
		t.Errorf("slice[4] (HomeTeam) = %v, want Team1", slice[4])
	}
	if slice[5] != "Team2" {
		t.Errorf("slice[5] (AwayTeam) = %v, want Team2", slice[5])
	}
}
