package size

import "go.uber.org/fx"

var Module = fx.Module("size",
	fx.Provide(
		NewGenerateAction,
		NewHandler,
	),
)
