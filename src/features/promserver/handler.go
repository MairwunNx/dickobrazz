package promserver

import (
	"dickobrazz/src/shared/metrics"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func newMetricsHandler() http.Handler {
	registry := prometheus.NewRegistry()

	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		collectors.NewBuildInfoCollector(),
	)

	return promhttp.HandlerFor(
		prometheus.Gatherers{registry, metrics.Registry()},
		promhttp.HandlerOpts{},
	)
}
