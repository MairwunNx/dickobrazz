package seasons

import "go.uber.org/fx"

var Module = fx.Module("seasons",
	fx.Provide(
		NewGetAction,
		NewHandler,
		NewCallbackHandler,
	),
)
