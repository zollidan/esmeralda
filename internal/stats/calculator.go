package stats

import (
	"fmt"
	"sort"
	"time"

	"github.com/zollidan/esmeralda/internal/api"
)

func Calculate(match api.Match, team1Matches []api.Match, team2Matches []api.Match) MatchStats {
	date, _ := time.Parse("2006-01-02", match.DateEvent)
	homeTeamID := match.HomeTeam.ID
	awayTeamID := match.AwayTeam.ID

	// Фильтруем матчи по времени (ДО текущей даты)
	team1All := filterTeam(team1Matches, homeTeamID, date)
	team1Home := filterHome(team1Matches, homeTeamID, date)
	team1Away := filterAway(team1Matches, homeTeamID, date)
	team2All := filterTeam(team2Matches, awayTeamID, date)
	team2Away := filterAway(team2Matches, awayTeamID, date)

	// H2H фильтры
	combined := combineUniqueMatches(team1Matches, team2Matches)
	h2hAny := filterH2H(combined, homeTeamID, awayTeamID, date)
	h2hHome := filterH2HHomeField(combined, homeTeamID, awayTeamID, date)

	// Сортируем по дате (новые первыми)
	team1All = sortByDateDesc(team1All)
	team1Home = sortByDateDesc(team1Home)
	team1Away = sortByDateDesc(team1Away)
	team2All = sortByDateDesc(team2All)
	team2Away = sortByDateDesc(team2Away)
	h2hAny = sortByDateDesc(h2hAny)
	h2hHome = sortByDateDesc(h2hHome)

	return MatchStats{
		// Базовые статистики (для отдельных колонок)
		H2H15:     calcTotal(h2hHome, homeTeamID, 15),
		H2H25:     calcTotal(h2hAny, homeTeamID, 25),
		Home25:    calcTotal(team1Home, homeTeamID, 25),
		Away15:    calcTotal(team1Away, homeTeamID, 15),
		Away25:    calcTotal(team1Away, homeTeamID, 25),
		AwayK2_25: calcTotal(team2Away, awayTeamID, 25),

		// Расширенные данные (JSON)
		Details: buildDetails(h2hAny, h2hHome, team1All, team1Home, team1Away, team2All, team2Away),
	}
}

// buildDetails создает расширенные статистики для JSON
func buildDetails(h2hAny, h2hHome, team1All, team1Home, team1Away, team2All, team2Away []api.Match) StatsDetails {
	var d StatsDetails

	// H2H stats - любое поле
	d.H2H.AllField.Over25 = calcOver25(h2hAny, 25)
	d.H2H.AllField.Goals = calcGoalsWindow(h2hAny, 25)
	d.H2H.AllField.Win5 = calcWinWindow(h2hAny, 5)
	d.H2H.AllField.Win3 = calcWinWindow(h2hAny, 3)

	// H2H stats - на поле хозяев
	d.H2H.HomeField.Over25 = calcOver25(h2hHome, 25)
	d.H2H.HomeField.Goals = calcGoalsWindow(h2hHome, 25)
	d.H2H.HomeField.Win5 = calcWinWindow(h2hHome, 5)
	d.H2H.HomeField.Win3 = calcWinWindow(h2hHome, 3)

	// Home team stats (Team 1 at home)
	d.Home.AllField.Over25 = calcOver25(team1All, 25)
	d.Home.AllField.Win5 = calcWinWindow(team1All, 5)
	d.Home.AllField.Win3 = calcWinWindow(team1All, 3)
	d.Home.AllField.Goals = calcGoalsWindow(team1All, 25)

	d.Home.HomeField.Over25 = calcOver25(team1Home, 25) // все матчи дома уже на своем поле
	d.Home.HomeField.Win5 = calcWinWindow(team1Home, 5)
	d.Home.HomeField.Win3 = calcWinWindow(team1Home, 3)
	d.Home.HomeField.Goals = calcGoalsWindow(team1Home, 25)

	// Away team stats (Team 2 away)
	d.Away.AllField.Over25 = calcOver25(team2All, 25)
	d.Away.AllField.Win5 = calcWinWindow(team2All, 5)
	d.Away.AllField.Win3 = calcWinWindow(team2All, 3)
	d.Away.AllField.Goals = calcGoalsWindow(team2All, 25)

	d.Away.AwayField.Over25 = calcOver25(team2Away, 25) // все матчи в гостях уже на выезде
	d.Away.AwayField.Win5 = calcWinWindow(team2Away, 5)
	d.Away.AwayField.Win3 = calcWinWindow(team2Away, 3)
	d.Away.AwayField.Goals = calcGoalsWindow(team2Away, 25)

	return d
}

