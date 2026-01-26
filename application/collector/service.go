package collector

import (
	"context"
	"dickobrazz/application/database"
	"dickobrazz/application/datetime"
	"dickobrazz/application/logging"
	"dickobrazz/application/metrics"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

const defaultStatsInterval = time.Minute

type StatsCollector struct {
	ctx      context.Context
	log      *logging.Logger
	db       *mongo.Client
	redis    *redis.Client
	interval time.Duration
	startTime time.Time
}

func NewStatsCollector(ctx context.Context, log *logging.Logger, db *mongo.Client, redis *redis.Client, startTime time.Time) *StatsCollector {
	return &StatsCollector{
		ctx:      ctx,
		log:      log,
		db:       db,
		redis:    redis,
		interval: defaultStatsInterval,
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

	if count, err := s.getTotalUsersCount(); err == nil {
		metrics.SetTotalUsers(float64(count))
	} else {
		s.log.E("Failed to collect total users stats", logging.InnerError, err)
	}

	now := datetime.NowTime()
	metrics.SetUptimeSeconds(time.Since(s.startTime).Seconds())

	if count, err := s.getActiveUsersCount(now.Add(-24 * time.Hour)); err == nil {
		metrics.SetDAU(float64(count))
	} else {
		s.log.E("Failed to collect DAU stats", logging.InnerError, err)
	}

	if count, err := s.getActiveUsersCount(now.Add(-30 * 24 * time.Hour)); err == nil {
		metrics.SetMAU(float64(count))
	} else {
		s.log.E("Failed to collect MAU stats", logging.InnerError, err)
	}

	s.collectSizeDistribution()
	s.updateAvailability(now)
}

func (s *StatsCollector) collectSizeDistribution() {
	collection := database.CollectionCocks(s.db)
	cursor, err := collection.Aggregate(s.ctx, database.PipelineSizeDistribution())
	if err != nil {
		s.log.E("Failed to collect size distribution stats", logging.InnerError, err)
		return
	}
	defer cursor.Close(s.ctx)

	type bucketRow struct {
		ID    any `bson:"_id"`
		Count int `bson:"count"`
	}
	rows := make([]bucketRow, 0)
	if err := cursor.All(s.ctx, &rows); err != nil {
		s.log.E("Failed to decode size distribution stats", logging.InnerError, err)
		return
	}

	boundaries := []int{0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55, 60, 62}
	counts := make(map[int]float64, len(boundaries))
	for _, row := range rows {
		idInt, ok := row.ID.(int32)
		if !ok {
			if id64, ok := row.ID.(int64); ok {
				idInt = int32(id64)
			} else {
				continue
			}
		}
		counts[int(idInt)] = float64(row.Count)
	}

	for i := 0; i < len(boundaries)-1; i++ {
		start := boundaries[i]
		end := boundaries[i+1] - 1
		label := strconv.Itoa(start) + "-" + strconv.Itoa(end)
		metrics.SetSizeDistribution(label, counts[start])
	}
}

func (s *StatsCollector) updateAvailability(now time.Time) {
	if s.redis == nil {
		return
	}

	firstSeenKey := "metrics:availability:first_seen"
	lastSeenKey := "metrics:availability:last_seen"
	downSecondsKey := "metrics:availability:down_seconds"

	firstSeen, err := s.getOrSetUnix(firstSeenKey, now.Unix())
	if err != nil {
		s.log.E("Failed to read availability first_seen", logging.InnerError, err)
		return
	}

	lastSeenRaw, err := s.redis.Get(s.ctx, lastSeenKey).Result()
	if err != nil && err != redis.Nil {
		s.log.E("Failed to read availability last_seen", logging.InnerError, err)
		return
	}

	downSeconds, err := s.getFloat64(downSecondsKey)
	if err != nil {
		s.log.E("Failed to read availability down_seconds", logging.InnerError, err)
		return
	}

	if lastSeenRaw != "" {
		lastSeenUnix, err := strconv.ParseInt(lastSeenRaw, 10, 64)
		if err == nil {
			delta := now.Sub(time.Unix(lastSeenUnix, 0))
			if delta > s.interval*2 {
				downSeconds += delta.Seconds() - s.interval.Seconds()
			}
		}
	}

	if err := s.redis.Set(s.ctx, lastSeenKey, now.Unix(), 0).Err(); err != nil {
		s.log.E("Failed to write availability last_seen", logging.InnerError, err)
		return
	}
	if err := s.redis.Set(s.ctx, downSecondsKey, downSeconds, 0).Err(); err != nil {
		s.log.E("Failed to write availability down_seconds", logging.InnerError, err)
		return
	}

	totalSeconds := now.Sub(time.Unix(firstSeen, 0)).Seconds()
	if totalSeconds <= 0 {
		return
	}
	availability := 100 * (1 - downSeconds/totalSeconds)
	if availability < 0 {
		availability = 0
	}
	if availability > 100 {
		availability = 100
	}
	metrics.SetAvailabilityPercent(availability)
}

func (s *StatsCollector) getOrSetUnix(key string, value int64) (int64, error) {
	raw, err := s.redis.Get(s.ctx, key).Result()
	if err == redis.Nil {
		if err := s.redis.Set(s.ctx, key, value, 0).Err(); err != nil {
			return 0, err
		}
		return value, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(raw, 10, 64)
}

func (s *StatsCollector) getFloat64(key string) (float64, error) {
	raw, err := s.redis.Get(s.ctx, key).Result()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return strconv.ParseFloat(raw, 64)
}

func (s *StatsCollector) getTotalUsersCount() (int, error) {
	collection := database.CollectionCocks(s.db)
	cursor, err := collection.Aggregate(s.ctx, database.PipelineTotalCockersCount())
	if err != nil {
		return 0, err
	}
	defer cursor.Close(s.ctx)

	var result struct {
		Total int `bson:"total"`
	}
	if cursor.Next(s.ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.Total, nil
	}
	return 0, nil
}

func (s *StatsCollector) getActiveUsersCount(since time.Time) (int, error) {
	collection := database.CollectionCocks(s.db)
	cursor, err := collection.Aggregate(s.ctx, database.PipelineActiveUsersSince(since))
	if err != nil {
		return 0, err
	}
	defer cursor.Close(s.ctx)

	var result struct {
		Total int `bson:"total"`
	}
	if cursor.Next(s.ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, err
		}
		return result.Total, nil
	}
	return 0, nil
}
