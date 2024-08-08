package util

import (
	"os"

	"github.com/redis/go-redis/v9"
)

func GetInstance() *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	return client
}
