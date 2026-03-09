package stats

import (
	"sort"
	"time"

	"github.com/zollidan/esmeralda/internal/api"
)

type MatchStats struct {
	H2H15     TotalStats
	H2H25     TotalStats
	Home25    TotalStats // команда_1 дома, последние 25
	Away15    TotalStats // команда_1 в гостях, последние 15
	Away25    TotalStats // команда_1 в гостях, последние 25
	AwayK2_25 TotalStats // команда_2 в гостях, последние 25
	GoalsH2H  GoalStats
	GoalsHome GoalStats
	GoalsAway GoalStats
}

func Calculate(match api.Match, team1Matches []api.Match, team2Matches []api.Match) MatchStats {
	date, _ := time.Parse("2006-01-02", match.DateEvent)

	homeTeamID := match.HomeTeam.ID
	awayTeamID := match.AwayTeam.ID

	h2h := sortByDateDesc(
		filterH2HHomeField(team1Matches, homeTeamID, awayTeamID, date),
	)

	home := sortByDateDesc(filterHome(team1Matches, homeTeamID, date))
	away1 := sortByDateDesc(filterAway(team1Matches, homeTeamID, date))
	away2 := sortByDateDesc(filterAway(team2Matches, awayTeamID, date))

	return MatchStats{
		H2H15: calcTotal(h2h, homeTeamID, 15),
		H2H25: calcTotal(h2h, homeTeamID, 25),

		Home25: calcTotal(home, homeTeamID, 25),

		Away15: calcTotal(away1, homeTeamID, 15),

		Away25: calcTotal(away1, homeTeamID, 25),

		AwayK2_25: calcTotal(away2, awayTeamID, 25),

		GoalsH2H:  calcGoals(h2h, 25),
		GoalsHome: calcGoals(home, 25),
		GoalsAway: calcGoals(away2, 25),
	}
}

func filterH2H(matches []api.Match, team1ID, team2ID int, before time.Time) []api.Match {
	var result []api.Match

	for _, m := range matches {
		date, err := time.Parse("2006-01-02", m.DateEvent)
		if err != nil {
			continue
		}

		if !date.Before(before) {
			continue
		}

		if (m.HomeTeam.ID == team1ID && m.AwayTeam.ID == team2ID) ||
			(m.HomeTeam.ID == team2ID && m.AwayTeam.ID == team1ID) {
			result = append(result, m)
		}
	}

	return result
}

func filterHome(matches []api.Match, teamID int, before time.Time) []api.Match {
	var result []api.Match
	for _, m := range matches {
		d, err := time.Parse("2006-01-02", m.DateEvent)
		if err != nil {
			continue
		}
		if m.HomeTeam.ID == teamID && d.Before(before) {
			result = append(result, m)
		}
	}
	return result
}

func filterAway(matches []api.Match, teamID int, before time.Time) []api.Match {
	var result []api.Match
	for _, m := range matches {
		d, err := time.Parse("2006-01-02", m.DateEvent)
		if err != nil {
			continue
		}
		if m.AwayTeam.ID == teamID && d.Before(before) {
			result = append(result, m)
		}
	}
	return result
}

func calcTotal(matches []api.Match, teamID int, limit int) TotalStats {
	var s TotalStats
	for _, m := range matches {
		if s.Matches >= limit {
			break
		}
		home := m.HomeScore.Current
		away := m.AwayScore.Current

		s.Matches++
		switch {
		case home > away:
			if m.HomeTeam.ID == teamID {
				s.Win1++
			} else {
				s.Win2++
			}
		case away > home:
			if m.AwayTeam.ID == teamID {
				s.Win1++
			} else {
				s.Win2++
			}
		default:
			s.Draw++
		}
	}
	return s
}

func calcGoals(matches []api.Match, limit int) GoalStats {
	var s GoalStats
	count := 0
	for _, m := range matches {
		if count >= limit {
			break
		}
		count++
		total := m.HomeScore.Current + m.AwayScore.Current

		if total > 2 {
			s.Over25Matches++
			s.Over25Total += total
		}
		if total > 3 {
			s.Over3Matches++
			s.Over3Total += total
		}
		if total > 5 {
			s.Over5Matches++
			s.Over5Total += total
		}
	}
	return s
}

func sortByDateDesc(matches []api.Match) []api.Match {
	sorted := make([]api.Match, len(matches))
	copy(sorted, matches)

	sort.Slice(sorted, func(i, j int) bool {
		di, _ := time.Parse("2006-01-02", sorted[i].DateEvent)
		dj, _ := time.Parse("2006-01-02", sorted[j].DateEvent)
		return di.After(dj) // новые сначала
	})
	return sorted
}

func filterH2HHomeField(matches []api.Match, homeID, awayID int, before time.Time) []api.Match {
	var result []api.Match

	for _, m := range matches {
		d, err := time.Parse("2006-01-02", m.DateEvent)
		if err != nil {
			continue
		}

		if !d.Before(before) {
			continue
		}

		if m.HomeTeam.ID == homeID && m.AwayTeam.ID == awayID {
			result = append(result, m)
		}
	}

	return result
}
