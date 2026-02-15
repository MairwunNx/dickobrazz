package privacy

import "go.uber.org/fx"

var Module = fx.Module("privacy",
	fx.Provide(NewHandler),
)
