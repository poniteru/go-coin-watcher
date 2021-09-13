package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// 使用自定义type后 idea的查找实现类功能无法定位到，故放弃
// type MsgHandlerFunc func(update *tgbotapi.Update, botBean *BotBean)

type IBotMsgHandler interface {
	GetMsgHandler(command string) func(update *tgbotapi.Update, botBean *BotBean)
}

type BotMsgHandler struct {
	botEventMap map[string]func(update *tgbotapi.Update, botBean *BotBean)
}

func (botMsgHandler *BotMsgHandler) GetMsgHandler(command string) func(update *tgbotapi.Update, botBean *BotBean) {
	return botMsgHandler.botEventMap[command]
}

func NewBotMsgHandler(botEventMap map[string]func(update *tgbotapi.Update, botBean *BotBean)) *BotMsgHandler {
	return &BotMsgHandler{botEventMap: botEventMap}
}
