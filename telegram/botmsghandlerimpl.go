package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/poniteru/go-coin-watcher/app/dao"
	"github.com/poniteru/go-coin-watcher/app/service"
	"github.com/poniteru/go-coin-watcher/digitcoin/digicoinutil"
	"github.com/poniteru/go-coin-watcher/digitcoin/digiconst"
	market2 "github.com/poniteru/go-coin-watcher/digitcoin/market"
	service2 "github.com/poniteru/go-coin-watcher/digitcoin/service"
	"strconv"
	"strings"
)

func NewBotMsgHandlerImpl() *BotMsgHandler {
	botEventMap := make(map[string]func(update *tgbotapi.Update, botBean *BotBean))
	botEventMap["/start"] = Start
	botEventMap["/ping"] = Ping
	botEventMap["/invitecode"] = InviteCodeInput
	botEventMap["/setreminder"] = SetReminder
	botEventMap["/listreminder"] = ListReminder
	botEventMap["/delreminder"] = DelReminder
	return NewBotMsgHandler(botEventMap)
}

func Start(update *tgbotapi.Update, botBean *BotBean) {
	if "/start" == update.Message.Text {
		text := fmt.Sprintf("Hello %s,\nYour chat ID is:%d", update.Message.From.String(), update.Message.Chat.ID)
		botBean.SendMsg(update.Message.Chat.ID, text)
	} else if strings.HasPrefix(update.Message.Text, "/start ") {
		token, err := service.Register(update.Message.Chat.ID, strings.TrimPrefix(update.Message.Text, "/start "))
		if err != nil {
			return
		}
		text := fmt.Sprintf("Hello %s,\nYour chat ID is:%d,\nYour token is:%s", update.Message.From.String(), update.Message.Chat.ID, token)
		botBean.SendMsg(update.Message.Chat.ID, text)
	}
}

func Ping(update *tgbotapi.Update, botBean *BotBean) {
	if "/ping" == update.Message.Text {
		text := fmt.Sprintf("pong")
		botBean.SendMsg(update.Message.Chat.ID, text)
	}
}

func InviteCodeInput(update *tgbotapi.Update, botBean *BotBean) {
	if strings.HasPrefix(update.Message.Text, "/invitecode ") {
		// registerService
		token, err := service.Register(update.Message.Chat.ID, strings.TrimPrefix(update.Message.Text, "/invitecode "))
		if err != nil {
			return
		}
		text := fmt.Sprintf("Hello %s,\nYour token is:%s", update.Message.From.String(), token)
		botBean.SendMsg(update.Message.Chat.ID, text)
	}
}

func SetReminder(update *tgbotapi.Update, botBean *BotBean) {
	exists, err := dao.SelectUserExists(update.Message.Chat.ID)
	if err != nil || exists != 1 {
		botBean.SendMsg(update.Message.Chat.ID, "操作失败")
		return
	}
	// /setreminder eth up 1234.56
	cmds := strings.Fields(update.Message.Text)
	if len(cmds) < 4 {
		botBean.SendMsg(update.Message.Chat.ID, "操作失败")
		return
	}
	coinPair := digicoinutil.ToDefaultCurrencyPair(cmds[1])
	if _, ok := digiconst.CurrencyPairMap[coinPair]; !ok {
		botBean.SendMsg(update.Message.Chat.ID, "操作失败")
		return
	}
	direction := market2.GetDirection(cmds[2])
	if direction == market2.UNDEFINED {
		botBean.SendMsg(update.Message.Chat.ID, "操作失败")
		return
	}
	priceStr := cmds[3]
	_, err = strconv.ParseFloat(priceStr, 64)
	if err != nil {
		botBean.SendMsg(update.Message.Chat.ID, "操作失败")
		return
	}
	err = service2.SetReminder(digicoinutil.ToDefaultCurrencyPair(coinPair), direction, priceStr, strconv.FormatInt(update.Message.Chat.ID, 10))
	if err != nil {
		botBean.SendMsg(update.Message.Chat.ID, "操作失败")
		return
	}
	text := fmt.Sprintf("操作成功")
	botBean.SendMsg(update.Message.Chat.ID, text)
}

func ListReminder(update *tgbotapi.Update, botBean *BotBean) {
	// 检查用户是否存在
	exists, err := dao.SelectUserExists(update.Message.Chat.ID)
	if err != nil || exists != 1 {
		botBean.SendMsg(update.Message.Chat.ID, "操作失败")
		return
	}
	// 检查命令参数是否正确
	// /listreminder eth up
	cmds := strings.Fields(update.Message.Text)
	if len(cmds) < 3 {
		botBean.SendMsg(update.Message.Chat.ID, "操作失败")
		return
	}
	coinPair := digicoinutil.ToDefaultCurrencyPair(cmds[1])
	if _, ok := digiconst.CurrencyPairMap[coinPair]; !ok {
		botBean.SendMsg(update.Message.Chat.ID, "操作失败")
		return
	}
	direction := market2.GetDirection(cmds[2])
	if direction == market2.UNDEFINED {
		botBean.SendMsg(update.Message.Chat.ID, "操作失败")
		return
	}
	reminders, err := service2.GetReminders(coinPair, direction, strconv.FormatInt(update.Message.Chat.ID, 10))
	if err != nil {
		botBean.SendMsg(update.Message.Chat.ID, "操作失败")
		return
	}
	if reminders == nil {
		botBean.SendMsg(update.Message.Chat.ID, "无")
		return
	}
	var builder strings.Builder
	for _, z := range reminders {
		builder.WriteString(strconv.FormatFloat(z.Score, 'f', -1, 64) + "\n")
	}
	//str := builder.String()
	botBean.SendMsg(update.Message.Chat.ID, builder.String())
}

func DelReminder(update *tgbotapi.Update, botBean *BotBean) {
	// 检查用户是否存在
	exists, err := dao.SelectUserExists(update.Message.Chat.ID)
	if err != nil || exists != 1 {
		botBean.SendMsg(update.Message.Chat.ID, "操作失败")
		return
	}
	// 检查命令参数是否正确
	// /delreminder eth up 1234.1234
	cmds := strings.Fields(update.Message.Text)
	if len(cmds) < 4 {
		botBean.SendMsg(update.Message.Chat.ID, "操作失败")
		return
	}
	coinPair := digicoinutil.ToDefaultCurrencyPair(cmds[1])
	if _, ok := digiconst.CurrencyPairMap[coinPair]; !ok {
		botBean.SendMsg(update.Message.Chat.ID, "操作失败")
		return
	}
	direction := market2.GetDirection(cmds[2])
	if direction == market2.UNDEFINED {
		botBean.SendMsg(update.Message.Chat.ID, "操作失败")
		return
	}
	priceStr := cmds[3]
	_, err = strconv.ParseFloat(priceStr, 64)
	if err != nil {
		botBean.SendMsg(update.Message.Chat.ID, "操作失败")
		return
	}
	err = service2.DelReminder(coinPair, direction, priceStr, strconv.FormatInt(update.Message.Chat.ID, 10))
	if err != nil {
		botBean.SendMsg(update.Message.Chat.ID, "操作失败")
		return
	}
	botBean.SendMsg(update.Message.Chat.ID, "操作成功")
}
