package metrics

import (
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

const metricsNamespace = "dickobrazz"

var (
	appRegistry = prometheus.NewRegistry()
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
	totalUsers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: metricsNamespace,
			Name:      "total_users",
			Help:      "Total number of unique users",
		},
	)
	statsDAU = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: metricsNamespace,
			Name:      "stats_dau",
			Help:      "Daily Active Users (last 24h)",
		},
	)
	statsMAU = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: metricsNamespace,
			Name:      "stats_mau",
			Help:      "Monthly Active Users (last 30d)",
		},
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
	sizeDistribution = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: metricsNamespace,
			Name:      "size_distribution",
			Help:      "Distribution of sizes",
		},
		[]string{"bucket"},
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
			totalUsers,
			statsDAU,
			statsMAU,
			detectedLanguages,
			uptimeSeconds,
			availabilityPercent,
			sizeDistribution,
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

func SetTotalUsers(value float64) {
	totalUsers.Set(value)
}

func SetDAU(value float64) {
	statsDAU.Set(value)
}

func SetMAU(value float64) {
	statsMAU.Set(value)
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

func SetSizeDistribution(bucket string, value float64) {
	sizeDistribution.WithLabelValues(bucket).Set(value)
}

func Registry() *prometheus.Registry {
	return appRegistry
}
