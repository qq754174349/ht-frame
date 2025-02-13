package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/qq754174349/ht-frame/config"
)

var Redis *redis.Client

type AutoConfig struct{}

func (AutoConfig) Init(cfg interface{}) error {
	redisConfig := cfg.(config.RedisConfig)
	Redis = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Username: redisConfig.User,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
	return nil
}
