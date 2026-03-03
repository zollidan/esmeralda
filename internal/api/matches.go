package api

import "fmt"

func (c *Client) GetMatches(filter MatchesFilter) ([]Match, int, error) {
	var resp MatchesResponse
	if err := c.get("/football/matches", filter.ToParams(), &resp); err != nil {
		return nil, 0, fmt.Errorf("get matches: %w", err)
	}
	return resp.Matches, resp.TotalMatches, nil
}
