package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
)

type (
	BotBean struct {
		BotToken      string
		BotHookPath   string
		Bot           *tgbotapi.BotAPI
		BotMsgHandler IBotMsgHandler
	}
)

func NewBotBean(botToken string, botHookPath string) *BotBean {
	return &BotBean{BotToken: botToken, BotHookPath: botHookPath}
}

func (botBean *BotBean) SetBotMsgHandler(botMsgHandler IBotMsgHandler) *BotBean {
	botBean.BotMsgHandler = botMsgHandler
	return botBean
}

func (botBean *BotBean) SendMsg(chatID int64, text string) {
	_, err := botBean.Bot.Send(tgbotapi.NewMessage(chatID, text))
	if err != nil {
		log.Printf("%+v\n", err)
	}
}
