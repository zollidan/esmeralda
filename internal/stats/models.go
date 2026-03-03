package stats

import "time"

type TotalStats struct {
    Matches int
    Win1    int
    Draw    int
    Win2    int
}

type GoalStats struct {
    Over25Matches int
    Over25Total   int
    Over3Matches  int
    Over3Total    int
    Over5Matches  int
    Over5Total    int
}

type Row struct {
    Day      int
    Month    time.Month
    Year     int
    Time     string
    HomeTeam string
    AwayTeam string
    League   string

    H2H15     TotalStats
    H2H25     TotalStats
    Home25    TotalStats
    Away15    TotalStats
    Away25    TotalStats     
    AwayK2_25 TotalStats     

    GoalsH2H  GoalStats
    GoalsHome GoalStats
    GoalsAway GoalStats
}