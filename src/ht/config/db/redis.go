package db

import (
	"github.com/go-redis/redis/v8"
	"ht-crm/autoconfigure"
)

var Redis *redis.Client

func init() {
	redisConfig := autoconfigure.GetAppCig().Datasource.Redis
	Redis = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Username: redisConfig.User,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})
}
