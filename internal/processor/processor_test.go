package processor

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/zollidan/esmeralda/internal/api"
	"github.com/zollidan/esmeralda/internal/models"
	"github.com/zollidan/esmeralda/internal/queue"
	"github.com/zollidan/esmeralda/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// --- Mock MatchFetcher ---

type mockFetcher struct {
	matches      []api.Match
	totalMatches int
	err          error
	getCalls     int
}

func (m *mockFetcher) GetMatches(filter api.MatchesFilter) ([]api.Match, int, error) {
	m.getCalls++
	return m.matches, m.totalMatches, m.err
}

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	if err := db.AutoMigrate(&models.Task{}, &models.Game{}); err != nil {
		t.Fatalf("auto migrate: %v", err)
	}
	return db
}

func TestBuildGame_Success(t *testing.T) {
	match := api.Match{
		DateEvent:      "2025-01-15",
		StartTimestamp: time.Date(2025, 1, 15, 18, 30, 0, 0, time.UTC).UnixMilli(),
		HomeTeam:       api.Team{ID: 1, Name: "Team A"},
		AwayTeam:       api.Team{ID: 2, Name: "Team B"},
		Tournament:     api.Tournament{Name: "Premier League"},
	}

	fetcher := &mockFetcher{
		matches: []api.Match{
			{
				DateEvent: "2024-06-01",
				HomeTeam:  api.Team{ID: 1, Name: "Team A"},
				AwayTeam:  api.Team{ID: 2, Name: "Team B"},
				HomeScore: api.Score{Current: 2},
				AwayScore: api.Score{Current: 1},
			},
		},
		totalMatches: 1,
	}

	game, err := buildGame(fetcher, match)
	if err != nil {
		t.Fatalf("buildGame: %v", err)
	}

	if game.HomeTeam != "Team A" {
		t.Errorf("HomeTeam = %q, want %q", game.HomeTeam, "Team A")
	}
	if game.AwayTeam != "Team B" {
		t.Errorf("AwayTeam = %q, want %q", game.AwayTeam, "Team B")
	}
	if game.League != "Premier League" {
		t.Errorf("League = %q, want %q", game.League, "Premier League")
	}
	if game.Day != 15 {
		t.Errorf("Day = %d, want 15", game.Day)
	}
	if game.Month != 1 {
		t.Errorf("Month = %d, want 1", game.Month)
	}
	if game.Year != 2025 {
		t.Errorf("Year = %d, want 2025", game.Year)
	}
}

func TestBuildGame_InvalidDate(t *testing.T) {
	match := api.Match{
		DateEvent: "not-a-date",
		HomeTeam:  api.Team{ID: 1},
		AwayTeam:  api.Team{ID: 2},
	}

	fetcher := &mockFetcher{matches: []api.Match{}}

	_, err := buildGame(fetcher, match)
	if err == nil {
		t.Fatal("expected error for invalid date")
	}
}

func TestBuildGame_APIError(t *testing.T) {
	match := api.Match{
		DateEvent: "2025-01-15",
		HomeTeam:  api.Team{ID: 1},
		AwayTeam:  api.Team{ID: 2},
	}

	fetcher := &mockFetcher{err: errors.New("api unavailable")}

	_, err := buildGame(fetcher, match)
	if err == nil {
		t.Fatal("expected error when API fails")
	}
}

func TestBuildGame_StatsCalculation(t *testing.T) {
	match := api.Match{
		DateEvent: "2025-03-01",
		HomeTeam:  api.Team{ID: 10, Name: "Home"},
		AwayTeam:  api.Team{ID: 20, Name: "Away"},
	}

	fetcher := &mockFetcher{
		matches: []api.Match{
			{
				DateEvent: "2024-12-01",
				HomeTeam:  api.Team{ID: 10},
				AwayTeam:  api.Team{ID: 20},
				HomeScore: api.Score{Current: 3},
				AwayScore: api.Score{Current: 0},
			},
			{
				DateEvent: "2024-11-01",
				HomeTeam:  api.Team{ID: 10},
				AwayTeam:  api.Team{ID: 20},
				HomeScore: api.Score{Current: 1},
				AwayScore: api.Score{Current: 1},
			},
		},
		totalMatches: 2,
	}

	game, err := buildGame(fetcher, match)
	if err != nil {
		t.Fatalf("buildGame: %v", err)
	}

	if game.H2H15.Matches != 2 {
		t.Errorf("H2H15.Matches = %d, want 2", game.H2H15.Matches)
	}
}

