package shared

import (
	"dickobrazz/src/shared/api"
	"dickobrazz/src/shared/config"
	"dickobrazz/src/shared/localization"
	"dickobrazz/src/shared/logging"

	"go.uber.org/fx"
)

var Module = fx.Module("shared",
	fx.Provide(
		logging.NewLogger,
		config.NewConfiguration,
		localization.NewLocalizationManager,
		func(cfg *config.Configuration) *api.APIClient {
			return api.NewAPIClient(cfg.Bot.Server.BaseURL, cfg.Bot.CSOT)
		},
	),
)
