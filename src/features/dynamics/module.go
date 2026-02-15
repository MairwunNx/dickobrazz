package dynamics

import "go.uber.org/fx"

var Module = fx.Module("dynamics",
	fx.Provide(
		NewGetAction,
		NewHandler,
	),
)
