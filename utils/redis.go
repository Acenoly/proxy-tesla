package utils

import (
	"tesla/config"
	"time"
)
import "context"
import "github.com/go-redis/redis/v8"

var (
	RedisClient *redis.Client
	RedisWriteClient *redis.Client
	RedisSessionWriteClient *redis.Client
	RedisSessionClient *redis.Client

)

func init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.RedisUrl,
		Password: "",                  // no password set
		DB:       config.AppConfig.DB, // use default DB
		PoolSize: 100,                  // 连接池大小
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	RedisWriteClient = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.RedisWriteUrl,
		Password: "",                  // no password set
		DB:       config.AppConfig.RedisWriteDB, // use default DB
		PoolSize: 100,                  // 连接池大小
	})

	ctxy, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = RedisWriteClient.Ping(ctxy).Result()
	if err != nil {
		panic(err)
	}

	RedisSessionWriteClient = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.RedisSessionWriteUrl,
		Password: "",                  // no password set
		DB:       config.AppConfig.RedisSessionWriteDB, // use default DB
		PoolSize: 100,                  // 连接池大小
	})

	ctxyz, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = RedisSessionWriteClient.Ping(ctxyz).Result()
	if err != nil {
		panic(err)
	}


	RedisSessionClient = redis.NewClient(&redis.Options{
		Addr:     config.AppConfig.RedisSessionUrl,
		Password: "",                  // no password set
		DB:       config.AppConfig.RedisSessionDB, // use default DB
		PoolSize: 100,                  // 连接池大小
	})

	ctxyq, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = RedisSessionClient.Ping(ctxyq).Result()
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
	_, err = RedisWriteClient.Set(ctx, key, value, t).Result()
	return
}

func GetRedisSession(key string) (val string, err error) {
	ctx := context.Background()
	val, err = RedisSessionClient.Get(ctx, key).Result()
	return
}

func SetRedisSession(key string, value string, t time.Duration) (err error) {
	ctx := context.Background()
	_, err = RedisSessionWriteClient.Set(ctx, key, value, t).Result()
	return
}