package application

import (
	"context"
	"fmt"
	"github.com/go-redis/cache/v9"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"strings"
	"time"
)

type UserCock struct {
	UserId   int64  `json:"user_id,omitempty"`
	UserName string `json:"user_name,omitempty"`
	Size     int    `json:"size,omitempty"`
}

func InitializeRedisConnection(log *Logger) (*redis.Client, *cache.Cache) {
	port := os.Getenv("REDIS_PORT")
	if port == "" {
		log.F("REDIS_PORT does not have value, set it in .env file")
	}

	password := os.Getenv("REDIS_PASSWORD")
	if password == "" {
		log.F("REDIS_PASSWORD does not have value, set it in .env file")
	}

	address := "cache:" + port

	client := redis.NewClient(&redis.Options{Addr: address, Password: password, DB: 0})

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		log.F("Failed to connect to Redis", InnerError, err)
	}

	redisCache := cache.New(&cache.Options{
		Redis:      client,
		LocalCache: cache.NewTinyLFU(100, time.Minute),
	})

	log.I("Successfully connected to Redis!")
	return client, redisCache
}

func (app *Application) GetCockSizeFromCache(log *Logger, userID int64) *int {
	key := GetCockCacheKey(userID)

	var cock UserCock
	if err := app.cache.Get(app.ctx, key, &cock); err != nil {
		return nil
	}

	log.I("Successfully fetched cock from redis")
	return &cock.Size
}

func (app *Application) SaveCockToCache(log *Logger, userID int64, userName string, size int) {
	now := NowTime()
	midnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, NowLocation())
	ttl := time.Until(midnight)

	if err := app.cache.Set(&cache.Item{Ctx: app.ctx, Key: GetCockCacheKey(userID), Value: &UserCock{UserId: userID, UserName: userName, Size: size}, TTL: ttl}); err != nil {
		log.E("Failed to save cock Size to Redis", InnerError, err)
	} else {
		log.I("Successfully saved cock to redis")
	}
}

func GetCockCacheKey(userID int64) string {
	return fmt.Sprintf("cock_size:%d", userID)
}

func (app *Application) GetCockSizesFromCache(log *Logger) []UserCock {
	var cockSizes []UserCock

	iter := app.redis.Scan(app.ctx, 0, "cock_size:*", 0).Iterator()
	for iter.Next(app.ctx) {
		key := iter.Val()
		var cock UserCock

		err := app.cache.Get(app.ctx, key, &cock)
		if err != nil {
			return nil
		}

		userID, _ := strconv.ParseInt(strings.TrimPrefix(key, "cock_size:"), 10, 64)
		cock.UserId = userID
		cockSizes = append(cockSizes, cock)
	}

	if err := iter.Err(); err != nil {
		log.E("Failed to iterate over cock keys in Redis", InnerError, err)
		panic(err)
	}

	log.I("Successfully fetched all cock sizes from Redis")
	return cockSizes
}
