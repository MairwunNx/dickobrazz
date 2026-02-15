package collector

import (
	"context"
	"dickobrazz/src/shared/metrics"
	"time"
)

const defaultStatsInterval = time.Minute

type StatsCollector struct {
	ctx       context.Context
	interval  time.Duration
	startTime time.Time
}

func NewStatsCollector(ctx context.Context, startTime time.Time) *StatsCollector {
	return &StatsCollector{
		ctx:       ctx,
		interval:  defaultStatsInterval,
		startTime: startTime,
	}
}

func (s *StatsCollector) Start() {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.collectStats()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.collectStats()
		}
	}
}

func (s *StatsCollector) collectStats() {
	if s.ctx.Err() != nil {
		return
	}

	metrics.SetUptimeSeconds(time.Since(s.startTime).Seconds())
}
