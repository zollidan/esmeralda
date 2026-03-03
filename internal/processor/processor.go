package processor

import (
	"fmt"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/zollidan/esmeralda-ru-api-fetcher/internal/api"
	"github.com/zollidan/esmeralda-ru-api-fetcher/internal/export"
	"github.com/zollidan/esmeralda-ru-api-fetcher/internal/stats"
)

// сделать контекст для управления горутинами

const workers = 20

type result struct {
	index int
	row   stats.Row
	err   error
}

func ProcessMatches(client *api.Client, writer *export.Writer, matches []api.Match, totalMatches int) error {

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

	for i, match := range matches {
		// fmt.Printf("Processing %d/%d (id=%d)\n", i+1, totalMatches, match.ID)

		// i, match := i, match
		wg.Add(1)

		go func() {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			row, err := buildRow(client, match, bar)
			results <- result{index: i, row: row, err: err}
		}()
	}

	fmt.Println("Все запросы отправлены, ожидаем завершения...")

	go func() {
		wg.Wait()
		close(results)
	}()

	rows := make([]stats.Row, len(matches))
	for r := range results {
		if r.err != nil {
			return fmt.Errorf("match error: %w", r.err)
		}
		rows[r.index] = r.row
	}

	for _, row := range rows {
		if err := writer.WriteRow(row); err != nil {
			return fmt.Errorf("write row error: %w", err)
		}
	}

	return nil
}

func buildRow(client *api.Client, match api.Match, bar *progressbar.ProgressBar) (stats.Row, error) {
	// start := time.Now()
	defer func() {
		// fmt.Printf("match %d processed in %v\n", match.ID, time.Since(start))
		bar.Add(1)
	}()

	date, err := time.Parse("2006-01-02", match.DateEvent)
	if err != nil {
		return stats.Row{}, fmt.Errorf("parse date: %w", err)
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
		return stats.Row{}, fmt.Errorf("get home team matches: %w", r1.err)
	}
	if r2.err != nil {
		return stats.Row{}, fmt.Errorf("get away team matches: %w", r2.err)
	}

	statistic := stats.Calculate(match, r1.matches, r2.matches)
	return stats.Row{
		Day:       date.Day(),
		Month:     date.Month(),
		Year:      date.Year(),
		Time:      time.UnixMilli(match.StartTimestamp).Format("15:04"),
		HomeTeam:  match.HomeTeam.Name,
		AwayTeam:  match.AwayTeam.Name,
		League:    match.Tournament.Name,

		H2H15:     statistic.H2H15,
		H2H25:     statistic.H2H25,

		Home25:    statistic.Home25,
		Away15:    statistic.Away15,
		Away25:    statistic.Away25,
		AwayK2_25: statistic.AwayK2_25,  

		GoalsH2H:  statistic.GoalsH2H,
		GoalsHome: statistic.GoalsHome,
		GoalsAway: statistic.GoalsAway,
	}, nil
}