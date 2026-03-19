package stats

// MatchStats represents all statistics for a single match
type MatchStats struct {
	// Основные статистики (базовые колонки)
	H2H15     TotalStats
	H2H25     TotalStats
	Home25    TotalStats
	Away15    TotalStats
	Away25    TotalStats
	AwayK2_25 TotalStats

	// Все остальное в JSON (расширенные данные)
	Details StatsDetails
}

// TotalStats represents wins, draws, losses for a team over N matches
type TotalStats struct {
	Matches int
	Win1    int
	Draw    int
	Win2    int
}

// StatsDetails holds all extended statistics as JSON
type StatsDetails struct {
	// H2H stats with different time windows
	H2H struct {
		HomeField struct {
			Over25 Over25Stats `json:"over25"`
			Goals  GoalsWindow `json:"goals"`
			Win5   WinStats    `json:"win5"` // last 5 matches
			Win3   WinStats    `json:"win3"` // last 3 matches
		} `json:"home_field"`
		AllField struct {
			Over25 Over25Stats `json:"over25"`
			Goals  GoalsWindow `json:"goals"`
			Win5   WinStats    `json:"win5"`
			Win3   WinStats    `json:"win3"`
		} `json:"all_field"`
	} `json:"h2h"`

	// Home team stats (Team 1 at home)
	Home struct {
		AllField struct {
			Over25 Over25Stats `json:"over25"`
			Win5   WinStats    `json:"win5"`
			Win3   WinStats    `json:"win3"`
			Goals  GoalsWindow `json:"goals"`
		} `json:"all_field"`
		HomeField struct {
			Over25 Over25Stats `json:"over25"`
			Win5   WinStats    `json:"win5"`
			Win3   WinStats    `json:"win3"`
			Goals  GoalsWindow `json:"goals"`
		} `json:"home_field"`
	} `json:"home"`

	// Away team stats (Team 2 away)
	Away struct {
		AllField struct {
			Over25 Over25Stats `json:"over25"`
			Win5   WinStats    `json:"win5"`
			Win3   WinStats    `json:"win3"`
			Goals  GoalsWindow `json:"goals"`
		} `json:"all_field"`
		AwayField struct {
			Over25 Over25Stats `json:"over25"`
			Win5   WinStats    `json:"win5"`
			Win3   WinStats    `json:"win3"`
			Goals  GoalsWindow `json:"goals"`
		} `json:"away_field"`
	} `json:"away"`
}

// Over25Stats - total matches, over 2.5 goals, under 2.5 goals
type Over25Stats struct {
	Total   int `json:"total"`
	Over25  int `json:"over25"`
	Under25 int `json:"under25"`
}

// WinStats - matches count and total goals (for 5 and 3 match windows)
type WinStats struct {
	Matches int `json:"matches"`
	Goals   int `json:"goals"`
}

// GoalsWindow - statistics for goals in a specific window (25, 5, or 3 matches)
type GoalsWindow struct {
	Matches int `json:"matches"`
	Goals   int `json:"goals"`
}
