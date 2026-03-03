package main

import (
	"log"
	"os"

	"github.com/zollidan/esmeralda-ru-api-fetcher/internal/api"
	"github.com/zollidan/esmeralda-ru-api-fetcher/internal/config"
	"github.com/zollidan/esmeralda-ru-api-fetcher/internal/export"
	"github.com/zollidan/esmeralda-ru-api-fetcher/internal/processor"
	"github.com/zollidan/esmeralda-ru-api-fetcher/internal/utils"
)

func main() {

	cfg := config.InitConfig()

	client := api.NewClient(cfg.SportAPIRU.BaseURL, cfg.SportAPIRU.Token)

	date, err := utils.InputDate(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	matches, totalMatches, err := client.GetMatches(api.MatchesFilter{
		Date: date,
	})
	if err != nil {
		log.Fatal(err)
	}

	writer, err := export.NewWriter(cfg.Excel.FileName)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
        if err := writer.Save(cfg.Excel.FileName); err != nil {
            log.Printf("save error: %v", err)
        }
    }()

	err = processor.ProcessMatches(client, writer, matches, totalMatches)
	if err != nil {
		log.Fatal(err)
	}
}

