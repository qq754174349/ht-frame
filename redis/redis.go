package redis

import (
	"github.com/go-redis/redis/v8"
	"github.com/qq754174349/ht-frame/config"
)

var Redis *redis.Client

type AutoConfig struct{}

func (AutoConfig) Init(cfg *config.AppConfig) error {
	redisConfig := cfg.Datasource.Redis
	Redis = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Username: redisConfig.User,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
	return nil
}
