package app

import "go.uber.org/fx"

var Module = fx.Module("app",
	fx.Provide(
		InitializeTelegramBot,
		NewRouter,
		NewApplication,
	),
)
