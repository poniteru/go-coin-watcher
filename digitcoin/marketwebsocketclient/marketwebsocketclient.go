package marketwebsocketclient

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/huobirdcenter/huobi_golang/config"
	"github.com/huobirdcenter/huobi_golang/logging/applogger"
	"github.com/huobirdcenter/huobi_golang/pkg/client/marketwebsocketclient"
	"github.com/huobirdcenter/huobi_golang/pkg/model/market"
	config2 "github.com/poniteru/go-coin-watcher/app/config"
	"github.com/poniteru/go-coin-watcher/digitcoin/digicoinutil"
	"github.com/poniteru/go-coin-watcher/digitcoin/digiconst"
	market2 "github.com/poniteru/go-coin-watcher/digitcoin/market"
	"github.com/poniteru/go-coin-watcher/digitcoin/model"
	"github.com/poniteru/go-coin-watcher/digitcoin/service"
	"github.com/poniteru/go-coin-watcher/telegram"
	"strconv"
	"strings"
)

func RunAll() {
	//reqAndSubscribeCandlestick()
	//reqAndSubscribeDepth()
	//reqAndSubscribe150LevelMarketByPrice()
	//subscribeFullMarketByPrice()
	//reqAndSubscribeMarketByPriceTick()
	//subscribeBBO()
	//reqAndSubscribeTrade()

	reqAndSubscribeLast24hCandlestickAll()
}

func reqAndSubscribeLast24hCandlestickAll() {
	for k, v := range digiconst.CurrencyPairMap {
		reqAndSubscribeLast24hCandlestickV2(k, v)
	}
}

func reqAndSubscribeLast24hCandlestickV2(currencyPair string, clientId string) {

	/*	mockMarketRemindInfo := &model.MarketRemindInfo{
			Direction:   market2.UP,
			TargetPrice: decimal.NewFromFloat(35780.00),
		}

		//mockMarketRemindInfos := make([]*model.MarketRemindInfo, 8)
		mockMarketRemindInfos := []*model.MarketRemindInfo{mockMarketRemindInfo}*/

	// Initialize a new instance
	client := new(marketwebsocketclient.Last24hCandlestickWebSocketClient).Init(config.Host)

	// Set the callback handlers
	client.SetHandler(
		// Connected handler
		func() {
			client.Request(digicoinutil.ToHuobiWsDefaultCurrencyPair(currencyPair), clientId)

			client.Subscribe(digicoinutil.ToHuobiWsDefaultCurrencyPair(currencyPair), clientId)
		},
		// Response handler
		func(resp interface{}) {
			candlestickResponse, ok := resp.(market.SubscribeLast24hCandlestickResponse)
			if ok {
				if &candlestickResponse != nil {
					if candlestickResponse.Tick != nil {
						t := candlestickResponse.Tick
						//applogger.Info("WebSocket received candlestick update, id: %d, count: %v, volume: %v [%v-%v-%v-%v]",
						//	t.Id, t.Count, t.Vol, t.Open, t.Close, t.Low, t.High)
						go remindV2(t, currencyPair)
					}

					if candlestickResponse.Data != nil {
						t := candlestickResponse.Data
						//applogger.Info("WebSocket received candlestick data, id: %d, count: %v, volume: %v [%v-%v-%v-%v]",
						//	t.Id, t.Count, t.Vol, t.Open, t.Close, t.Low, t.High)
						go remindV2(t, currencyPair)
					}
				}
			} else {
				applogger.Warn("Unknown response: %v", resp)
			}
		})

	// Connect to the server and wait for the handler to handle the response
	client.Connect(true)

	//fmt.Println("Press ENTER to unsubscribe and stop...")
	//fmt.Scanln()
	//
	//client.UnSubscribe("btcusdt", "1608")
	//
	//client.Close()
	//applogger.Info("Client closed")
}

func remind(mockMarketRemindInfos *[]*model.MarketRemindInfo, t *market.Candlestick) {
	for i, info := range *mockMarketRemindInfos {
		if info.Direction == market2.UP {
			if t.Close.GreaterThanOrEqual(info.TargetPrice) {
				text := fmt.Sprintf("提醒条件已触发UP: %s", info.TargetPrice)
				applogger.Info(text)
				*mockMarketRemindInfos = append((*mockMarketRemindInfos)[:i], (*mockMarketRemindInfos)[i+1:]...)
			}
		} else if info.Direction == market2.DOWN {
			if t.Close.LessThanOrEqual(info.TargetPrice) {
				text := fmt.Sprintf("提醒条件已触发DOWN: %s", info.TargetPrice)
				applogger.Info(text)
				*mockMarketRemindInfos = append((*mockMarketRemindInfos)[:i], (*mockMarketRemindInfos)[i+1:]...)
			}
		}
	}
}

func remindV2(t *market.Candlestick, currencyPair string) {
	up, down := service.GetAndDelReminders(digicoinutil.ToDefaultCurrencyPair(currencyPair), t.Close.String())
	for _, z := range up {
		//fmt.Printf("%s\n", z.Member.(string))
		//text := fmt.Sprintf("%v设置的升到%v的提醒条件已触发", z.Member, z.Score)
		//applogger.Info(text)
		sendReminder(currencyPair, market2.UP, &z, t)
	}
	for _, z := range down {
		//fmt.Printf("%s\n", z.Member.(string))
		//text := fmt.Sprintf("%v设置的降到%v的提醒条件已触发", z.Member, z.Score)
		//applogger.Info(text)
		sendReminder(currencyPair, market2.DOWN, &z, t)
	}
}

func sendReminder(currencyPair string, direction market2.Direction, z *redis.Z, t *market.Candlestick) {
	var logText string
	var remindText string = "%v 已%s到 %v, 当前价格: %v"
	if direction == market2.UP {
		logText = fmt.Sprintf("[%v]设置的[%v]涨到[%v]的提醒条件已触发", z.Member, currencyPair, z.Score)
		remindText = fmt.Sprintf(remindText, digicoinutil.ToDisplayCurrencyPair(currencyPair), "涨", z.Score, t.Close.String())
	} else if direction == market2.DOWN {
		logText = fmt.Sprintf("[%v]设置的[%v]跌到[%v]的提醒条件已触发", z.Member, currencyPair, z.Score)
		remindText = fmt.Sprintf(remindText, digicoinutil.ToDisplayCurrencyPair(currencyPair), "跌", z.Score, t.Close.String())
	}
	applogger.Info(logText)
	memberStr := z.Member.(string)
	if strings.Contains(memberStr, "_") {
		memberStr = strings.Split(memberStr, "_")[0]
	}
	chatId, err := strconv.ParseInt(memberStr, 10, 64)
	if err != nil {
		return
	}
	//applogger.Info("[%v],[%v]", chatId, remindText)
	telegram.GetBot(config2.BotToken1).SendMsg(chatId, remindText)
}
