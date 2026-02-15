package collector

import (
	"context"
	"time"

	"go.uber.org/fx"
)

var Module = fx.Module("collector",
	fx.Invoke(func(lc fx.Lifecycle) {
		ctx, cancel := context.WithCancel(context.Background())
		sc := NewStatsCollector(ctx, time.Now())
		lc.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				go sc.Start()
				return nil
			},
			OnStop: func(_ context.Context) error {
				cancel()
				return nil
			},
		})
	}),
)
