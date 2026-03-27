package models

import "time"

type Game struct {
	ID       uint    `gorm:"primaryKey"      json:"id"`
	TaskID   *string `gorm:"type:uuid;index" json:"task_id"`
	Day      int     `gorm:"not null"        json:"day"`
	Month    int     `gorm:"not null"        json:"month"`
	Year     int     `gorm:"not null"        json:"year"`
	Time     string  `gorm:"not null"        json:"time"`
	HomeTeam string  `gorm:"not null"        json:"home_team"`
	AwayTeam string  `gorm:"not null"        json:"away_team"`
	League   string  `gorm:"not null"        json:"league"`

	// Сумма очных игр на поле команда 1 + W/D/L
	H2HHomeMatches int `gorm:"not null;default:0" json:"h2h_home_matches"`
	H2HHomeWin1    int `gorm:"not null;default:0" json:"h2h_home_win1"`
	H2HHomeDraw    int `gorm:"not null;default:0" json:"h2h_home_draw"`
	H2HHomeWin2    int `gorm:"not null;default:0" json:"h2h_home_win2"`

	// Общее количество матчей дома первой команды + W/D/L
	HomeTeamHomeMatches int `gorm:"not null;default:0" json:"home_team_home_matches"`
	HomeTeamHomeWin     int `gorm:"not null;default:0" json:"home_team_home_win"`
	HomeTeamHomeDraw    int `gorm:"not null;default:0" json:"home_team_home_draw"`
	HomeTeamHomeLoss    int `gorm:"not null;default:0" json:"home_team_home_loss"`

	// Общее количество матчей в гостях второй команды + L/D/W
	AwayTeamAwayMatches int `gorm:"not null;default:0" json:"away_team_away_matches"`
	AwayTeamAwayLoss    int `gorm:"not null;default:0" json:"away_team_away_loss"`
	AwayTeamAwayDraw    int `gorm:"not null;default:0" json:"away_team_away_draw"`
	AwayTeamAwayWin     int `gorm:"not null;default:0" json:"away_team_away_win"`

	// Тотал 2.5: очные на поле хозяев
	H2HHomeOver25Total int `gorm:"not null;default:0" json:"h2h_home_over25_total"`
	H2HHomeOver25Over  int `gorm:"not null;default:0" json:"h2h_home_over25_over"`
	H2HHomeOver25Under int `gorm:"not null;default:0" json:"h2h_home_over25_under"`

	// Тотал 2.5: все встречи (команда 1 дома)
	HomeAllOver25Total int `gorm:"not null;default:0" json:"home_all_over25_total"`
	HomeAllOver25Over  int `gorm:"not null;default:0" json:"home_all_over25_over"`
	HomeAllOver25Under int `gorm:"not null;default:0" json:"home_all_over25_under"`

	// Тотал 2.5: все встречи (команда 2 в гостях)
	AwayAllOver25Total int `gorm:"not null;default:0" json:"away_all_over25_total"`
	AwayAllOver25Over  int `gorm:"not null;default:0" json:"away_all_over25_over"`
	AwayAllOver25Under int `gorm:"not null;default:0" json:"away_all_over25_under"`

	// Кол-во игр и сумма мячей: очн. на любом поле (25/5/3)
	H2HAnyGames25 int `gorm:"not null;default:0" json:"h2h_any_games_25"`
	H2HAnyGoals25 int `gorm:"not null;default:0" json:"h2h_any_goals_25"`
	H2HAnyGames5  int `gorm:"not null;default:0" json:"h2h_any_games_5"`
	H2HAnyGoals5  int `gorm:"not null;default:0" json:"h2h_any_goals_5"`
	H2HAnyGames3  int `gorm:"not null;default:0" json:"h2h_any_games_3"`
	H2HAnyGoals3  int `gorm:"not null;default:0" json:"h2h_any_goals_3"`

	// Кол-во игр и сумма мячей: хозяева на любом поле (25/5/3)
	HomeAnyGames25 int `gorm:"not null;default:0" json:"home_any_games_25"`
	HomeAnyGoals25 int `gorm:"not null;default:0" json:"home_any_goals_25"`
	HomeAnyGames5  int `gorm:"not null;default:0" json:"home_any_games_5"`
	HomeAnyGoals5  int `gorm:"not null;default:0" json:"home_any_goals_5"`
	HomeAnyGames3  int `gorm:"not null;default:0" json:"home_any_games_3"`
	HomeAnyGoals3  int `gorm:"not null;default:0" json:"home_any_goals_3"`

	// Кол-во игр и сумма мячей: гости на любом поле (25/5/3)
	AwayAnyGames25 int `gorm:"not null;default:0" json:"away_any_games_25"`
	AwayAnyGoals25 int `gorm:"not null;default:0" json:"away_any_goals_25"`
	AwayAnyGames5  int `gorm:"not null;default:0" json:"away_any_games_5"`
	AwayAnyGoals5  int `gorm:"not null;default:0" json:"away_any_goals_5"`
	AwayAnyGames3  int `gorm:"not null;default:0" json:"away_any_games_3"`
	AwayAnyGoals3  int `gorm:"not null;default:0" json:"away_any_goals_3"`

	// Кол-во игр и сумма мячей: очн. на поле хозяев (25/5/3)
	H2HHomeGames25 int `gorm:"not null;default:0" json:"h2h_home_games_25"`
	H2HHomeGoals25 int `gorm:"not null;default:0" json:"h2h_home_goals_25"`
	H2HHomeGames5  int `gorm:"not null;default:0" json:"h2h_home_games_5"`
	H2HHomeGoals5  int `gorm:"not null;default:0" json:"h2h_home_goals_5"`
	H2HHomeGames3  int `gorm:"not null;default:0" json:"h2h_home_games_3"`
	H2HHomeGoals3  int `gorm:"not null;default:0" json:"h2h_home_goals_3"`

	// Кол-во игр и сумма мячей: хозяева на поле хозяев (25/5/3)
	HomeHomeGames25 int `gorm:"not null;default:0" json:"home_home_games_25"`
	HomeHomeGoals25 int `gorm:"not null;default:0" json:"home_home_goals_25"`
	HomeHomeGames5  int `gorm:"not null;default:0" json:"home_home_games_5"`
	HomeHomeGoals5  int `gorm:"not null;default:0" json:"home_home_goals_5"`
	HomeHomeGames3  int `gorm:"not null;default:0" json:"home_home_games_3"`
	HomeHomeGoals3  int `gorm:"not null;default:0" json:"home_home_goals_3"`

	// Кол-во игр и сумма мячей: гости на поле гостей (25/5/3)
	AwayAwayGames25 int `gorm:"not null;default:0" json:"away_away_games_25"`
	AwayAwayGoals25 int `gorm:"not null;default:0" json:"away_away_goals_25"`
	AwayAwayGames5  int `gorm:"not null;default:0" json:"away_away_games_5"`
	AwayAwayGoals5  int `gorm:"not null;default:0" json:"away_away_goals_5"`
	AwayAwayGames3  int `gorm:"not null;default:0" json:"away_away_games_3"`
	AwayAwayGoals3  int `gorm:"not null;default:0" json:"away_away_goals_3"`

	CreatedAt time.Time `json:"created_at"`
}
