package model

import (
	"github.com/poniteru/go-coin-watcher/digitcoin/market"
	"github.com/shopspring/decimal"
)

type MarketRemindInfo struct {
	Direction   market.Direction `json:"direction"`
	TargetPrice decimal.Decimal  `json:"targetPrice"`
}
