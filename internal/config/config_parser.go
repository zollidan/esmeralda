package config

import (
	"flag"
	"log"
)

type Config struct {
    SportAPIRU  SportAPIRU
    Excel       Excel
    Tech        Tech
    DatabaseDSN string
}

type SportAPIRU struct {
    BaseURL         string
    BaseFootballURL string
    Token           string
}

type Excel struct {
    FilePath string
    FileName string
}

type Tech struct {
    GamesLimit int
    Workers int
} 

func InitConfig() Config {
    token := flag.String("token", "", "API token for api-sport.ru")
    dsn := flag.String("dsn", "host=localhost user=esmeralda password=secret dbname=esmeralda port=5432 sslmode=disable", "PostgreSQL DSN")
    flag.Parse()

    if *token == "" {
        log.Fatal("token is required: --token=YOUR_SECRET_TOKEN")
    }

    return Config{
        DatabaseDSN: *dsn,
        SportAPIRU: SportAPIRU{
            BaseURL:         "https://api.api-sport.ru/v2",
            BaseFootballURL: "https://api.api-sport.ru/v2/football",
            Token:           *token,
        },
        Excel: Excel{
            FilePath: "./output/",
            FileName: "esmeralda-ru.xlsx",
        },
        Tech: Tech{
            Workers:    20,
        },
    }
}