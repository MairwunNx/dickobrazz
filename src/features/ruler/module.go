package ruler

import "go.uber.org/fx"

var Module = fx.Module("ruler",
	fx.Provide(
		NewGetAction,
		NewHandler,
	),
)
