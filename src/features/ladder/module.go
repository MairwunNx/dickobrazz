package ladder

import "go.uber.org/fx"

var Module = fx.Module("ladder",
	fx.Provide(
		NewGetAction,
		NewHandler,
	),
)
