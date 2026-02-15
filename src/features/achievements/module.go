package achievements

import "go.uber.org/fx"

var Module = fx.Module("achievements",
	fx.Provide(
		NewGetAction,
		NewHandler,
		NewCallbackHandler,
	),
)
