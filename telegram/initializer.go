package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/poniteru/go-coin-watcher/app/config"
	"github.com/poniteru/go-coin-watcher/common/constant"
	"log"
)

var BotMap = make(map[string]*BotBean)

func init() {
	initTgBot()
}

func initTgBot() {
	var bots []*BotBean
	bots = append(bots, NewBotBean(config.BotToken1, config.BotHookPath1).SetBotMsgHandler(NewBotMsgHandlerImpl()))
	for _, botBean := range bots {
		BotMap[botBean.BotToken] = botBean
	}
	for _, botBean := range BotMap {
		bot, err := tgbotapi.NewBotAPI(botBean.BotToken)
		if err != nil {
			log.Panic(err)
		}
		botBean.Bot = bot
		bot.Debug = true
		log.Printf("Authorized on account %s", bot.Self.UserName)

		//删除webhook
		//apiResponse, err := bot.RemoveWebhook()
		//log.Printf("%+v\n", apiResponse)
		//if botBean.BotToken == constant.BotToken2 {
		setupWebhook(botBean)
		//}
	}
}

func setupWebhook(botBean *BotBean) {
	_, err := botBean.Bot.SetWebhook(tgbotapi.NewWebhook(config.TgWebHookHost + config.BasePath + constant.V1 + config.TgWebHookPath + botBean.BotHookPath))
	if err != nil {
		log.Fatal(err)
	}
	info, err := botBean.Bot.GetWebhookInfo()
	if err != nil {
		log.Fatal(err)
	}
	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
}

func GetBot(botToken string) *BotBean {
	return BotMap[botToken]
}
