package api

import (
	"encoding/json"
	"testing"
)

func TestDecodeMatchesResponse(t *testing.T) {
	// Minimal inline fixture: tests should not depend on external files.
	const matchesResponseMockJSON = `{
  "totalMatches": 1,
  "matches": [
    {
      "id": 123,
      "status": "finished",
      "dateEvent": "2024-01-01",
      "homeTeam": { "id": 1, "name": "Home" },
      "awayTeam": { "id": 2, "name": "Away" },
      "homeScore": { "current": 2 },
      "awayScore": { "current": 1 }
    }
  ]
}`

	var resp MatchesResponse
	if err := json.Unmarshal([]byte(matchesResponseMockJSON), &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	t.Logf("OK: decoded %d matches", len(resp.Matches))
}
