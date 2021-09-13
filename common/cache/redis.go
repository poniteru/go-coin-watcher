package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/poniteru/go-coin-watcher/app/config"
	"time"
)

var (
	Rdb *redis.Client
	ctx = context.Background()
)

func init() {
	if err := initClient(); err != nil {
		panic(err)
	}
}

func initClient() (err error) {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPwd, // no password set
		DB:       0,               // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	_, err = Rdb.Ping(ctx).Result()
	return err
}

func V8Example() {
	//if err := initClient(); err != nil {
	//	return
	//}

	err := Rdb.Set(ctx, "key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := Rdb.Get(ctx, "key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

	val2, err := Rdb.Get(ctx, "key2").Result()
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		panic(err)
	} else {
		fmt.Println("key2", val2)
	}
	// Output: key value
	// key2 does not exist
}

func ZADD() {
	val, err := Rdb.ZRangeWithScores(ctx, "mykey", 0, -1).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)
}
