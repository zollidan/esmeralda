package processor

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/zollidan/esmeralda/internal/api"
	"github.com/zollidan/esmeralda/internal/models"
	"github.com/zollidan/esmeralda/internal/queue"
	"github.com/zollidan/esmeralda/internal/repository"
	"github.com/zollidan/esmeralda/internal/stats"
)

type result struct {
	index int
	game  models.Game
	err   error
}

func (p *Processor) ProcessMatches(ctx context.Context, client api.MatchFetcher, games *repository.GameRepository, matches []api.Match, totalMatches int, taskID string) error {
	limit := len(matches)
	results := make(chan result, limit)
	sem := make(chan struct{}, p.workers)

	var wg sync.WaitGroup

	for i, match := range matches[:limit] {
		wg.Add(1)

		go func() {
			defer wg.Done()

			select {
			case sem <- struct{}{}:
			case <-ctx.Done():
				results <- result{index: i, err: ctx.Err()}
				return
			}
			defer func() { <-sem }()

			game, err := buildGame(client, match)
			if err != nil {
				results <- result{index: i, err: fmt.Errorf("build game: %w", err)}
				return
			}

			game.TaskID = &taskID

			results <- result{index: i, game: game}
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	gameSlice := make([]models.Game, limit)
	processed := 0
	for r := range results {
		if r.err != nil {
			return fmt.Errorf("match error: %w", r.err)
		}
		gameSlice[r.index] = r.game
		processed++
		p.publishProgress(ctx, taskID, queue.StatusProcessing, totalMatches, processed)
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
		m, err := fetchMatches(client, api.MatchesFilter{
			TeamID: match.HomeTeam.ID,
			Status: api.MatchStatusFinished,
		}, fetchNeed{
			NeedHome:       25,
			NeedAway:       25,
			NeedH2HAny:     25,
			NeedH2HHome:    15,
			OpponentTeamID: match.AwayTeam.ID,
			Before:         date,
		})
		ch1 <- res{m, err}
	}()
	go func() {
		m, err := fetchMatches(client, api.MatchesFilter{
			TeamID: match.AwayTeam.ID,
			Status: api.MatchStatusFinished,
		}, fetchNeed{
			NeedHome: 25,
			NeedAway: 25,
			Before:   date,
		})
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
	d := s.Details

	return models.Game{
		Day:      date.Day(),
		Month:    int(date.Month()),
		Year:     date.Year(),
		Time:     time.UnixMilli(match.StartTimestamp).Format("15:04"),
		HomeTeam: match.HomeTeam.Name,
		AwayTeam: match.AwayTeam.Name,
		League:   formatLeagueName(match),

		H2HHomeMatches: s.H2H15.Matches,
		H2HHomeWin1:    s.H2H15.Win1,
		H2HHomeDraw:    s.H2H15.Draw,
		H2HHomeWin2:    s.H2H15.Win2,

		HomeTeamHomeMatches: s.Home25.Matches,
		HomeTeamHomeWin:     s.Home25.Win1,
		HomeTeamHomeDraw:    s.Home25.Draw,
		HomeTeamHomeLoss:    s.Home25.Win2,

		AwayTeamAwayMatches: s.AwayK2_25.Matches,
		AwayTeamAwayLoss:    s.AwayK2_25.Win2,
		AwayTeamAwayDraw:    s.AwayK2_25.Draw,
		AwayTeamAwayWin:     s.AwayK2_25.Win1,

		H2HHomeOver25Total: d.H2H.HomeField.Over25.Total,
		H2HHomeOver25Over:  d.H2H.HomeField.Over25.Over25,
		H2HHomeOver25Under: d.H2H.HomeField.Over25.Under25,

		HomeAllOver25Total: d.Home.HomeField.Over25.Total,
		HomeAllOver25Over:  d.Home.HomeField.Over25.Over25,
		HomeAllOver25Under: d.Home.HomeField.Over25.Under25,

		AwayAllOver25Total: d.Away.AwayField.Over25.Total,
		AwayAllOver25Over:  d.Away.AwayField.Over25.Over25,
		AwayAllOver25Under: d.Away.AwayField.Over25.Under25,

		H2HAnyGames25: d.H2H.AllField.Goals.Matches,
		H2HAnyGoals25: d.H2H.AllField.Goals.Goals,
		H2HAnyGames5:  d.H2H.AllField.Win5.Matches,
		H2HAnyGoals5:  d.H2H.AllField.Win5.Goals,
		H2HAnyGames3:  d.H2H.AllField.Win3.Matches,
		H2HAnyGoals3:  d.H2H.AllField.Win3.Goals,

		HomeAnyGames25: d.Home.AllField.Goals.Matches,
		HomeAnyGoals25: d.Home.AllField.Goals.Goals,
		HomeAnyGames5:  d.Home.AllField.Win5.Matches,
		HomeAnyGoals5:  d.Home.AllField.Win5.Goals,
		HomeAnyGames3:  d.Home.AllField.Win3.Matches,
		HomeAnyGoals3:  d.Home.AllField.Win3.Goals,

		AwayAnyGames25: d.Away.AllField.Goals.Matches,
		AwayAnyGoals25: d.Away.AllField.Goals.Goals,
		AwayAnyGames5:  d.Away.AllField.Win5.Matches,
		AwayAnyGoals5:  d.Away.AllField.Win5.Goals,
		AwayAnyGames3:  d.Away.AllField.Win3.Matches,
		AwayAnyGoals3:  d.Away.AllField.Win3.Goals,

		H2HHomeGames25: d.H2H.HomeField.Goals.Matches,
		H2HHomeGoals25: d.H2H.HomeField.Goals.Goals,
		H2HHomeGames5:  d.H2H.HomeField.Win5.Matches,
		H2HHomeGoals5:  d.H2H.HomeField.Win5.Goals,
		H2HHomeGames3:  d.H2H.HomeField.Win3.Matches,
		H2HHomeGoals3:  d.H2H.HomeField.Win3.Goals,

		HomeHomeGames25: d.Home.HomeField.Goals.Matches,
		HomeHomeGoals25: d.Home.HomeField.Goals.Goals,
		HomeHomeGames5:  d.Home.HomeField.Win5.Matches,
		HomeHomeGoals5:  d.Home.HomeField.Win5.Goals,
		HomeHomeGames3:  d.Home.HomeField.Win3.Matches,
		HomeHomeGoals3:  d.Home.HomeField.Win3.Goals,

		AwayAwayGames25: d.Away.AwayField.Goals.Matches,
		AwayAwayGoals25: d.Away.AwayField.Goals.Goals,
		AwayAwayGames5:  d.Away.AwayField.Win5.Matches,
		AwayAwayGoals5:  d.Away.AwayField.Win5.Goals,
		AwayAwayGames3:  d.Away.AwayField.Win3.Matches,
		AwayAwayGoals3:  d.Away.AwayField.Win3.Goals,
	}, nil
}

type fetchNeed struct {
	NeedHome       int
	NeedAway       int
	NeedH2HAny     int
	NeedH2HHome    int
	OpponentTeamID int
	Before         time.Time
}

func fetchMatches(client api.MatchFetcher, filter api.MatchesFilter, need fetchNeed) ([]api.Match, error) {
	var all []api.Match
	homeCount, awayCount := 0, 0
	h2hAnyCount, h2hHomeCount := 0, 0
	teamID := filter.TeamID

	for page := 1; page <= 20; page++ {
		filter.Page = page
		filter.PageSize = 25

		matches, total, err := client.GetMatches(filter)
		if err != nil {
			return nil, err
		}

		all = append(all, matches...)

		for _, m := range matches {
			matchDate, err := time.Parse("2006-01-02", m.DateEvent)
			if err != nil {
				continue
			}
			if !need.Before.IsZero() && !matchDate.Before(need.Before) {
				continue
			}

			if m.HomeTeam.ID == teamID {
				homeCount++
			}
			if m.AwayTeam.ID == teamID {
				awayCount++
			}

			if need.OpponentTeamID == 0 {
				continue
			}

			isH2HAny := (m.HomeTeam.ID == teamID && m.AwayTeam.ID == need.OpponentTeamID) ||
				(m.HomeTeam.ID == need.OpponentTeamID && m.AwayTeam.ID == teamID)
			if isH2HAny {
				h2hAnyCount++
			}

			if m.HomeTeam.ID == teamID && m.AwayTeam.ID == need.OpponentTeamID {
				h2hHomeCount++
			}
		}

		if homeCount >= need.NeedHome && awayCount >= need.NeedAway && h2hAnyCount >= need.NeedH2HAny && h2hHomeCount >= need.NeedH2HHome {
			break
		}

		if len(all) >= total {
			break
		}
	}

	return all, nil
}

func formatLeagueName(match api.Match) string {
	league := strings.TrimSpace(match.Tournament.Name)
	season := strings.TrimSpace(match.Season.Name)
	if season == "" || strings.Contains(league, season) {
		return league
	}
	return strings.TrimSpace(league + " " + season)
}
