package stats

import (
	"testing"
	"time"
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

	// Verify all headers are non-empty strings
	for i, h := range headers {
		s, ok := h.(string)
		if !ok {
			t.Errorf("headers[%d] is not string: %T", i, h)
			continue
		}
		if s == "" {
			t.Errorf("headers[%d] is empty", i)
		}
	}
}

func TestRow_ToSlice(t *testing.T) {
	row := Row{
		Day:      15,
		Month:    time.March,
		Year:     2025,
		Time:     "18:00",
		HomeTeam: "Team A",
		AwayTeam: "Team B",
		League:   "League 1",
		H2H15:    TotalStats{Matches: 10, Win1: 5, Draw: 3, Win2: 2},
	}

	slice := row.ToSlice()

	if len(slice) != len(Columns) {
		t.Fatalf("slice len = %d, want %d", len(slice), len(Columns))
	}

	// Verify basic fields
	if slice[0] != 15 { // Day
		t.Errorf("Day = %v, want 15", slice[0])
	}
	if slice[1] != 3 { // Month as int
		t.Errorf("Month = %v, want 3", slice[1])
	}
	if slice[2] != 2025 { // Year
		t.Errorf("Year = %v, want 2025", slice[2])
	}
	if slice[3] != "18:00" { // Time
		t.Errorf("Time = %v, want 18:00", slice[3])
	}
	if slice[4] != "Team A" { // HomeTeam
		t.Errorf("HomeTeam = %v, want Team A", slice[4])
	}
	if slice[5] != "Team B" { // AwayTeam
		t.Errorf("AwayTeam = %v, want Team B", slice[5])
	}
	if slice[6] != "League 1" { // League
		t.Errorf("League = %v, want League 1", slice[6])
	}

	// H2H15 stats (columns 7-10)
	if slice[7] != 10 { // H2H15 Matches
		t.Errorf("H2H15.Matches = %v, want 10", slice[7])
	}
	if slice[8] != 5 { // H2H15 Win1
		t.Errorf("H2H15.Win1 = %v, want 5", slice[8])
	}
	if slice[9] != 3 { // H2H15 Draw
		t.Errorf("H2H15.Draw = %v, want 3", slice[9])
	}
	if slice[10] != 2 { // H2H15 Win2
		t.Errorf("H2H15.Win2 = %v, want 2", slice[10])
	}
}

func TestRow_ToSlice_ZeroValues(t *testing.T) {
	row := Row{}
	slice := row.ToSlice()

	if len(slice) != len(Columns) {
		t.Fatalf("slice len = %d, want %d", len(slice), len(Columns))
	}

	// All stat columns should be 0
	for i := 7; i < len(slice); i++ {
		if v, ok := slice[i].(int); ok && v != 0 {
			t.Errorf("slice[%d] = %d, want 0", i, v)
		}
	}
}

func TestHeaders_ColumnCount(t *testing.T) {
	headers := Headers()
	// Based on column.go there should be at least 7 base + 24 total stats + 18 goals = 49+ columns
	if len(headers) < 40 {
		t.Errorf("headers count = %d, expected at least 40", len(headers))
	}
}

func TestRow_ToSlice_Length_MatchesHeaders(t *testing.T) {
	row := Row{}
	slice := row.ToSlice()
	headers := Headers()

	if len(slice) != len(headers) {
		t.Errorf("ToSlice len = %d, Headers len = %d, should match", len(slice), len(headers))
	}
}

func TestRow_ToSlice_AllGoalStats(t *testing.T) {
	row := Row{
		GoalsH2H: GoalStats{
			Over25Matches: 10, Over25Total: 40,
			Over3Matches: 5, Over3Total: 25,
			Over5Matches: 2, Over5Total: 12,
		},
		GoalsHome: GoalStats{
			Over25Matches: 8, Over25Total: 32,
			Over3Matches: 4, Over3Total: 20,
			Over5Matches: 1, Over5Total: 7,
		},
		GoalsAway: GoalStats{
			Over25Matches: 6, Over25Total: 24,
			Over3Matches: 3, Over3Total: 15,
			Over5Matches: 0, Over5Total: 0,
		},
	}

	slice := row.ToSlice()

	// Verify the slice has non-zero values (goal stats are in the latter part)
	hasNonZero := false
	for _, v := range slice {
		if num, ok := v.(int); ok && num > 0 {
			hasNonZero = true
			break
		}
	}
	if !hasNonZero {
		t.Error("expected at least some non-zero values in slice")
	}
}
