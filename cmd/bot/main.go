package main

import (
	"log"

	"github.com/zollidan/esmeralda/internal/bot"
	"github.com/zollidan/esmeralda/internal/config"
)

func main() {
	cfg := config.InitConfig()

	b := bot.NewBot(cfg)

	err := b.Run()
	if err != nil {
		log.Fatal("Error running bot: ", err)
	}
}
