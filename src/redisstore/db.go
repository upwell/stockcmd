package redisstore

import (
	"strconv"

	"github.com/go-redis/redis"
	"hehan.net/my/stockcmd/logger"
	"hehan.net/my/stockcmd/store"
)

var Redis *redis.Client

func init() {
	logger.InitLogger()

	redisAddr, existed, _ := store.RunningConfig.GetString("redisAddr")
	if !existed {
		logger.SugarLog.Error("redis address not set")
		return
	}
	redisDBStr, _ := store.RunningConfig.GetStringOrDefault("redisDB", "0")
	redisDB, _ := strconv.Atoi(redisDBStr)
	options := &redis.Options{
		Addr: redisAddr,
		DB:   redisDB,
	}
	redisPass, existed, _ := store.RunningConfig.GetString("redisPassword")
	if existed {
		options.Password = redisPass
	}

	client := redis.NewClient(options)
	Redis = client
	_, err := Redis.Ping().Result()
	if err != nil {
		logger.SugarLog.Error(err.Error())
	}
}
