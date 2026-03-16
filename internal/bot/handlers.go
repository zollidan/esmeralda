package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/zollidan/esmeralda/internal/queue"
)

const (
	taskPollInterval = 10 * time.Second
	taskWaitTimeout  = 30 * time.Minute
)

func replyToStartCommand(bot *tgbotapi.BotAPI, chatID int64, msgID int) error {
	msg := tgbotapi.NewMessage(chatID, "Welcome to Esmeralda Bot! Use /tasks to see your tasks, /date to start a task, and /help for more information.")
	msg.ReplyToMessageID = msgID

	bot.Send(msg)

	return nil
}

func (b *Bot) replyStartTask(bot *tgbotapi.BotAPI, chatID int64, msgID int, text string) error {
	parts := strings.Split(text, " ")

	if len(parts) < 2 {
		msg := tgbotapi.NewMessage(chatID, "Usage: /date YYYY-MM-DD")
		msg.ReplyToMessageID = msgID
		_, err := bot.Send(msg)
		return err
	}

	date := parts[1]

	task, err := b.CreateTask(date)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Failed to create task")
		msg.ReplyToMessageID = msgID
		bot.Send(msg)
		return err
	}

	msg := tgbotapi.NewMessage(chatID, fmt.Sprintf("Task created: %s\nDate: %s\nWaiting for result...", task.ID, task.Date))
	msg.ReplyToMessageID = msgID

	if _, err = bot.Send(msg); err != nil {
		return err
	}

	go b.waitAndSendResult(bot, chatID, task.ID, task.Date)

	return nil
}

func (b *Bot) replyHelp(bot *tgbotapi.BotAPI, chatID int64, msgID int) error {
	msg := tgbotapi.NewMessage(chatID,
		"Commands:\n"+
			"/date YYYY-MM-DD — create task\n"+
			"/tasks — show tasks",
	)

	msg.ReplyToMessageID = msgID
	_, err := bot.Send(msg)
	return err
}

func (b *Bot) replyDefault(bot *tgbotapi.BotAPI, chatID int64, msgID int) error {
	msg := tgbotapi.NewMessage(chatID, "Unknown command. Use /help")
	msg.ReplyToMessageID = msgID

	_, err := bot.Send(msg)
	return err
}

func (b *Bot) waitAndSendResult(bot *tgbotapi.BotAPI, chatID int64, taskID, date string) {
	task, err := b.WaitForTask(taskID, taskPollInterval, taskWaitTimeout)
	if err != nil {
		log.Printf("wait for task %s: %v", taskID, err)
		b.sendText(bot, chatID, fmt.Sprintf("Task %s: failed while waiting for result: %v", taskID, err))
		return
	}

	switch task.Status {
	case string(queue.StatusDone):
		content, filename, err := b.DownloadExport(date)
		if err != nil {
			log.Printf("download export for task %s: %v", taskID, err)
			b.sendText(bot, chatID, fmt.Sprintf("Task %s finished, but export download failed: %v", taskID, err))
			return
		}

		doc := tgbotapi.NewDocument(chatID, tgbotapi.FileBytes{Name: filename, Bytes: content})
		doc.Caption = fmt.Sprintf("Task %s completed successfully", taskID)
		if _, err := bot.Send(doc); err != nil {
			log.Printf("send document for task %s: %v", taskID, err)
			b.sendText(bot, chatID, fmt.Sprintf("Task %s completed, but file delivery failed: %v", taskID, err))
		}
	case string(queue.StatusError), string(queue.StatusFailed), string(queue.StatusCancelled):
		b.sendText(bot, chatID, fmt.Sprintf("Task %s finished with status: %s", taskID, task.Status))
	default:
		b.sendText(bot, chatID, fmt.Sprintf("Task %s finished with unexpected status: %s", taskID, task.Status))
	}
}

func (b *Bot) sendText(bot *tgbotapi.BotAPI, chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := bot.Send(msg); err != nil {
		log.Printf("send telegram message: %v", err)
	}
}
