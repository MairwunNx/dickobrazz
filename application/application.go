package application

import (
	"context"
	"dickobrazz/application/logging"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/go-redis/cache/v9"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	ctx         context.Context
	cancel      context.CancelFunc
	log         *logging.Logger
	bot         *tgbotapi.BotAPI
	rnd         *Random
	db          *mongo.Client
	redis       *redis.Client
	cache       *cache.Cache
	healthcheck *HealthcheckServer
	wg          sync.WaitGroup
	startTime   time.Time
}

func NewApplication() *Application {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	log := logging.NewLogger()

	bot := InitializeTelegramBot(log)
	rnd := InitializeRandom(log)
	db := InitializeMongoConnection(ctx, log)
	client, redisCache := InitializeRedisConnection(log)

	app := &Application{
		ctx:         ctx,
		cancel:      cancel,
		log:         log,
		bot:         bot,
		rnd:         rnd,
		db:          db,
		redis:       client,
		cache:       redisCache,
		startTime:   time.Now(),
	}
	app.healthcheck = InitializeHealthcheckServer(log, &app.wg)

	return app
}

func (app *Application) Shutdown() {
	app.cancel()

	if app.healthcheck != nil {
		if err := app.healthcheck.Shutdown(app.ctx); err != nil {
			app.log.E("Failed to shutdown healthcheck server", logging.InnerError, err)
		}
	}

	app.wg.Wait()

	if err := app.db.Disconnect(app.ctx); err != nil {
		app.log.E("Failed to disconnect MongoDB", logging.InnerError, err)
	}

	if err := app.redis.Close(); err != nil {
		app.log.E("Failed to close Redis connection", logging.InnerError, err)
	}

	app.log.I("Gracefully shutting down... Bye!")
}

func (app *Application) Run() {
	if app.healthcheck != nil {
		app.healthcheck.Start()
	}

	updates := tgbotapi.NewUpdate(0)
	updates.Timeout = 60

	updatesChan := app.bot.GetUpdatesChan(updates)

	app.log.I("Bot started, waiting for updates...")

	for {
		select {
		case <-app.ctx.Done():
			app.log.I("Received shutdown signal, stopping bot...")
			app.bot.StopReceivingUpdates()
			return

		case update, ok := <-updatesChan:
			if !ok {
				app.log.I("Updates channel closed, stopping bot...")
				return
			}

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

			if callback := update.CallbackQuery; callback != nil {
				user := update.SentFrom()
				log := app.log.With(
					logging.UserId, user.ID,
					logging.UserName, user.UserName,
					"callback_id", callback.ID,
					"callback_data", callback.Data,
				)

				app.HandleCallbackQuery(log, callback)
			}
		}
	}
}
