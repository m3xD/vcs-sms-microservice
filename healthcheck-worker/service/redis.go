package service

import (
	"context"
	"healthcheck-worker/util"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	RedisClient *redis.Client
}

func NewRedis() *Redis {
	return &Redis{RedisClient: util.GetInstance()}
}

func (redis *Redis) SetServer(key string, value interface{}, duration time.Duration) error {
	err := redis.RedisClient.Set(context.Background(), key, value, duration).Err()
	return err
}

func (redis *Redis) GetServer(key string) (string, error) {
	return redis.RedisClient.Get(context.Background(), key).Result()
}
