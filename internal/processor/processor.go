package processor

import (
	"fmt"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/zollidan/esmeralda/internal/api"
	"github.com/zollidan/esmeralda/internal/models"
	"github.com/zollidan/esmeralda/internal/stats"

	"gorm.io/gorm"
)

// сделать контекст для управления горутинами

const workers = 20

type result struct {
	index int
	game  models.Game
	err   error
}

func ProcessMatches(client *api.Client, db *gorm.DB, matches []api.Match, totalMatches int) error {

	bar := progressbar.NewOptions(
		len(matches),
		progressbar.OptionSetDescription("Processing matches"),
		progressbar.OptionSetWidth(40),
		progressbar.OptionShowCount(),
		progressbar.OptionShowIts(),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "=",
			SaucerHead:    ">",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)	

	results := make(chan result, len(matches))
	sem := make(chan struct{}, workers)

	var wg sync.WaitGroup

	for i, match := range matches[:20] {
		// fmt.Printf("Processing %d/%d (id=%d)\n", i+1, totalMatches, match.ID)

		// i, match := i, match
		wg.Add(1)

		go func() {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			game, err := buildGame(client, match, bar)
			results <- result{index: i, game: game, err: err}
		}()
	}

	fmt.Println("Все запросы отправлены, ожидаем завершения...")

	go func() {
		wg.Wait()
		close(results)
	}()

	games := make([]models.Game, len(matches))
	for r := range results {
		if r.err != nil {
			return fmt.Errorf("match error: %w", r.err)
		}
		games[r.index] = r.game
	}

	if err := db.CreateInBatches(games, 100).Error; err != nil {
		return fmt.Errorf("save games to db: %w", err)
	}

	return nil
}

func buildGame(client *api.Client, match api.Match, bar *progressbar.ProgressBar) (models.Game, error) {
	// start := time.Now()
	defer func() {
		// fmt.Printf("match %d processed in %v\n", match.ID, time.Since(start))
		bar.Add(1)
	}()

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
