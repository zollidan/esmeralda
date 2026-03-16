package bot

import (
	"log"
	"net/http"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zollidan/esmeralda/internal/config"
)

type Bot struct {
	cfg config.Config
	client *http.Client
}

func NewBot(cfg config.Config) *Bot {
	return &Bot{
		cfg:    cfg,
		client: &http.Client{},
	}
}

func (b *Bot) Run() error {
	bot, err := tgbotapi.NewBotAPI(b.cfg.TelegramBot.Token)
	if err != nil {
		return err
	}

	bot.Debug, err = strconv.ParseBool(b.cfg.TelegramBot.Debug)
	if err != nil {
		return err
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)

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

			case "date":
				if err := b.replyStartTask(bot, chatID, msgID, text); err != nil {
					log.Println(err)
				}

			case "help":
				if err := b.replyHelp(bot, chatID, msgID); err != nil {
					log.Println(err)
				}

			default:
				if err := b.replyDefault(bot, chatID, msgID); err != nil {
					log.Println(err)
				}
			}
		}
	}
	return nil
}
