package db

import (
	"github.com/go-redis/redis/v8"
	"ht-crm/src/ht/config"
)

var Redis *redis.Client

func init() {
	redisConfig := config.GetEnvCfg().Redis
	Redis = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Username: redisConfig.User,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
}
