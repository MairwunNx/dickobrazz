package api

import (
	"dickobrazz/src/shared/config"

	"go.uber.org/fx"
)

var Module = fx.Module("api",
	fx.Provide(func(cfg *config.Configuration) *APIClient {
		return NewAPIClient(cfg.Bot.Server.BaseURL, cfg.Bot.CSOT)
	}),
)
