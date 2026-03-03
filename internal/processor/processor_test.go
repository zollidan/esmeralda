package processor

import (
	"testing"

	"github.com/zollidan/esmeralda-ru-api-fetcher/internal/api"
)

func TestBuldRow_Unit(t *testing.T) {
	_ = api.Match{
		Status: api.MatchStatusFinished,
		HomeTeam: api.Team{
			ID: 1079350,
		},
		AwayTeam: api.Team{
			ID: 1079347,
		},
		HomeScore: api.Score{
			Current: 2,
			Period1: 2,
			Period2: 0,
		},
		AwayScore: api.Score{
			Current: 3,
			Period1: 2,
			Period2: 1,
		},
	}

	
}