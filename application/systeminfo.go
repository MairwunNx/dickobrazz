package application

import (
	"context"
	"dickobrazz/application/localization"
	"dickobrazz/application/logging"
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.mongodb.org/mongo-driver/bson"
)

type SystemInfo struct {
	Uptime   string
	Version  string
	BuildRev string
	BuildAt  string

	OS            string
	Arch          string
	GoVersion     string
	MemoryUsed    uint64
	MemoryTotal   uint64
	MemoryPercent float64

	MongoVersion string
	RedisVersion string

	UserID   int64
	Username string
	BotID    int64
}

func (app *Application) GetSystemInfo(log *logging.Logger, localizer *i18n.Localizer, userID int64, username string) *SystemInfo {
	info := &SystemInfo{
		UserID:   userID,
		Username: username,
	}

	info.Uptime = FormatUptime(app.localization, localizer, time.Since(app.startTime))

	info.Version = logging.Version
	info.BuildRev = logging.BuildRv
	info.BuildAt = logging.BuildAt

	info.OS = runtime.GOOS
	info.Arch = runtime.GOARCH
	info.GoVersion = logging.GoVersion

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	info.MemoryUsed = m.Alloc / 1024 / 1024 // МБ
	info.MemoryTotal = m.Sys / 1024 / 1024  // МБ
	if info.MemoryTotal > 0 {
		info.MemoryPercent = float64(info.MemoryUsed) / float64(info.MemoryTotal) * 100
	}

	info.MongoVersion = app.GetMongoVersion(log, localizer)
	info.RedisVersion = app.GetRedisVersion(log, localizer)

	info.BotID = app.bot.Self.ID
	return info
}

func (app *Application) GetMongoVersion(log *logging.Logger, localizer *i18n.Localizer) string {
	ctx, cancel := context.WithTimeout(app.ctx, 5*time.Second)
	defer cancel()

	var result bson.M
	err := app.db.Database("admin").RunCommand(ctx, bson.D{{Key: "buildInfo", Value: 1}}).Decode(&result)
	if err != nil {
		log.E("Failed to get MongoDB version", logging.InnerError, err)
		return app.localization.Localize(localizer, MsgUnknownValue, nil)
	}

	if version, ok := result["version"].(string); ok {
		return version
	}

	return app.localization.Localize(localizer, MsgUnknownValue, nil)
}

func (app *Application) GetRedisVersion(log *logging.Logger, localizer *i18n.Localizer) string {
	ctx, cancel := context.WithTimeout(app.ctx, 5*time.Second)
	defer cancel()

	info, err := app.redis.Info(ctx, "server").Result()
	if err != nil {
		log.E("Failed to get Redis version", logging.InnerError, err)
		return app.localization.Localize(localizer, MsgUnknownValue, nil)
	}

	// Парсим info строку для получения версии
	lines := strings.Split(info, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "redis_version:") {
			return strings.TrimPrefix(line, "redis_version:")
		}
	}

	return app.localization.Localize(localizer, MsgUnknownValue, nil)
}

func FormatUptime(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	minutes := int(d.Minutes()) % 60

	if days > 0 {
		dayStr := localizationManager.Localize(localizer, UptimeDayShort, map[string]any{"Count": days})
		hourStr := localizationManager.Localize(localizer, UptimeHourShort, map[string]any{"Count": hours})
		minStr := localizationManager.Localize(localizer, UptimeMinuteShort, map[string]any{"Count": minutes})
		return fmt.Sprintf("%s %s %s", dayStr, hourStr, minStr)
	} else if hours > 0 {
		hourStr := localizationManager.Localize(localizer, UptimeHourShort, map[string]any{"Count": hours})
		minStr := localizationManager.Localize(localizer, UptimeMinuteShort, map[string]any{"Count": minutes})
		return fmt.Sprintf("%s %s", hourStr, minStr)
	} else {
		return localizationManager.Localize(localizer, UptimeMinuteShort, map[string]any{"Count": minutes})
	}
}
