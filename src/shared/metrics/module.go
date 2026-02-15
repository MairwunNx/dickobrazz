package metrics

import (
	"context"

	"go.uber.org/fx"
)

var Module = fx.Module("metrics",
	fx.Invoke(func(lc fx.Lifecycle) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return Register()
			},
		})
	}),
)
