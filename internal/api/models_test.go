package api

import (
	"encoding/json"
	"os"
	"testing"
)

func TestDecodeMatchesResponse(t *testing.T) {
	data, err := os.ReadFile("../../real_matches_response.json")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}

	var resp MatchesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	t.Logf("OK: decoded %d matches", len(resp.Matches))
}