func TestBuildGame_FetchesBothTeams(t *testing.T) {
	match := api.Match{
		DateEvent: "2025-01-15",
		HomeTeam:  api.Team{ID: 1},
		AwayTeam:  api.Team{ID: 2},
	}

	fetcher := &mockFetcher{matches: []api.Match{}, totalMatches: 0}

	_, err := buildGame(fetcher, match)
	if err != nil {
		t.Fatalf("buildGame: %v", err)
	}

	// GetMatches should be called 2 times: once for home team, once for away team
	if fetcher.getCalls != 2 {
		t.Errorf("getCalls = %d, want 2", fetcher.getCalls)
	}
}

func TestInit_CreatesProcessor(t *testing.T) {
	db := setupTestDB(t)
	gamesRepo := repository.NewGameRepository(db)
	fetcher := &mockFetcher{}

	p := Init(fetcher, gamesRepo, &queue.Producer{}, &queue.Producer{}, 5)
	if p == nil {
		t.Fatal("Init returned nil")
	}
	if p.apiClient == nil {
		t.Error("apiClient is nil")
	}
	if p.games == nil {
		t.Error("games is nil")
	}
}

func TestProcessMatches_SavesGamesToDB(t *testing.T) {
	db := setupTestDB(t)
	gamesRepo := repository.NewGameRepository(db)

	// Test the batch save path that ProcessMatches uses
	taskID := "task-1"
	gameSlice := make([]models.Game, 20)
	for i := range gameSlice {
		gameSlice[i] = models.Game{
			TaskID:   &taskID,
			Day:      15,
			Month:    1,
			Year:     2025,
			Time:     "18:00",
			HomeTeam: "Home",
			AwayTeam: "Away",
			League:   "League",
		}
	}

	ctx := context.Background()
	if err := gamesRepo.CreateInBatches(ctx, gameSlice, 100); err != nil {
		t.Fatalf("save games: %v", err)
	}

	saved, err := gamesRepo.FindByTaskID(ctx, "task-1")
	if err != nil {
		t.Fatalf("find by task: %v", err)
	}
	if len(saved) != 20 {
		t.Errorf("saved = %d, want 20", len(saved))
	}
}

func TestPublishResult_MessageFormat(t *testing.T) {
	result := queue.TaskResult{
		TaskID: "task-1",
		Status: queue.StatusDone,
		Error:  "",
	}

	data, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded queue.TaskResult
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.TaskID != "task-1" {
		t.Errorf("TaskID = %q, want %q", decoded.TaskID, "task-1")
	}
	if decoded.Status != queue.StatusDone {
		t.Errorf("Status = %q, want %q", decoded.Status, queue.StatusDone)
	}
}

func TestPublishProgress_MessageFormat(t *testing.T) {
	progress := queue.TaskProgress{
		TaskID:       "task-1",
		Status:       queue.StatusProcessing,
		TotalMatches: 20,
		CurrentMatch: 5,
	}

	data, err := json.Marshal(progress)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var decoded queue.TaskProgress
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if decoded.TotalMatches != 20 {
		t.Errorf("TotalMatches = %d, want 20", decoded.TotalMatches)
	}
	if decoded.CurrentMatch != 5 {
		t.Errorf("CurrentMatch = %d, want 5", decoded.CurrentMatch)
	}
}

func TestBuildGame_TimeFormatting(t *testing.T) {
	// StartTimestamp 2025-01-15 18:30 UTC
	ts := time.Date(2025, 1, 15, 18, 30, 0, 0, time.UTC).UnixMilli()
	match := api.Match{
		DateEvent:      "2025-01-15",
		StartTimestamp: ts,
		HomeTeam:       api.Team{ID: 1, Name: "A"},
		AwayTeam:       api.Team{ID: 2, Name: "B"},
	}

	fetcher := &mockFetcher{matches: []api.Match{}}

	game, err := buildGame(fetcher, match)
	if err != nil {
		t.Fatalf("buildGame: %v", err)
	}

	// Time should be formatted as HH:MM
	if len(game.Time) != 5 || game.Time[2] != ':' {
		t.Errorf("Time format unexpected: %q", game.Time)
	}
}

func TestBuildGame_EmptyHistory(t *testing.T) {
	match := api.Match{
		DateEvent: "2025-01-15",
		HomeTeam:  api.Team{ID: 1, Name: "A"},
		AwayTeam:  api.Team{ID: 2, Name: "B"},
	}

	fetcher := &mockFetcher{matches: []api.Match{}, totalMatches: 0}

	game, err := buildGame(fetcher, match)
	if err != nil {
		t.Fatalf("buildGame: %v", err)
	}

	// All stats should be zero with no history
	if game.H2H15.Matches != 0 {
		t.Errorf("H2H15.Matches = %d, want 0", game.H2H15.Matches)
	}
	if game.H2H25.Matches != 0 {
		t.Errorf("H2H25.Matches = %d, want 0", game.H2H25.Matches)
	}
	if game.Home25.Matches != 0 {
		t.Errorf("Home25.Matches = %d, want 0", game.Home25.Matches)
	}
}
