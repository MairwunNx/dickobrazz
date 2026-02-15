package app

import (
	"context"
	"dickobrazz/src/shared/logging"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Poller struct {
	bot    *tgbotapi.BotAPI
	router *Router
	log    *logging.Logger
}

func NewPoller(bot *tgbotapi.BotAPI, router *Router, log *logging.Logger) *Poller {
	return &Poller{bot: bot, router: router, log: log}
}

func (p *Poller) Start(ctx context.Context) {
	updates := tgbotapi.NewUpdate(0)
	updates.Timeout = 60

	updatesChan := p.bot.GetUpdatesChan(updates)

	p.log.I("Bot started, waiting for updates...")

	for {
		select {
		case <-ctx.Done():
			p.log.I("Received shutdown signal, stopping bot...")
			p.bot.StopReceivingUpdates()
			return

		case update, ok := <-updatesChan:
			if !ok {
				p.log.I("Updates channel closed, stopping bot...")
				return
			}

			p.router.HandleUpdate(update)
		}
	}
}
