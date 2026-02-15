package help

import "go.uber.org/fx"

var Module = fx.Module("help",
	fx.Provide(NewHandler),
)
