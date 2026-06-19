package config

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func ConnectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("redis_host"),
		Password: "",
		DB:       0,
		Protocol: 2,
	})
	return client

}
