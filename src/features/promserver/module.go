package promserver

import (
	"context"
	"dickobrazz/src/shared/logging"
	"net/http"
	"time"

	"go.uber.org/fx"
)

const metricsAddr = ":9091"

var Module = fx.Module("promserver",
	fx.Invoke(func(lc fx.Lifecycle, log *logging.Logger) {
		mux := http.NewServeMux()
		mux.Handle("/metrics", newMetricsHandler())
		server := &http.Server{
			Addr:         metricsAddr,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}

		lc.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					log.I("HTTP server is starting", "kind", "metrics", "addr", metricsAddr)
					if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						log.E("HTTP server failed", "kind", "metrics", logging.InnerError, err)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return server.Shutdown(ctx)
			},
		})
	}),
)
