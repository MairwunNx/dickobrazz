package application

import (
	"context"
	"dickobrazz/application/logging"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-redis/cache/v9"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	ctx    context.Context
	cancel context.CancelFunc
	log    *logging.Logger
	bot    *tgbotapi.BotAPI
	rnd    *Random
	db     *mongo.Client
	redis  *redis.Client
	cache  *cache.Cache
}

func NewApplication() *Application {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	log := logging.NewLogger()

	bot := InitializeTelegramBot(log)
	rnd := InitializeRandom(log)
	db := InitializeMongoConnection(ctx, log)
	client, redisCache := InitializeRedisConnection(log)

	return &Application{ctx, cancel, log, bot, rnd, db, client, redisCache}
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
	updates.Timeout = 60

	for update := range app.bot.GetUpdatesChan(updates) {
		if msg := update.Message; msg != nil {
			user := update.SentFrom()
			app.log.With(
				logging.UserId, user.ID,
				logging.UserName, user.UserName,
				logging.ChatType, msg.Chat.Type,
				logging.ChatId, msg.Chat.ID,
			).I("Received message")
		}

		if query := update.InlineQuery; query != nil {
			user := update.SentFrom()
			log := app.log.With(
				logging.UserId, user.ID,
				logging.UserName, user.UserName,
				logging.QueryId, query.ID,
				logging.ChatType, query.ChatType,
			)

			app.HandleInlineQuery(log, query)
		}
	}
}
