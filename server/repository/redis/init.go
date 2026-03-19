package redis

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/sw5005-sus/ceramicraft-user-mservice/server/config"
)

var redisClient *redis.Client

func InitRedis() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.Config.RedisConfig.Host, config.Config.RedisConfig.Port),
		DB:   0, // Use default DB
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}
