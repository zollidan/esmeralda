package bot

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zollidan/esmeralda/internal/config"
)

func Run(cfg config.Config) error {
	bot, err := tgbotapi.NewBotAPI(cfg.TelegramBot.Token)
	if err != nil {
		return err
	}

	bot.Debug, err = strconv.ParseBool(cfg.TelegramBot.Debug)
	if err != nil {
		return err
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

	svc := NewRequestClient(cfg.TelegramBot.APIBaseURL, cfg.TelegramBot.APIToken)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			chatID := update.Message.Chat.ID
			msgID := update.Message.MessageID
			text := update.Message.Text

			switch update.Message.Command() {
			case "start":
				if err := replyToStartCommand(bot, chatID, msgID); err != nil {
					log.Println(err)
				}
			case "tasks":
				if err := replyTasks(bot, chatID, msgID, svc); err != nil {
					log.Println(err)
				}
			case "date":
				if err := replyStartTask(bot, chatID, msgID, text, svc); err != nil {
					log.Println(err)
				}
			case "help":
				if err := replyHelp(bot, chatID, msgID); err != nil {
					log.Println(err)
				}
			default:
				if err := replyDefault(bot, chatID, msgID); err != nil {
					log.Println(err)
				}
			}
		}
	}
	return nil
}
