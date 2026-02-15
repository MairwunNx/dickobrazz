package app

import (
	"context"
	"dickobrazz/src/shared/collector"
	"dickobrazz/src/shared/logging"
	"dickobrazz/src/shared/metrics"
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
	outsiders      *OutsiderServers
	statsCollector *collector.StatsCollector
	router         *Router
	wg             sync.WaitGroup
}

func NewApplication(log *logging.Logger, bot *tgbotapi.BotAPI, router *Router) *Application {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	if err := metrics.Register(); err != nil {
		log.F("Failed to register metrics", logging.InnerError, err)
	}

	startTime := time.Now()

	app := &Application{
		ctx:    ctx,
		cancel: cancel,
		log:    log,
		bot:    bot,
		router: router,
	}
	app.outsiders = InitializeOutsiderServers(log, &app.wg)
	app.statsCollector = collector.NewStatsCollector(app.ctx, startTime)

	// Set context on the router
	app.router.ctx = ctx

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

			app.router.HandleUpdate(update)
		}
	}
}
