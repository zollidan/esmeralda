package processor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/zollidan/esmeralda/internal/api"
	"github.com/zollidan/esmeralda/internal/models"
	"github.com/zollidan/esmeralda/internal/queue"
	"github.com/zollidan/esmeralda/internal/repository"
	"github.com/zollidan/esmeralda/internal/stats"
)

// сделать контекст для управления горутинами

const workers = 20

type result struct {
	index int
	game  models.Game
	err   error
}

func (p *Processor) ProcessMatches(ctx context.Context, client api.MatchFetcher, games *repository.GameRepository, matches []api.Match, totalMatches int, taskID string) error {
	limit := 20
	results := make(chan result, limit)
	sem := make(chan struct{}, workers)

	var wg sync.WaitGroup

	for i, match := range matches[:limit] {
		wg.Add(1)

		go func() {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			game, err := buildGame(client, match)
			if err != nil {
				results <- result{index: i, err: fmt.Errorf("build game: %w", err)}
				return
			}

			game.TaskID = &taskID

			if err := p.publishProgress(ctx, taskID, queue.StatusPending, limit, i); err != nil {
				results <- result{index: i, err: fmt.Errorf("publish progress: %w", err)}
				return
			}

			results <- result{index: i, game: game}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	gameSlice := make([]models.Game, limit)
	for r := range results {
		if r.err != nil {
			return fmt.Errorf("match error: %w", r.err)
		}
		gameSlice[r.index] = r.game
	}

	if err := games.CreateInBatches(ctx, gameSlice, 100); err != nil {
		return fmt.Errorf("save games to db: %w", err)
	}

	return nil
}

func buildGame(client api.MatchFetcher, match api.Match) (models.Game, error) {
	date, err := time.Parse("2006-01-02", match.DateEvent)
	if err != nil {
		return models.Game{}, fmt.Errorf("parse date: %w", err)
	}

	type res struct {
		matches []api.Match
		err     error
	}

	ch1 := make(chan res, 1)
	ch2 := make(chan res, 1)

	go func() {
		m, _, err := client.GetMatches(api.MatchesFilter{TeamID: match.HomeTeam.ID, Status: api.MatchStatusFinished})
		ch1 <- res{m, err}
	}()
	go func() {
		m, _, err := client.GetMatches(api.MatchesFilter{TeamID: match.AwayTeam.ID, Status: api.MatchStatusFinished})
		ch2 <- res{m, err}
	}()

	r1, r2 := <-ch1, <-ch2
	if r1.err != nil {
		return models.Game{}, fmt.Errorf("get home team matches: %w", r1.err)
	}
	if r2.err != nil {
		return models.Game{}, fmt.Errorf("get away team matches: %w", r2.err)
	}

	s := stats.Calculate(match, r1.matches, r2.matches)
	return models.Game{
		Day:      date.Day(),
		Month:    int(date.Month()),
		Year:     date.Year(),
		Time:     time.UnixMilli(match.StartTimestamp).Format("15:04"),
		HomeTeam: match.HomeTeam.Name,
		AwayTeam: match.AwayTeam.Name,
		League:   match.Tournament.Name,

		H2H15:     models.TotalStats(s.H2H15),
		H2H25:     models.TotalStats(s.H2H25),
		Home25:    models.TotalStats(s.Home25),
		Away15:    models.TotalStats(s.Away15),
		Away25:    models.TotalStats(s.Away25),
		AwayK2_25: models.TotalStats(s.AwayK2_25),

		GoalsH2H:  models.GoalStats(s.GoalsH2H),
		GoalsHome: models.GoalStats(s.GoalsHome),
		GoalsAway: models.GoalStats(s.GoalsAway),
	}, nil
}
