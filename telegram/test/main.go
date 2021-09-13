package main

import (
	"github.com/poniteru/go-coin-watcher/app/config"
	"log"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// 此为使用轮询模式示例
func main() {
	bot, err := tgbotapi.NewBotAPI(config.BotToken1)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
		bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text+"1"))
	}
}
