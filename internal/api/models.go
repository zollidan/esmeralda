package api

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type MatchStatus string

const (
	MatchStatusNotStarted   MatchStatus = "notstarted"
	MatchStatusInProgress   MatchStatus = "inprogress"
	MatchStatusFinished     MatchStatus = "finished"
	MatchStatusCanceled     MatchStatus = "canceled"
	MatchStatusPostponed    MatchStatus = "postponed"
	MatchStatusInterrupted  MatchStatus = "interrupted"
	MatchStatusSuspended    MatchStatus = "suspended"
	MatchStatusDelayed      MatchStatus = "delayed"
	MatchStatusWillContinue MatchStatus = "willcontinue"
)

type MatchesFilter struct {
	Date          time.Time
	TournamentIDs []int
	SeasonID      int
	TeamID        int
	Status        MatchStatus
}

func (f MatchesFilter) ToParams() url.Values {
	params := url.Values{}

	if !f.Date.IsZero() {
		params.Set("date", f.Date.Format("2006-01-02"))
	}
	if len(f.TournamentIDs) > 0 {
		ids := make([]string, len(f.TournamentIDs))
		for i, id := range f.TournamentIDs {
			ids[i] = strconv.Itoa(id)
		}
		params.Set("tournament_id", strings.Join(ids, ","))
	}
	if f.SeasonID != 0 {
		params.Set("season_id", fmt.Sprintf("%d", f.SeasonID))
	}
	if f.TeamID != 0 {
		params.Set("team_id", fmt.Sprintf("%d", f.TeamID))
	}
	if f.Status != "" {
		params.Set("status", string(f.Status))
	}

	return params
}

type MatchesResponse struct {
	TotalMatches int     `json:"totalMatches"`
	Matches      []Match `json:"matches"`
}

type Match struct {
	ID        int         `json:"id"`
	Status    MatchStatus `json:"status"`
	DateEvent string      `json:"dateEvent"`
	// yyyy-mm-dd
	StartTimestamp int64 `json:"startTimestamp"`
	// in millis
	CurrentMatchMinute int           `json:"currentMatchMinute"`
	CurrentMatchSecond int           `json:"currentMatchSecond"`
	Tournament         Tournament    `json:"tournament"`
	Category           Category      `json:"category"`
	RoundInfo          RoundInfo     `json:"roundInfo"`
	Season             Season        `json:"season"`
	Venue              *Venue        `json:"venue"`
	Referee            *Referee      `json:"referee"`
	HomeTeam           Team          `json:"homeTeam"`
	AwayTeam           Team          `json:"awayTeam"`
	HomeScore          Score         `json:"homeScore"`
	AwayScore          Score         `json:"awayScore"`
	LiveEvents         []LiveEvent   `json:"liveEvents"`
	MatchStatistics    []StatsPeriod `json:"matchStatistics"`
	OddsBase           []OddsMarket  `json:"oddsBase"`
	Highlights         []Highlight   `json:"highlights"`
}

type Translations struct {
	RU string `json:"ru"`
}

type Tournament struct {
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	Translations Translations `json:"translations"`
	Image        string       `json:"image"`
}

type Category struct {
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	Translations Translations `json:"translations"`
	Image        string       `json:"image"`
}

type RoundInfo struct {
	Name  string `json:"name"`
	Round int    `json:"round"`
}

type Season struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Year string `json:"year"`
}

type Venue struct {
	Name     string      `json:"name"`
	Capacity int         `json:"capacity"`
	City     NamedEntity `json:"city"`
	Country  NamedEntity `json:"country"`
}

type NamedEntity struct {
	Name string `json:"name"`
}

type Referee struct {
	Name           string       `json:"name"`
	YellowCards    int          `json:"yellowCards"`
	RedCards       int          `json:"redCards"`
	YellowRedCards int          `json:"yellowRedCards"`
	Games          int          `json:"games"`
	Country        NamedEntity  `json:"country"`
	Translations   Translations `json:"translations"`
}

type Team struct {
	ID           int          `json:"id"`
	Name         string       `json:"name"`
	Translations Translations `json:"translations"`
	Gender       string       `json:"gender"`
	Country      string       `json:"country"`
	Manager      *Manager     `json:"manager"`
	Image        string       `json:"image"`
	Lineup       *Lineup      `json:"lineup"`
}

type Manager struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Lineup struct {
	Players   []Player `json:"players"`
	Formation string   `json:"formation"`
}

type Player struct {
	ID          int                    `json:"id"`
	Position    string                 `json:"position"`
	Statistics  map[string]interface{} `json:"statistics"`
	ShirtNumber int                    `json:"shirtNumber"`
	Captain     bool                   `json:"captain"`
	Substitute  bool                   `json:"substitute"`
}

type Score struct {
	Current int `json:"current"`
	Period1 int `json:"period1"`
	Period2 int `json:"period2"`
}

type LiveEvent struct {
	Time        int       `json:"time"`
	TimeSeconds int       `json:"timeSeconds"`
	Type        string    `json:"type"`
	Class       string    `json:"class"`
	Team        string    `json:"team"`
	Player      PlayerRef `json:"player"`
	PlayerIn    PlayerRef `json:"playerIn"`
	PlayerOut   PlayerRef `json:"playerOut"`
	Reason      string    `json:"reason"`
	From        string    `json:"from"`
	HomeScore   int       `json:"homeScore"`
	AwayScore   int       `json:"awayScore"`
}

type PlayerRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type StatsPeriod struct {
	Period string       `json:"period"`
	Groups []StatsGroup `json:"groups"`
}

type StatsGroup struct {
	GroupName       string      `json:"groupName"`
	StatisticsItems []StatsItem `json:"statisticsItems"`
}

type StatsItem struct {
	Name      string          `json:"name"`
	HomeValue json.RawMessage `json:"homeValue"`
	AwayValue json.RawMessage `json:"awayValue"`
	Key       string          `json:"key"`
}

type OddsMarket struct {
	Name      string       `json:"name"`
	Group     string       `json:"group"`
	Period    string       `json:"period"`
	IsLive    bool         `json:"isLive"`
	Suspended bool         `json:"suspended"`
	Choices   []OddsChoice `json:"choices"`
}

type OddsChoice struct {
	Name           string  `json:"name"`
	Decimal        float64 `json:"decimal"`
	InitialDecimal float64 `json:"initialDecimal"`
	Winning        bool    `json:"winning"`
	Change         float64 `json:"change"`
}

type Highlight struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Image string `json:"image"`
}
