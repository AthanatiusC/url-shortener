package db

import (
	"context"
	"fmt"
	"time"

	"github.com/AthanatiusC/url-shortener/config"
	"github.com/AthanatiusC/url-shortener/helper/logger"
	"github.com/redis/go-redis/v9"
)

func InitRedis(ctx context.Context, config *config.RedisConfig) *redis.Client {
	ctx = context.WithValue(ctx, "source", "redis")
	logger.Info("Connecting to Redis")
	redis := redis.NewClient(&redis.Options{
		Addr:        fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password:    config.Password,
		DB:          config.DB,
		DialTimeout: time.Duration(config.TimeOut),
	})
	err := redis.Ping(ctx).Err()
	if err != nil {
		panic(err)
	}
	logger.Info("Connected to Redis")
	return redis
}
