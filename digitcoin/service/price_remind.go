package service

import (
	"github.com/go-redis/redis/v8"
	"github.com/poniteru/go-coin-watcher/digitcoin/market"
	"github.com/poniteru/go-coin-watcher/digitcoin/repository"
)

func SetReminder(currencyPair string, direction market.Direction, price string, member string) (err error) {
	return repository.SetReminder(currencyPair, direction, price, member)
}

// currencyPair btc:usdt
func GetAndDelReminders(currencyPair string, price string) ([]redis.Z, []redis.Z) {
	valUp, err := repository.GetAndDelUpReminders(currencyPair, price)
	if err != nil {
		return nil, nil
	}
	valDown, err := repository.GetAndDelDownReminders(currencyPair, price)
	if err != nil {
		return valUp, nil
	}
	return valUp, valDown
}
