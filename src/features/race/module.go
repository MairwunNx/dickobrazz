package race

import "go.uber.org/fx"

var Module = fx.Module("race",
	fx.Provide(
		NewGetAction,
		NewHandler,
	),
)
