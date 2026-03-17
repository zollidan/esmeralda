package stats

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/zollidan/esmeralda/internal/api"
)

func loadData(t *testing.T) []api.Match {
	t.Helper()

	// Minimal inline fixture: tests should not depend on external files.
	const team1ResponseMockJSON = `{
  "totalMatches": 2,
  "matches": [
    {
      "id": 1001,
      "status": "finished",
      "dateEvent": "2024-01-01",
      "homeTeam": { "id": 1079350, "name": "Team1" },
      "awayTeam": { "id": 2, "name": "Team2" },
      "homeScore": { "current": 1 },
      "awayScore": { "current": 0 }
    },
    {
      "id": 1002,
      "status": "finished",
      "dateEvent": "2024-01-02",
      "homeTeam": { "id": 3, "name": "Other" },
      "awayTeam": { "id": 1079350, "name": "Team1" },
      "homeScore": { "current": 0 },
      "awayScore": { "current": 2 }
    }
  ]
}`

	var resp api.MatchesResponse
	if err := json.Unmarshal([]byte(team1ResponseMockJSON), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	return resp.Matches
}

func TestFilterH2HHomeField(t *testing.T) {
	before := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	matches := []api.Match{
		{DateEvent: "2024-01-01", HomeTeam: api.Team{ID: 1}, AwayTeam: api.Team{ID: 2}},

		{DateEvent: "2024-01-02", HomeTeam: api.Team{ID: 2}, AwayTeam: api.Team{ID: 1}},
	}

	res := filterH2HHomeField(matches, 1, 2, before)

	if len(res) != 1 {
		t.Fatalf("expected 1 match, got %d", len(res))
	}
}

func TestFilterH2H_Unit(t *testing.T) {
	before := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	matches := []api.Match{
		{
			DateEvent: "2024-01-01",
			HomeTeam:  api.Team{ID: 1},
			AwayTeam:  api.Team{ID: 2},
		},
		{
			DateEvent: "2024-01-02",
			HomeTeam:  api.Team{ID: 2},
			AwayTeam:  api.Team{ID: 1},
		},
		{
			DateEvent: "2026-01-01", // позже cutoff
			HomeTeam:  api.Team{ID: 1},
			AwayTeam:  api.Team{ID: 2},
		},
	}

	res := filterH2H(matches, 1, 2, before)

	if len(res) != 2 {
		t.Fatalf("expected 2 H2H matches, got %d", len(res))
	}
}

func TestFilterHome_Unit(t *testing.T) {
	before := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	matches := []api.Match{
		{DateEvent: "2024-01-01", HomeTeam: api.Team{ID: 1}},
		{DateEvent: "2024-01-02", HomeTeam: api.Team{ID: 2}},
		{DateEvent: "2026-01-01", HomeTeam: api.Team{ID: 1}}, // позже cutoff
	}

	res := filterHome(matches, 1, before)

	if len(res) != 1 {
		t.Fatalf("expected 1 home match, got %d", len(res))
	}
}

func TestFilterAway_Unit(t *testing.T) {
	before := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	matches := []api.Match{
		{DateEvent: "2024-01-01", AwayTeam: api.Team{ID: 1}},
		{DateEvent: "2024-01-02", AwayTeam: api.Team{ID: 2}},
		{DateEvent: "2026-01-01", AwayTeam: api.Team{ID: 1}},
	}

	res := filterAway(matches, 1, before)

	if len(res) != 1 {
		t.Fatalf("expected 1 away match, got %d", len(res))
	}
}

func TestCalcTotal(t *testing.T) {
	teamID := 1

	matches := []api.Match{
		{
			HomeTeam:  api.Team{ID: 1},
			AwayTeam:  api.Team{ID: 2},
			HomeScore: api.Score{Current: 2},
			AwayScore: api.Score{Current: 1},
		},
		{
			HomeTeam:  api.Team{ID: 2},
			AwayTeam:  api.Team{ID: 1},
			HomeScore: api.Score{Current: 0},
			AwayScore: api.Score{Current: 3},
		},
		{
			HomeTeam:  api.Team{ID: 1},
			AwayTeam:  api.Team{ID: 2},
			HomeScore: api.Score{Current: 1},
			AwayScore: api.Score{Current: 1},
		},
	}

	stats := calcTotal(matches, teamID, 10)

	if stats.Matches != 3 {
		t.Fatalf("expected 3 matches, got %d", stats.Matches)
	}
	if stats.Win1 != 2 {
		t.Fatalf("expected 2 wins, got %d", stats.Win1)
	}
	if stats.Draw != 1 {
		t.Fatalf("expected 1 draw, got %d", stats.Draw)
	}
}

func TestCalcTotal_Limit(t *testing.T) {
	teamID := 1

	var matches []api.Match
	for i := 0; i < 20; i++ {
		matches = append(matches, api.Match{
			HomeTeam:  api.Team{ID: 1},
			AwayTeam:  api.Team{ID: 2},
			HomeScore: api.Score{Current: 1},
			AwayScore: api.Score{Current: 0},
		})
	}

	stats := calcTotal(matches, teamID, 5)

	if stats.Matches != 5 {
		t.Fatalf("expected 5 matches due to limit, got %d", stats.Matches)
	}
	if stats.Win1 != 5 {
		t.Fatalf("expected 5 wins, got %d", stats.Win1)
	}
}

func TestIntegration_WithRealJSON(t *testing.T) {
	matches := loadData(t)

	if len(matches) == 0 {
		t.Fatal("expected matches from JSON, got 0")
	}

	before := time.Now()

	res := filterHome(matches, 1079350, before)

	for _, m := range res {
		if m.HomeTeam.ID != 1079350 {
			t.Fatalf("unexpected home team id: %d", m.HomeTeam.ID)
		}
	}
}

func TestCalcTotal_Invariant(t *testing.T) {
	teamID := 1

	matches := []api.Match{
		{HomeTeam: api.Team{ID: 1}, AwayTeam: api.Team{ID: 2}, HomeScore: api.Score{Current: 1}, AwayScore: api.Score{Current: 0}},
		{HomeTeam: api.Team{ID: 1}, AwayTeam: api.Team{ID: 2}, HomeScore: api.Score{Current: 0}, AwayScore: api.Score{Current: 1}},
		{HomeTeam: api.Team{ID: 1}, AwayTeam: api.Team{ID: 2}, HomeScore: api.Score{Current: 1}, AwayScore: api.Score{Current: 1}},
	}

	stats := calcTotal(matches, teamID, 10)

	if stats.Matches != stats.Win1+stats.Draw+stats.Win2 {
		t.Fatal("broken invariant: matches != win1 + draw + win2")
	}
}

func TestSortByDateDesc(t *testing.T) {
	matches := []api.Match{
		{DateEvent: "2023-01-01"},
		{DateEvent: "2024-01-01"},
		{DateEvent: "2022-01-01"},
	}

	sorted := sortByDateDesc(matches)

	if sorted[0].DateEvent != "2024-01-01" {
		t.Fatal("sorting failed")
	}
}

func TestCalculate_Full(t *testing.T) {
	match := api.Match{
		DateEvent: "2025-01-10",
		HomeTeam:  api.Team{ID: 1},
		AwayTeam:  api.Team{ID: 2},
	}

	team1Matches := []api.Match{
		{
			DateEvent: "2024-01-01",
			HomeTeam:  api.Team{ID: 1},
			AwayTeam:  api.Team{ID: 2},
			HomeScore: api.Score{Current: 2},
			AwayScore: api.Score{Current: 1},
		},
	}

	team2Matches := []api.Match{
		{
			DateEvent: "2024-01-01",
			HomeTeam:  api.Team{ID: 3},
			AwayTeam:  api.Team{ID: 2},
			HomeScore: api.Score{Current: 0},
			AwayScore: api.Score{Current: 2},
		},
	}

	stats := Calculate(match, team1Matches, team2Matches)

	if stats.H2H15.Matches != 1 {
		t.Fatal("H2H broken")
	}

	if stats.AwayK2_25.Matches != 1 {
		t.Fatal("AwayK2_25 broken")
	}
}
