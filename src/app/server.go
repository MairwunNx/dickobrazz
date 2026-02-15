package app

import (
	"context"
	"dickobrazz/src/shared/logging"
	"dickobrazz/src/shared/metrics"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	healthAddr        = ":8080"
	systemMetricsAddr = ":9100"
	appMetricsAddr    = ":9091"
)

type OutsiderServers struct {
	health        *http.Server
	systemMetrics *http.Server
	appMetrics    *http.Server
	log           *logging.Logger
	wg            *sync.WaitGroup
}

func InitializeOutsiderServers(log *logging.Logger, wg *sync.WaitGroup) *OutsiderServers {
	healthMux := http.NewServeMux()
	healthMux.HandleFunc("/health", healthCheckHandler)
	healthServer := &http.Server{
		Addr:         healthAddr,
		Handler:      healthMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	systemRegistry := prometheus.NewRegistry()
	systemRegistry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewBuildInfoCollector(),
	)
	systemMux := http.NewServeMux()
	systemMux.Handle("/metrics", promhttp.HandlerFor(systemRegistry, promhttp.HandlerOpts{}))
	systemServer := &http.Server{
		Addr:         systemMetricsAddr,
		Handler:      systemMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	appMux := http.NewServeMux()
	appMux.Handle("/metrics", promhttp.HandlerFor(metrics.Registry(), promhttp.HandlerOpts{}))
	appServer := &http.Server{
		Addr:         appMetricsAddr,
		Handler:      appMux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	return &OutsiderServers{
		health:        healthServer,
		systemMetrics: systemServer,
		appMetrics:    appServer,
		log:           log,
		wg:            wg,
	}
}

func (os *OutsiderServers) Start() {
	os.startServer("healthcheck", os.health)
	os.startServer("system_metrics", os.systemMetrics)
	os.startServer("application_metrics", os.appMetrics)
}

func (os *OutsiderServers) Shutdown(ctx context.Context) error {
	os.log.I("Shutting down outsider servers")
	var firstErr error
	for _, server := range []*http.Server{os.health, os.systemMetrics, os.appMetrics} {
		if err := server.Shutdown(ctx); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

func (os *OutsiderServers) startServer(kind string, server *http.Server) {
	os.wg.Add(1)
	go func() {
		defer os.wg.Done()
		os.log.I("HTTP server is starting", "kind", kind, "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			os.log.E("HTTP server failed", "kind", kind, logging.InnerError, err)
		}
	}()
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","service":"dickobrazz"}`))
}
