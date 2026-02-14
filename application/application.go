package application

import (
	"context"
	"dickobrazz/application/api"
	"dickobrazz/application/collector"
	"dickobrazz/application/localization"
	"dickobrazz/application/logging"
	"dickobrazz/application/metrics"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Application struct {
	ctx            context.Context
	cancel         context.CancelFunc
	log            *logging.Logger
	bot            *tgbotapi.BotAPI
	localization   *localization.LocalizationManager
	api            *api.APIClient
	outsiders      *OutsiderServers
	statsCollector *collector.StatsCollector
	wg             sync.WaitGroup
	startTime      time.Time
}

func NewApplication() *Application {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	log := logging.NewLogger()
	cfg := LoadConfiguration(log)
	if err := metrics.Register(); err != nil {
		log.F("Failed to register metrics", logging.InnerError, err)
	}

	bot := InitializeTelegramBot(log, cfg)
	localizationManager, err := localization.NewLocalizationManager(log)
	if err != nil {
		log.F("Failed to initialize localization manager", logging.InnerError, err)
	}
	apiClient := api.NewAPIClient(cfg.Bot.Server.BaseURL, cfg.Bot.CSOT)
	startTime := time.Now()

	app := &Application{
		ctx:          ctx,
		cancel:       cancel,
		log:          log,
		bot:          bot,
		localization: localizationManager,
		api:          apiClient,
		startTime:    startTime,
	}
	app.outsiders = InitializeOutsiderServers(log, &app.wg)
	app.statsCollector = collector.NewStatsCollector(app.ctx, startTime)

	return app
}

func (app *Application) Shutdown() {
	app.cancel()

	if app.outsiders != nil {
		if err := app.outsiders.Shutdown(app.ctx); err != nil {
			app.log.E("Failed to shutdown outsider servers", logging.InnerError, err)
		}
	}

	app.wg.Wait()

	app.log.I("Gracefully shutting down... Bye!")
}

func (app *Application) Run() {
	if app.outsiders != nil {
		app.outsiders.Start()
	}
	if app.statsCollector != nil {
		app.wg.Add(1)
		go func() {
			defer app.wg.Done()
			app.statsCollector.Start()
		}()
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

			processingStarted := time.Now()
			handledKinds := 0
			updateKind := "ignored"
			if _, detectedLang := app.localization.LocalizerByUpdate(&update); detectedLang != "" {
				metrics.IncDetectedLanguage(detectedLang)
			}

			if msg := update.Message; msg != nil {
				user := update.SentFrom()
				log := app.log.With(
					logging.UserId, user.ID,
					logging.UserName, user.UserName,
					logging.ChatType, msg.Chat.Type,
					logging.ChatId, msg.Chat.ID,
				)
				log.I("Received message")
				metrics.IncMessagesHandled("message")
				handledKinds++
				updateKind = "message"

				if msg.IsCommand() {
					app.HandleCommand(log, &update)
				}
			}

			if query := update.InlineQuery; query != nil {
				user := update.SentFrom()
				log := app.log.With(
					logging.UserId, user.ID,
					logging.UserName, user.UserName,
					logging.QueryId, query.ID,
					logging.ChatType, query.ChatType,
				)

				app.HandleInlineQuery(log, &update)
				metrics.IncMessagesHandled("inline_query")
				handledKinds++
				updateKind = "inline_query"
			}

			if callback := update.CallbackQuery; callback != nil {
				user := update.SentFrom()
				log := app.log.With(
					logging.UserId, user.ID,
					logging.UserName, user.UserName,
					"callback_id", callback.ID,
					"callback_data", callback.Data,
				)

				app.HandleCallbackQuery(log, &update)
				metrics.IncMessagesHandled("callback_query")
				handledKinds++
				updateKind = "callback_query"
			}

			if handledKinds == 0 {
				metrics.IncMessagesIgnored("unsupported_update")
				updateKind = "ignored"
			} else if handledKinds > 1 {
				updateKind = "multiple"
			}
			metrics.ObserveUpdateDuration(updateKind, time.Since(processingStarted))
		}
	}
}
