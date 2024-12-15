package application

import (
	"context"
	"github.com/go-redis/cache/v9"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"os/signal"
	"syscall"
)

var Version = "unknown" // 1
var BuildAt = "unknown" // 2024-06-29_15:48:20
var BuildRv = "unknown" // 5b329d0

type Application struct {
	ctx    context.Context
	cancel context.CancelFunc
	log    *Logger
	bot    *tgbotapi.BotAPI
	db     *mongo.Client
	redis  *redis.Client
	cache  *cache.Cache
}

func NewApplication() *Application {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	log := NewLogger()

	bot := InitializeTelegramBot(log)
	db := InitializeMongoConnection(ctx, log)
	client, redisCache := InitializeRedisConnection(log)

	return &Application{ctx, cancel, log, bot, db, client, redisCache}
}

func (app *Application) Shutdown() {
	app.cancel()
	if err := app.db.Disconnect(app.ctx); err != nil {
		app.log.E("Failed to disconnect MongoDB", err)
	}
	app.log.I("Gracefully shutting down... Bye!")
}

func (app *Application) Run() {
	updates := tgbotapi.NewUpdate(0)
	updates.Timeout = 120

	for update := range app.bot.GetUpdatesChan(updates) {
		if query := update.InlineQuery; query != nil {
			user := update.SentFrom()
			log := app.log.With(
				UserId, user.ID,
				UserName, user.UserName,
				QueryId, query.ID,
				ChatType, query.ChatType,
			)

			app.HandleInlineQuery(log, query)
		}
	}
}
