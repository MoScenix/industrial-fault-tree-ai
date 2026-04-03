package redis

import (
	"context"
	"fmt"
	"os"

	"github.com/MoScenix/industrial-fault-tree-ai/app/ai/conf"
	"github.com/redis/go-redis/v9"
)

var (
	RedisClient *redis.Client
)

func Init() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     conf.GetConf().Redis.Address,
		Username: conf.GetConf().Redis.Username,
		Password: fmt.Sprintf(conf.GetConf().Redis.Password, os.Getenv("REDIS_PASSWORD")),
		DB:       conf.GetConf().Redis.DB,
	})
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
}
