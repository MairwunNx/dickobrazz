package application

import (
	"context"
	"dickobrazz/application/logging"
	"fmt"
	"runtime"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type SystemInfo struct {
	Uptime       string
	Version      string
	BuildRev     string
	BuildAt      string
	
	OS           string
	Arch         string
	GoVersion    string
	MemoryUsed   uint64
	MemoryTotal  uint64
	MemoryPercent float64
	
	MongoVersion string
	RedisVersion string
	
	UserID       int64
	Username     string
	BotID        int64
}

func (app *Application) GetSystemInfo(log *logging.Logger, userID int64, username string) *SystemInfo {
	info := &SystemInfo{
		UserID:   userID,
		Username: username,
	}
	
	info.Uptime = FormatUptime(time.Since(app.startTime))
	
	info.Version = logging.Version
	info.BuildRev = logging.BuildRv
	info.BuildAt = logging.BuildAt
	
	info.OS = runtime.GOOS
	info.Arch = runtime.GOARCH
	info.GoVersion = logging.GoVersion
	
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	info.MemoryUsed = m.Alloc / 1024 / 1024 // МБ
	info.MemoryTotal = m.Sys / 1024 / 1024   // МБ
	if info.MemoryTotal > 0 {
		info.MemoryPercent = float64(info.MemoryUsed) / float64(info.MemoryTotal) * 100
	}
	
	info.MongoVersion = app.GetMongoVersion(log)
	info.RedisVersion = app.GetRedisVersion(log)
	
	info.BotID = app.bot.Self.ID
	return info
}

func (app *Application) GetMongoVersion(log *logging.Logger) string {
	var result bson.M
	err := app.db.Database("admin").RunCommand(context.Background(), bson.D{{Key: "buildInfo", Value: 1}}).Decode(&result)
	if err != nil {
		log.E("Failed to get MongoDB version", logging.InnerError, err)
		return "неизвестно"
	}
	
	if version, ok := result["version"].(string); ok {
		return version
	}
	
	return "неизвестно"
}

func (app *Application) GetRedisVersion(log *logging.Logger) string {
	info, err := app.redis.Info(context.Background(), "server").Result()
	if err != nil {
		log.E("Failed to get Redis version", logging.InnerError, err)
		return "неизвестно"
	}
	
	// Парсим info строку для получения версии
	lines := strings.Split(info, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "redis_version:") {
			return strings.TrimPrefix(line, "redis_version:")
		}
	}
	
	return "неизвестно"
}

func FormatUptime(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60
	
	if days > 0 {
		return fmt.Sprintf("%dд %dч %dм", days, hours, minutes)
	} else if hours > 0 {
		return fmt.Sprintf("%dч %dм", hours, minutes)
	} else {
		return fmt.Sprintf("%dм", minutes)
	}
}
