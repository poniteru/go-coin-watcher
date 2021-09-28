package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/poniteru/go-coin-watcher/common/cache"
	"github.com/poniteru/go-coin-watcher/digitcoin/market"
	"math"
	"strconv"
	"strings"
	"sync"
)

var locks sync.Map

func GetMarketRemindersKey(currencyPair string, direction market.Direction) string {
	return fmt.Sprintf("marketReminders:%s", strings.ToLower(currencyPair+":"+direction.String()))
}

func SetReminder(currencyPair string, direction market.Direction, price string, member string) (err error) {
	score, err := strconv.ParseFloat(price, 64)
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var val int64
		lock, _ := locks.LoadOrStore(GetMarketRemindersKey(currencyPair, direction), new(sync.Mutex))
		lock.(*sync.Mutex).Lock()
		val, err = cache.Rdb.ZAdd(context.Background(), GetMarketRemindersKey(currencyPair, direction), &redis.Z{
			Score:  score,
			Member: member + "_" + price,
		}).Result()
		lock.(*sync.Mutex).Unlock()
		fmt.Println("val", val)
	}()
	wg.Wait()
	return err
}

func GetRemindersByRange(currencyPair string, direction market.Direction, start, stop int64) ([]redis.Z, error) {
	val, err := cache.Rdb.ZRangeWithScores(context.Background(), GetMarketRemindersKey(currencyPair, direction), start, stop).Result()
	if err != nil {
		return nil, err
	}
	return val, err
}

func GetReminders(currencyPair string, price string) ([]redis.Z, []redis.Z, error) {
	valUp, err := GetUpReminders(currencyPair, price)
	if err != nil {
		return nil, nil, err
	}
	valDown, err := GetDownReminders(currencyPair, price)
	if err != nil {
		return nil, nil, err
	}
	//fmt.Println("valUp", valUp)
	//fmt.Println("valDown", valDown)
	return valUp, valDown, err
}

func GetUpReminders(currencyPair string, price string) ([]redis.Z, error) {
	valUp, err := cache.Rdb.ZRangeByScoreWithScores(context.Background(), GetMarketRemindersKey(currencyPair, market.UP), &redis.ZRangeBy{
		Min: "0",   // 最小分数
		Max: price, // 最大分数
	}).Result()
	if err != nil {
		return nil, err
	}
	fmt.Println("valUp", valUp)
	return valUp, err
}

func GetDownReminders(currencyPair string, price string) ([]redis.Z, error) {
	valDown, err := cache.Rdb.ZRangeByScoreWithScores(context.Background(), GetMarketRemindersKey(currencyPair, market.DOWN), &redis.ZRangeBy{
		Min: price,                       // 最小分数
		Max: strconv.Itoa(math.MaxInt64), // 最大分数
	}).Result()
	if err != nil {
		return nil, err
	}
	fmt.Println("valDown", valDown)
	return valDown, err
}

func DelReminders(currencyPair string, price string) error {
	valUp, err := DelUpReminders(currencyPair, price)
	if err != nil {
		return err
	}
	fmt.Println("DelUp:", valUp)
	valDown, err := DelDownReminders(currencyPair, price)
	if err != nil {
		return err
	}
	fmt.Println("DelDown", valDown)
	return err
}

func DelUpReminders(currencyPair string, price string) (result int64, err error) {
	result, err = cache.Rdb.ZRemRangeByScore(context.Background(), GetMarketRemindersKey(currencyPair, market.UP),
		"0", price).Result()
	return result, err
}

func DelDownReminders(currencyPair string, price string) (result int64, err error) {
	result, err = cache.Rdb.ZRemRangeByScore(context.Background(), GetMarketRemindersKey(currencyPair, market.DOWN),
		price, strconv.Itoa(math.MaxInt64)).Result()
	return result, err
}

func GetAndDelUpReminders(currencyPair string, price string) (valUp []redis.Z, err error) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		lock, _ := locks.LoadOrStore(GetMarketRemindersKey(currencyPair, market.UP), new(sync.Mutex))
		lock.(*sync.Mutex).Lock()
		valUp, err = cache.Rdb.ZRangeByScoreWithScores(context.Background(), GetMarketRemindersKey(currencyPair, market.UP), &redis.ZRangeBy{
			Min: "0",   // 最小分数
			Max: price, // 最大分数
		}).Result()
		if err != nil {
			lock.(*sync.Mutex).Unlock()
			return
		}
		if len(valUp) != 0 {
			_, err = cache.Rdb.ZRemRangeByScore(context.Background(), GetMarketRemindersKey(currencyPair, market.UP),
				"0", price).Result()
		}
		lock.(*sync.Mutex).Unlock()
	}()
	wg.Wait()
	if err != nil {
		return nil, err
	}
	return valUp, err
}

func GetAndDelDownReminders(currencyPair string, price string) (valDown []redis.Z, err error) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		lock, _ := locks.LoadOrStore(GetMarketRemindersKey(currencyPair, market.DOWN), new(sync.Mutex))
		lock.(*sync.Mutex).Lock()
		valDown, err = cache.Rdb.ZRangeByScoreWithScores(context.Background(), GetMarketRemindersKey(currencyPair, market.DOWN), &redis.ZRangeBy{
			Min: price,                       // 最小分数
			Max: strconv.Itoa(math.MaxInt64), // 最大分数
		}).Result()
		if err != nil {
			lock.(*sync.Mutex).Unlock()
			return
		}
		if len(valDown) != 0 {
			_, err = cache.Rdb.ZRemRangeByScore(context.Background(), GetMarketRemindersKey(currencyPair, market.DOWN),
				price, strconv.Itoa(math.MaxInt64)).Result()
		}
		lock.(*sync.Mutex).Unlock()
	}()
	wg.Wait()
	if err != nil {
		return nil, err
	}
	return valDown, err
}
