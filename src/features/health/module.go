package health

import (
	"context"
	"dickobrazz/src/shared/logging"
	"net/http"
	"time"

	"go.uber.org/fx"
)

const healthAddr = ":8080"

var Module = fx.Module("health",
	fx.Invoke(func(lc fx.Lifecycle, log *logging.Logger) {
		mux := http.NewServeMux()
		mux.HandleFunc("/health", healthCheckHandler)
		server := &http.Server{
			Addr:         healthAddr,
			Handler:      mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  120 * time.Second,
		}

		lc.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					log.I("HTTP server is starting", "kind", "healthcheck", "addr", healthAddr)
					if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
						log.E("HTTP server failed", "kind", "healthcheck", logging.InnerError, err)
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