// combineUniqueMatches объединяет срезы матчей без дублей по match ID.
func combineUniqueMatches(slices ...[]api.Match) []api.Match {
	seen := make(map[int]struct{})
	result := make([]api.Match, 0)

	for _, slice := range slices {
		for _, m := range slice {
			if _, ok := seen[m.ID]; ok {
				continue
			}
			seen[m.ID] = struct{}{}
			result = append(result, m)
		}
	}

	return result
}

// calcTotal подсчитывает побед/ничьи/поражений за N матчей
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

// calcOver25 подсчитывает Over/Under 2.5 голов за N матчей
func calcOver25(matches []api.Match, limit int) Over25Stats {
	var s Over25Stats
	count := 0
	for _, m := range matches {
		if count >= limit {
			break
		}
		count++
		total := m.HomeScore.Current + m.AwayScore.Current
		s.Total++

		if total > 2 {
			s.Over25++
		} else {
			s.Under25++
		}
	}
	return s
}

// calcWinWindow подсчитывает матчи и голы за N матчей (для окон 5, 3)
func calcWinWindow(matches []api.Match, limit int) WinStats {
	var s WinStats
	for _, m := range matches {
		if s.Matches >= limit {
			break
		}
		s.Matches++
		s.Goals += m.HomeScore.Current + m.AwayScore.Current
	}
	return s
}

// calcGoalsWindow подсчитывает матчи и голы за N матчей
func calcGoalsWindow(matches []api.Match, limit int) GoalsWindow {
	var g GoalsWindow
	for _, m := range matches {
		if g.Matches >= limit {
			break
		}
		g.Matches++
		g.Goals += m.HomeScore.Current + m.AwayScore.Current
	}
	return g
}

// filterH2H возвращает все очные встречи (любое поле)
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

// filterH2HHomeField возвращает очные встречи только когда team1 играет дома
func filterH2HHomeField(matches []api.Match, team1ID, team2ID int, before time.Time) []api.Match {
	var result []api.Match
	for _, m := range matches {
		date, err := time.Parse("2006-01-02", m.DateEvent)
		if err != nil {
			continue
		}
		if !date.Before(before) {
			continue
		}
		if m.HomeTeam.ID == team1ID && m.AwayTeam.ID == team2ID {
			result = append(result, m)
		}
	}
	return result
}

// filterHome возвращает матчи team на доме
func filterHome(matches []api.Match, teamID int, before time.Time) []api.Match {
	var result []api.Match
	for _, m := range matches {
		date, err := time.Parse("2006-01-02", m.DateEvent)
		if err != nil {
			continue
		}
		if !date.Before(before) {
			continue
		}
		if m.HomeTeam.ID == teamID {
			result = append(result, m)
		}
	}
	return result
}

// filterTeam возвращает все матчи команды на любом поле
func filterTeam(matches []api.Match, teamID int, before time.Time) []api.Match {
	var result []api.Match
	for _, m := range matches {
		date, err := time.Parse("2006-01-02", m.DateEvent)
		if err != nil {
			continue
		}
		if !date.Before(before) {
			continue
		}
		if m.HomeTeam.ID == teamID || m.AwayTeam.ID == teamID {
			result = append(result, m)
		}
	}
	return result
}

// filterAway возвращает матчи team в гостях
func filterAway(matches []api.Match, teamID int, before time.Time) []api.Match {
	var result []api.Match
	for _, m := range matches {
		date, err := time.Parse("2006-01-02", m.DateEvent)
		if err != nil {
			continue
		}
		if !date.Before(before) {
			continue
		}
		if m.AwayTeam.ID == teamID {
			result = append(result, m)
		}
	}
	return result
}

// sortByDateDesc сортирует матчи по дате (новые первыми)
func sortByDateDesc(matches []api.Match) []api.Match {
	sorted := make([]api.Match, len(matches))
	copy(sorted, matches)
	sort.Slice(sorted, func(i, j int) bool {
		di, errI := time.Parse("2006-01-02", sorted[i].DateEvent)
		dj, errJ := time.Parse("2006-01-02", sorted[j].DateEvent)
		if errI != nil || errJ != nil {
			return fmt.Sprintf("%s-%d", sorted[i].DateEvent, sorted[i].ID) > fmt.Sprintf("%s-%d", sorted[j].DateEvent, sorted[j].ID)
		}
		return di.After(dj)
	})
	return sorted
}
