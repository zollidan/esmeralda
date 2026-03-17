CREATE TABLE games (
    id              BIGSERIAL PRIMARY KEY,
    source          VARCHAR(20) NOT NULL DEFAULT 'parser',

    -- базовые поля
    played_at       TIMESTAMPTZ NOT NULL,
    league          TEXT NOT NULL,
    team_home       TEXT NOT NULL,
    team_away       TEXT NOT NULL,
    score_home      SMALLINT,
    score_away      SMALLINT,

    -- букмекерские коэффициенты
    bk_name         TEXT,
    odds_1          NUMERIC(6,2),
    odds_x          NUMERIC(6,2),
    odds_2          NUMERIC(6,2),
    odds_over       NUMERIC(6,2),
    odds_under      NUMERIC(6,2),

    -- статистика очных встреч
    h2h_total       SMALLINT,
    h2h_w           SMALLINT,
    h2h_d           SMALLINT,
    h2h_l           SMALLINT,

    -- домашняя статистика хозяев
    home_total      SMALLINT,
    home_w          SMALLINT,
    home_d          SMALLINT,
    home_l          SMALLINT,

    -- выездная статистика гостей
    away_total      SMALLINT,
    away_w          SMALLINT,
    away_d          SMALLINT,
    away_l          SMALLINT,

    -- пороговая статистика (м_25, м_3, м_5)
    h2h_m25_yes     SMALLINT, h2h_m25_total SMALLINT,
    h2h_m3_yes      SMALLINT, h2h_m3_total  SMALLINT,
    h2h_m5_yes      SMALLINT, h2h_m5_total  SMALLINT,

    home_m25_yes    SMALLINT, home_m25_total SMALLINT,
    home_m3_yes     SMALLINT, home_m3_total  SMALLINT,
    home_m5_yes     SMALLINT, home_m5_total  SMALLINT,

    away_m25_yes    SMALLINT, away_m25_total SMALLINT,
    away_m3_yes     SMALLINT, away_m3_total  SMALLINT,
    away_m5_yes     SMALLINT, away_m5_total  SMALLINT,

    handicap        NUMERIC(5,2),
    created_at      TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE (played_at, team_home, team_away, league, bk_name)
);

CREATE INDEX idx_games_played_at ON games (played_at DESC);
CREATE INDEX idx_games_league    ON games (league);
CREATE INDEX idx_games_teams     ON games (team_home, team_away);