package storage

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
)

func ConnectRedis(ctx context.Context, log *slog.Logger, address, password string) (*redis.Client, *cache.Cache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       0,
	})

	ctxPing, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if _, err := client.Ping(ctxPing).Result(); err != nil {
		return nil, nil, err
	}

	redisCache := cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(100, time.Minute),
	})

	log.Info("Redis connection established")
	return client, redisCache, nil
}
