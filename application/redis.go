package application

import (
	"context"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

type UserCock struct {
	userId   int64  `json:"user_id,omitempty"`
	userName string `json:"user_name,omitempty"`
	size     int    `json:"size,omitempty"`
}

func InitializeRedisConnection(log *Logger) (*redis.Client, *cache.Cache) {
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		log.F("REDIS_PORT does not have value, set it in .env file")
	}

	address := "redis:" + port

	client := redis.NewClient(&redis.Options{Addr: address, Password: "", DB: 0})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.F("Failed to connect to Redis", InnerError, err)
	}

	redisCache := cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(100, time.Minute),
	})

	log.I("Successfully connected to Redis!")
	return client, redisCache
}
