package utils

import (
	"tesla/config"
	"time"
)
import "context"
import "github.com/go-redis/redis/v8"

var (
	RedisClient *redis.Client
)

func init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.RedisUrl,
		Password: "",                  // no password set
		DB:       config.AppConfig.DB, // use default DB
		PoolSize: 50,                  // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}

func GetRedisValueByPrefix(key string) (val string, err error) {
	ctx := context.Background()
	val, err = RedisClient.Get(ctx, config.AppConfig.KeyPrefix+key).Result()
	return
}

func SetRedisValueByPrefix(key string, value string, t time.Duration) (err error) {
	ctx := context.Background()
	_, err = RedisClient.Set(ctx, config.AppConfig.KeyPrefix+key, value, t).Result()
	return
}
