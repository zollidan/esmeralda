package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c := NewClient("https://example.com", "token-123")
	if c.baseURL != "https://example.com" {
		t.Errorf("baseURL = %q, want %q", c.baseURL, "https://example.com")
	}
	if c.token != "token-123" {
		t.Errorf("token = %q, want %q", c.token, "token-123")
	}
	if c.httpClient == nil {
		t.Fatal("httpClient is nil")
	}
}

func TestClient_Get_Success(t *testing.T) {
	expected := MatchesResponse{
		TotalMatches: 2,
		Matches: []Match{
			{ID: 1, Status: MatchStatusFinished},
			{ID: 2, Status: MatchStatusFinished},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "test-token" {
			t.Errorf("Authorization header = %q, want %q", r.Header.Get("Authorization"), "test-token")
		}
		if r.URL.Path != "/test-path" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/test-path")
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(expected)
		if err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	c := &Client{
		httpClient: server.Client(),
		baseURL:    server.URL,
		token:      "test-token",
	}

	var result MatchesResponse
	err := c.get("/test-path", nil, &result)
	if err != nil {
		t.Fatalf("get: %v", err)
	}

	if result.TotalMatches != 2 {
		t.Errorf("TotalMatches = %d, want 2", result.TotalMatches)
	}
	if len(result.Matches) != 2 {
		t.Errorf("len(Matches) = %d, want 2", len(result.Matches))
	}
}

func TestClient_Get_WithParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("date") != "2025-01-15" {
			t.Errorf("date param = %q, want %q", r.URL.Query().Get("date"), "2025-01-15")
		}
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(MatchesResponse{})
		if err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	c := &Client{
		httpClient: server.Client(),
		baseURL:    server.URL,
		token:      "token",
	}

	params := url.Values{}
	params.Set("date", "2025-01-15")

	var result MatchesResponse
	if err := c.get("/football/matches", params, &result); err != nil {
		t.Fatalf("get: %v", err)
	}
}

func TestClient_Get_Non200Status(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c := &Client{
		httpClient: server.Client(),
		baseURL:    server.URL,
		token:      "token",
	}

	var result MatchesResponse
	err := c.get("/test", nil, &result)
	if err == nil {
		t.Fatal("expected error for non-200 response")
	}
}

func TestClient_Get_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte("not json")); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	c := &Client{
		httpClient: server.Client(),
		baseURL:    server.URL,
		token:      "token",
	}

	var result MatchesResponse
	err := c.get("/test", nil, &result)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestClient_Get_ConnectionError(t *testing.T) {
	c := &Client{
		httpClient: &http.Client{Timeout: 1 * time.Millisecond},
		baseURL:    "http://127.0.0.1:1", // unreachable
		token:      "token",
	}

	var result MatchesResponse
	err := c.get("/test", nil, &result)
	if err == nil {
		t.Fatal("expected error for connection failure")
	}
}

func TestClient_GetMatches(t *testing.T) {
	expected := MatchesResponse{
		TotalMatches: 3,
		Matches: []Match{
			{ID: 10, Status: MatchStatusFinished, DateEvent: "2025-01-01"},
			{ID: 11, Status: MatchStatusFinished, DateEvent: "2025-01-01"},
			{ID: 12, Status: MatchStatusFinished, DateEvent: "2025-01-01"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/football/matches" {
			t.Errorf("path = %q, want %q", r.URL.Path, "/football/matches")
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(expected); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	c := &Client{
		httpClient: server.Client(),
		baseURL:    server.URL,
		token:      "test-token",
	}

	matches, total, err := c.GetMatches(MatchesFilter{
		Date: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatalf("GetMatches: %v", err)
	}
	if total != 3 {
		t.Errorf("total = %d, want 3", total)
	}
	if len(matches) != 3 {
		t.Errorf("len(matches) = %d, want 3", len(matches))
	}
}

func TestClient_GetMatches_APIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer server.Close()

	c := &Client{
		httpClient: server.Client(),
		baseURL:    server.URL,
		token:      "test-token",
	}

	_, _, err := c.GetMatches(MatchesFilter{})
	if err == nil {
		t.Fatal("expected error for API failure")
	}
}

func TestMatchesFilter_ToParams_Empty(t *testing.T) {
	f := MatchesFilter{}
	params := f.ToParams()
	if len(params) != 0 {
		t.Errorf("expected empty params, got %v", params)
	}
}

func TestMatchesFilter_ToParams_Date(t *testing.T) {
	f := MatchesFilter{
		Date: time.Date(2025, 3, 15, 0, 0, 0, 0, time.UTC),
	}
	params := f.ToParams()
	if params.Get("date") != "2025-03-15" {
		t.Errorf("date = %q, want %q", params.Get("date"), "2025-03-15")
	}
}

func TestMatchesFilter_ToParams_TeamID(t *testing.T) {
	f := MatchesFilter{
		TeamID: 42,
		Status: MatchStatusFinished,
	}
	params := f.ToParams()
	if params.Get("team_id") != "42" {
		t.Errorf("team_id = %q, want %q", params.Get("team_id"), "42")
	}
	if params.Get("status") != "finished" {
		t.Errorf("status = %q, want %q", params.Get("status"), "finished")
	}
}

func TestMatchesFilter_ToParams_TournamentIDs(t *testing.T) {
	f := MatchesFilter{
		TournamentIDs: []int{1, 2, 3},
	}
	params := f.ToParams()
	if params.Get("tournament_id") != "1,2,3" {
		t.Errorf("tournament_id = %q, want %q", params.Get("tournament_id"), "1,2,3")
	}
}

func TestMatchesFilter_ToParams_SeasonID(t *testing.T) {
	f := MatchesFilter{
		SeasonID: 100,
	}
	params := f.ToParams()
	if params.Get("season_id") != "100" {
		t.Errorf("season_id = %q, want %q", params.Get("season_id"), "100")
	}
}

func TestMatchesFilter_ToParams_AllFields(t *testing.T) {
	f := MatchesFilter{
		Date:          time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
		TournamentIDs: []int{10, 20},
		SeasonID:      50,
		TeamID:        7,
		Status:        MatchStatusInProgress,
	}
	params := f.ToParams()

	if params.Get("date") != "2025-06-01" {
		t.Errorf("date = %q", params.Get("date"))
	}
	if params.Get("tournament_id") != "10,20" {
		t.Errorf("tournament_id = %q", params.Get("tournament_id"))
	}
	if params.Get("season_id") != "50" {
		t.Errorf("season_id = %q", params.Get("season_id"))
	}
	if params.Get("team_id") != "7" {
		t.Errorf("team_id = %q", params.Get("team_id"))
	}
	if params.Get("status") != "inprogress" {
		t.Errorf("status = %q", params.Get("status"))
	}
}
