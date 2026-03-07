package models

import "time"

type TotalStats struct {
	Matches int `gorm:"not null;default:0" json:"matches"`
	Win1    int `gorm:"not null;default:0" json:"win1"`
	Draw    int `gorm:"not null;default:0" json:"draw"`
	Win2    int `gorm:"not null;default:0" json:"win2"`
}

type GoalStats struct {
	Over25Matches int `gorm:"not null;default:0" json:"over_25_matches"`
	Over25Total   int `gorm:"not null;default:0" json:"over_25_total"`
	Over3Matches  int `gorm:"not null;default:0" json:"over_3_matches"`
	Over3Total    int `gorm:"not null;default:0" json:"over_3_total"`
	Over5Matches  int `gorm:"not null;default:0" json:"over_5_matches"`
	Over5Total    int `gorm:"not null;default:0" json:"over_5_total"`
}

type Game struct {
	ID       uint   `gorm:"primaryKey"                json:"id"`
	TaskID   *string `gorm:"type:uuid;index"           json:"task_id"`
	Day      int    `gorm:"not null"                  json:"day"`
	Month    int    `gorm:"not null"                  json:"month"`
	Year     int    `gorm:"not null"                  json:"year"`
	Time     string `gorm:"not null"                  json:"time"`
	HomeTeam string `gorm:"not null"                  json:"home_team"`
	AwayTeam string `gorm:"not null"                  json:"away_team"`
	League   string `gorm:"not null"                  json:"league"`

	H2H15     TotalStats `gorm:"embedded;embeddedPrefix:h2h15_"      json:"h2h15"`
	H2H25     TotalStats `gorm:"embedded;embeddedPrefix:h2h25_"      json:"h2h25"`
	Home25    TotalStats `gorm:"embedded;embeddedPrefix:home25_"     json:"home25"`
	Away15    TotalStats `gorm:"embedded;embeddedPrefix:away15_"     json:"away15"`
	Away25    TotalStats `gorm:"embedded;embeddedPrefix:away25_"     json:"away25"`
	AwayK2_25 TotalStats `gorm:"embedded;embeddedPrefix:away_k2_25_" json:"away_k2_25"`

	GoalsH2H  GoalStats `gorm:"embedded;embeddedPrefix:goals_h2h_"  json:"goals_h2h"`
	GoalsHome GoalStats `gorm:"embedded;embeddedPrefix:goals_home_" json:"goals_home"`
	GoalsAway GoalStats `gorm:"embedded;embeddedPrefix:goals_away_" json:"goals_away"`

	CreatedAt time.Time `json:"created_at"`

	Task Task `gorm:"foreignKey:TaskID;references:ID" json:"-"`
}