package app

import (
	"context"
	"dickobrazz/src/shared/logging"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/fx"
)

type Application struct {
	ctx    context.Context
	cancel context.CancelFunc
	log    *logging.Logger
	poller *Poller
}

func NewApplication(lc fx.Lifecycle, log *logging.Logger, poller *Poller) *Application {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	app := &Application{
		ctx:    ctx,
		cancel: cancel,
		log:    log,
		poller: poller,
	}

	return app
}

func (app *Application) Shutdown() {
	app.cancel()
	app.log.I("Gracefully shutting down... Bye!")
}

func (app *Application) Run() {
	app.poller.Start(app.ctx)
}
