package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/huobirdcenter/huobi_golang/cmd/marketwebsocketclientexample"
	"github.com/huobirdcenter/huobi_golang/config"
	"github.com/huobirdcenter/huobi_golang/logging/applogger"
	"github.com/huobirdcenter/huobi_golang/pkg/client/marketwebsocketclient"
	"github.com/huobirdcenter/huobi_golang/pkg/model/market"
	"github.com/poniteru/go-coin-watcher/app/dao"
	"github.com/poniteru/go-coin-watcher/digitcoin"
	market2 "github.com/poniteru/go-coin-watcher/digitcoin/market"
	"github.com/poniteru/go-coin-watcher/digitcoin/repository"
	"github.com/poniteru/go-coin-watcher/digitcoin/service"
	"strings"
)

func main() {
	//runAll()
	//reqAndSubscribeLast24hCandlestick()
	//test3()
	//test6()
	//test4()
	//test5()
	digitcoin.RunAll()

	fmt.Println("Press ENTER to unsubscribe and stop...")
	fmt.Scanln()
}

func test3() {
	up, down, err := repository.GetReminders("btc:usdt", "10012.5")
	if err != nil {
		return
	}
	for _, z := range up {
		fmt.Printf("%v - %v\n", z.Member, z.Score)

	}
	for _, z := range down {
		fmt.Printf("%v - %v\n", z.Member, z.Score)

	}
}

func test6() {
	up, down := service.GetAndDelReminders("btc:usdt", "10012.5")
	for _, z := range up {
		fmt.Printf("%v - %v\n", z.Member, z.Score)

	}
	for _, z := range down {
		fmt.Printf("%v - %v\n", z.Member, z.Score)

	}
}

func test4() {
	err := repository.SetReminder("btc:usdt", market2.DOWN, "18003.88", "12345678")
	if err != nil {
		return
	}
}

func test5() {
	err := repository.DelReminders("btc:usdt", "10012.5")
	if err != nil {
		return
	}
}

func test2() {
	uuidStr := strings.ReplaceAll(uuid.New().String(), "-", "")
	dao.InsertUser(12365678, uuidStr)
}

func runAll() {
	//commonclientexample.RunAllExamples()
	//accountclientexample.RunAllExamples()
	//orderclientexample.RunAllExamples()
	//algoorderclientexample.RunAllExamples()
	//marketclientexample.RunAllExamples()
	//isolatedmarginclientexample.RunAllExamples()
	//crossmarginclientexample.RunAllExamples()
	//walletclientexample.RunAllExamples()
	//subuserclientexample.RunAllExamples()
	//stablecoinclientexample.RunAllExamples()
	//etfclientexample.RunAllExamples()
	marketwebsocketclientexample.RunAllExamples()
	//accountwebsocketclientexample.RunAllExamples()
	//orderwebsocketclientexample.RunAllExamples()
}

func reqAndSubscribeLast24hCandlestick() {
	// Initialize a new instance
	client := new(marketwebsocketclient.Last24hCandlestickWebSocketClient).Init(config.Host)

	// Set the callback handlers
	client.SetHandler(
		// Connected handler
		func() {
			client.Request("btcusdt", "1608")

			client.Subscribe("btcusdt", "1608")
		},
		// Response handler
		func(resp interface{}) {
			candlestickResponse, ok := resp.(market.SubscribeLast24hCandlestickResponse)
			if ok {
				if &candlestickResponse != nil {
					if candlestickResponse.Tick != nil {
						t := candlestickResponse.Tick
						applogger.Info("WebSocket received candlestick update, id: %d, count: %v, volume: %v [%v-%v-%v-%v]",
							t.Id, t.Count, t.Vol, t.Open, t.Close, t.Low, t.High)
					}

					if candlestickResponse.Data != nil {
						t := candlestickResponse.Data
						applogger.Info("WebSocket received candlestick data, id: %d, count: %v, volume: %v [%v-%v-%v-%v]",
							t.Id, t.Count, t.Vol, t.Open, t.Close, t.Low, t.High)
					}
				}
			} else {
				applogger.Warn("Unknown response: %v", resp)
			}
		})

	// Connect to the server and wait for the handler to handle the response
	client.Connect(true)

	fmt.Println("Press ENTER to unsubscribe and stop...")
	fmt.Scanln()

	client.UnSubscribe("btcusdt", "1608")

	client.Close()
	applogger.Info("Client closed")
}
