package metrics

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const metricsNamespace = "dickobrazz"

var (
	appRegistry     = prometheus.NewRegistry()
	messagesHandled = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "messages_handled_total",
			Help:      "Total number of handled updates",
		},
		[]string{"status"},
	)
	messagesIgnored = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "messages_ignored_total",
			Help:      "Total number of ignored updates",
		},
		[]string{"reason"},
	)
	updateDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: metricsNamespace,
			Name:      "update_duration_seconds",
			Help:      "Duration of update handling in seconds",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"type"},
	)
	detectedLanguages = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Name:      "detected_languages_total",
			Help:      "Total number of detected languages",
		},
		[]string{"language"},
	)
	uptimeSeconds = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: metricsNamespace,
			Name:      "uptime_seconds",
			Help:      "Service uptime in seconds",
		},
	)
	availabilityPercent = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: metricsNamespace,
			Name:      "availability_percent",
			Help:      "Service availability in percent",
		},
	)
	registerOnce sync.Once
)

func Register() error {
	var registerErr error
	registerOnce.Do(func() {
		collectors := []prometheus.Collector{
			messagesHandled,
			messagesIgnored,
			updateDuration,
			detectedLanguages,
			uptimeSeconds,
			availabilityPercent,
		}
		for _, collector := range collectors {
			if err := appRegistry.Register(collector); err != nil {
				if _, ok := err.(prometheus.AlreadyRegisteredError); ok {
					continue
				}
				registerErr = err
				return
			}
		}
	})
	return registerErr
}

func IncMessagesHandled(status string) {
	messagesHandled.WithLabelValues(status).Inc()
}

func IncMessagesIgnored(reason string) {
	messagesIgnored.WithLabelValues(reason).Inc()
}

func ObserveUpdateDuration(updateType string, duration time.Duration) {
	updateDuration.WithLabelValues(updateType).Observe(duration.Seconds())
}

func IncDetectedLanguage(language string) {
	detectedLanguages.WithLabelValues(language).Inc()
}

func SetUptimeSeconds(value float64) {
	uptimeSeconds.Set(value)
}

func SetAvailabilityPercent(value float64) {
	availabilityPercent.Set(value)
}

func Registry() *prometheus.Registry {
	return appRegistry
}
