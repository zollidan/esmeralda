package main

import (
	"log"
	"time"

	"github.com/zollidan/esmeralda/internal/api"
	"github.com/zollidan/esmeralda/internal/config"
	"github.com/zollidan/esmeralda/internal/db"
	"github.com/zollidan/esmeralda/internal/processor"
)

// в бд сохраеяются значения на все матчи нулями, сделать фикс
// добавить другой конфиг

func main() {

	cfg := config.InitConfig()

	client := api.NewClient(cfg.SportAPIRU.BaseURL, cfg.SportAPIRU.Token)

	database, err := db.New(cfg.DatabaseDSN)
	if err != nil {
		log.Fatal(err)
	}

	// date, err := utils.InputDate(os.Stdin)
	// if err != nil {
	// 	log.Fatal(err)
	// }


	// hardcodeed date for testing
	date, err := time.Parse("02.01.2006", "03.03.2026")
	if err != nil {
		log.Fatal(err)
	}

	matches, totalMatches, err := client.GetMatches(api.MatchesFilter{
		Date: date,
	})
	if err != nil {
		log.Fatal(err)
	}

	err = processor.ProcessMatches(client, database, matches, totalMatches)
	if err != nil {
		log.Fatal(err)
	}
}

