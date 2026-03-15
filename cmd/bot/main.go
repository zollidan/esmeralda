package main

import (
	"log"

	"github.com/zollidan/esmeralda/internal/bot"
	"github.com/zollidan/esmeralda/internal/config"
)

func main() {

	cfg := config.InitConfig()

	err := bot.Run(cfg)
	if err != nil {
		log.Fatal(err)
	}

}