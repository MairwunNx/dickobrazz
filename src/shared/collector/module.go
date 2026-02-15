package collector

import (
	"context"
	"time"

	"go.uber.org/fx"
)

var Module = fx.Module("collector",
	fx.Provide(func(lc fx.Lifecycle) *StatsCollector {
		sc := NewStatsCollector(context.Background(), time.Now())
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go sc.Start()
				return nil
			},
		})
		return sc
	}),
)
