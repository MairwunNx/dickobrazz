package application

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"sort"
	"strings"
	"time"
)

func (app *Application) HandleInlineQuery(log *Logger, query *tgbotapi.InlineQuery) {
	var results []any
	if query.From.ID == 362695653 {
		results = []any{
			app.InlineQueryCockSize(log, query),
			app.InlineQueryCockRace(log, query),
			app.InlineQueryCockRuler(log, query),
			//app.InlineQueryCockRaceImgStat(log, query),
			app.InlineQueryCockDynamic(log, query),
		}
	} else {
		results = []any{
			app.InlineQueryCockSize(log, query),
			app.InlineQueryCockRace(log, query),
			app.InlineQueryCockRuler(log, query),
			//app.InlineQueryCockRaceImgStat(log, query),
			app.InlineQueryCockDynamic(log, query),
		}
	}

	inlines := tgbotapi.InlineConfig{
		InlineQueryID: query.ID,
		IsPersonal:    true,
		CacheTime:     60,
		Results:       results,
	}

	if _, err := app.bot.Request(inlines); err != nil {
		log.E("Failed to send inline query", InnerError, err)
	} else {
		log.I("Inline query successfully sent")
	}
}

func (app *Application) InlineQueryCockSize(log *Logger, query *tgbotapi.InlineQuery) tgbotapi.InlineQueryResultArticle {
	var size int
	//
	//if daily, isPresent := app.GetDaylyCock(log, query.From.ID); isPresent {
	//
	//}

	if cached := app.GetCockSizeFromCache(log, query.From.ID); cached != nil {
		size = *cached
	} else {
		size = app.rnd.IntN(log, 60)

		cock := &Cock{
			ID:          uuid.NewString(),
			Size:        int32(size),
			Nickname:    query.From.UserName,
			UserID:      query.From.ID,
			RequestedAt: NowTime(),
		}

		app.SaveCockToCache(log, query.From.ID, query.From.UserName, size)
		app.SaveCockToMongo(log, cock)
	}

	emoji := EmojiFromSize(size)
	text := GenerateCockSizeText(size, emoji)

	return InitializeInlineQuery(
		"Размер кока",
		strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(text, ".", "\\."), "-", "\\-"), "!", "\\!"),
	)
}

func (app *Application) InlineQueryCockRace(log *Logger, query *tgbotapi.InlineQuery) tgbotapi.InlineQueryResultArticle {
	cocks := app.AggregateCockSizes(log)
	text := app.GenerateCockRaceScoreboard(log, query.From.ID, cocks)
	return InitializeInlineQuery("Гонка коков", strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(text, ".", "\\."), "-", "\\-"), "!", "\\!"))
}

const bigCockThreshold = 19

func (app *Application) InlineQueryCockDynamic(log *Logger, query *tgbotapi.InlineQuery) tgbotapi.InlineQueryResultArticle {
	/*collection := app.db.Database("dickbot_db").Collection("cocks")

	userPipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "user_id", Value: query.From.ID}}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$requested_at"},
			{Key: "totalSize", Value: bson.D{{Key: "$sum", Value: "$size"}}},
			{Key: "sizes", Value: bson.D{{Key: "$push", Value: "$size"}}},
			{Key: "avgSize", Value: bson.D{{Key: "$avg", Value: "$size"}}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
	}

	cursor, err := collection.Aggregate(app.ctx, userPipeline)
	if err != nil {
		log.E("Failed to aggregate user cock data", InnerError, err)
		return tgbotapi.InlineQueryResultArticle{}
	}

	var userResults []struct {
		Date      time.Time `bson:"_id"`
		TotalSize int       `bson:"totalSize"`
		Sizes     []int     `bson:"sizes"`
		AvgSize   float64   `bson:"avgSize"`
		Count     int       `bson:"count"`
	}
	if err := cursor.All(app.ctx, &userResults); err != nil {
		log.E("Failed to decode user cock data", InnerError, err)
		return tgbotapi.InlineQueryResultArticle{}
	}

	log.I("Successfully aggregated user cock data")

	averagePipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$user_id"},
			{Key: "avgSize", Value: bson.D{{Key: "$avg", Value: "$size"}}},
		}}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "totalAvgSize", Value: bson.D{{Key: "$avg", Value: "$avgSize"}}},
		}}},
	}

	averageCursor, err := collection.Aggregate(app.ctx, averagePipeline)
	if err != nil {
		log.E("Failed to aggregate global average cock size", InnerError, err)
		return tgbotapi.InlineQueryResultArticle{}
	}

	var globalResult struct {
		TotalAvgSize float64 `bson:"totalAvgSize"`
	}
	if averageCursor.Next(app.ctx) {
		if err := averageCursor.Decode(&globalResult); err != nil {
			log.E("Failed to decode global average data", InnerError, err)
			return tgbotapi.InlineQueryResultArticle{}
		}
	} else {
		log.E("No global average data found")
		return tgbotapi.InlineQueryResultArticle{}
	}

	log.I("Successfully calculated global average cock size", "TotalAvgSize", globalResult.TotalAvgSize)

	// Global total pipeline
	totalPipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "totalCock", Value: bson.D{{Key: "$sum", Value: "$size"}}},
		}}},
	}

	totalCursor, err := collection.Aggregate(app.ctx, totalPipeline)
	if err != nil {
		log.E("Failed to aggregate global total cock size", InnerError, err)
		return tgbotapi.InlineQueryResultArticle{}
	}

	var totalResult struct {
		TotalCock int `bson:"totalCock"`
	}
	if totalCursor.Next(app.ctx) {
		if err := totalCursor.Decode(&totalResult); err != nil {
			log.E("Failed to decode global total data", InnerError, err)
			return tgbotapi.InlineQueryResultArticle{}
		}
	} else {
		log.E("No global total data found")
		return tgbotapi.InlineQueryResultArticle{}
	}

	log.I("Successfully calculated global total cock size", "TotalCock", totalResult.TotalCock)

	// Pipeline для подсчёта уникальных пользователей
	uniqueUsersPipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$user_id"},
		}}},
		{{Key: "$count", Value: "count"}}, // Подсчитываем количество уникальных групп
	}

	// Выполняем агрегацию для подсчёта уникальных пользователей
	uniqueUsersCursor, err := collection.Aggregate(app.ctx, uniqueUsersPipeline)
	if err != nil {
		log.E("Failed to aggregate unique users", InnerError, err)
		return tgbotapi.InlineQueryResultArticle{}
	}

	var uniqueUsersResult struct {
		Count int `bson:"count"`
	}
	if uniqueUsersCursor.Next(app.ctx) {
		if err := uniqueUsersCursor.Decode(&uniqueUsersResult); err != nil {
			log.E("Failed to decode unique users data", InnerError, err)
			return tgbotapi.InlineQueryResultArticle{}
		}
	} else {
		log.E("No unique users data found")
		return tgbotapi.InlineQueryResultArticle{}
	}

	log.I("Successfully calculated unique users", "Count", uniqueUsersResult.Count)

	distributionPipeline := mongo.Pipeline{
		{{Key: "$facet", Value: bson.D{
			{Key: "bigCocks", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "size", Value: bson.D{{Key: "$gte", Value: bigCockThreshold}}}}}},
				bson.D{{Key: "$count", Value: "count"}},
			}},
			{Key: "smallCocks", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "size", Value: bson.D{{Key: "$lt", Value: bigCockThreshold}}}}}},
				bson.D{{Key: "$count", Value: "count"}},
			}},
		}}},
	}

	distributionCursor, err := collection.Aggregate(app.ctx, distributionPipeline)
	if err != nil {
		log.E("Failed to calculate cock size distribution", InnerError, err)
		return tgbotapi.InlineQueryResultArticle{}
	}

	var distributionResults []struct {
		BigCocks   []struct{ Count int } `bson:"bigCocks"`
		SmallCocks []struct{ Count int } `bson:"smallCocks"`
	}
	if err := distributionCursor.All(app.ctx, &distributionResults); err != nil {
		log.E("Failed to decode cock size distribution data", InnerError, err)
		return tgbotapi.InlineQueryResultArticle{}
	}

	log.I("Successfully calculated cock size distribution", "Results", distributionResults)

	// Подсчитываем "Большие" и "Маленькие" коки
	var bigCocks, smallCocks int
	if len(distributionResults) > 0 {
		if len(distributionResults[0].BigCocks) > 0 {
			bigCocks = distributionResults[0].BigCocks[0].Count
		}
		if len(distributionResults[0].SmallCocks) > 0 {
			smallCocks = distributionResults[0].SmallCocks[0].Count
		}
	}

	// Рассчитываем проценты
	var bigCocksPercent, smallCocksPercent float64
	totalCocks := bigCocks + smallCocks
	if totalCocks > 0 {
		bigCocksPercent = float64(bigCocks) / float64(totalCocks) * 100
		smallCocksPercent = float64(smallCocks) / float64(totalCocks) * 100
	}

	// Calculate metrics
	var totalCock, totalAvgCock, totalMedianCock int
	var userTotalCock, userAvgCock, userMaxCock, userYesterdayChangeCock int
	var userIrk, userYesterdayChangePercent, userDailyGrowth float64
	var userMaxCockDate time.Time

	globalCockPipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: nil},
			{Key: "sizes", Value: bson.D{{Key: "$push", Value: "$size"}}},
		}}},
	}

	var globalCocksResult struct {
		Sizes []int `bson:"sizes"`
	}

	globalCockCursor, err := collection.Aggregate(app.ctx, globalCockPipeline)
	if err != nil {
		log.E("Failed to aggregate global cock sizes for median", InnerError, err)
		return tgbotapi.InlineQueryResultArticle{}
	}

	if globalCockCursor.Next(app.ctx) {
		if err := globalCockCursor.Decode(&globalCocksResult); err != nil {
			log.E("Failed to decode global cock sizes for median", InnerError, err)
			return tgbotapi.InlineQueryResultArticle{}
		}
	}

	totalMedianCock = median(globalCocksResult.Sizes)

	var allCocks []int
	for _, result := range userResults {
		allCocks = append(allCocks, result.Sizes...)
		userTotalCock += result.TotalSize

		// Track max cock
		for _, size := range result.Sizes {
			if size >= userMaxCock {
				userMaxCock = size
				userMaxCockDate = result.Date
			}
		}
	}

	// Calculate overall metrics
	totalCock = totalResult.TotalCock

	if len(userResults) > 0 {
		userAvgCock = int(float64(userTotalCock) / float64(len(userResults)))
	}

	if len(allCocks) > 0 {
		totalAvgCock = int(globalResult.TotalAvgSize)
	}

	// Calculate IRK
	if totalCock > 0 && userTotalCock > 0 {
		// Нормализуем пользовательский общий кок относительно среднего общего размера
		normalizedUserCock := float64(userTotalCock) / float64(totalAvgCock)

		// Нормализуем количество записей пользователя
		normalizedUserRecords := float64(len(userResults)) / float64(len(allCocks))

		// Динамические веса (с ограничением)
		w1 := math.Max(1.0, math.Min(normalizedUserCock*2.0, 10.0))
		w2 := math.Max(1.0, math.Min(normalizedUserRecords*5.0, 10.0))

		// Вычисляем IRK
		rawIrk := normalizedUserCock / (1.0 + w1) * (normalizedUserRecords / (1.0 + w2))

		// Ограничиваем IRK в пределах [0.0, 1.0]
		userIrk = math.Max(0.0, math.Min(1.0, rawIrk))
	}

	// Calculate yesterday's change
	if len(userResults) > 1 {
		userYesterdayChangeCock = userResults[len(userResults)-1].TotalSize - userResults[len(userResults)-2].TotalSize
		if userResults[len(userResults)-2].TotalSize > 0 {
			userYesterdayChangePercent = float64(userYesterdayChangeCock) / float64(userResults[len(userResults)-2].TotalSize) * 100
		} else {
			userYesterdayChangePercent = 100
		}
	} else {
		userYesterdayChangePercent = 100
		userYesterdayChangeCock = 0
	}

	// Calculate daily growth
	var dailyGrowthSum float64
	for i := 1; i < len(userResults); i++ {
		dailyGrowthSum += float64(userResults[i].TotalSize - userResults[i-1].TotalSize)
	}
	if len(userResults) > 1 {
		userDailyGrowth = dailyGrowthSum / float64(len(userResults)-1)
	}

	// Pipeline для нахождения дня с самым большим коком
	maxCockDayPipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "year", Value: bson.D{{Key: "$year", Value: "$requested_at"}}},
				{Key: "month", Value: bson.D{{Key: "$month", Value: "$requested_at"}}},
				{Key: "day", Value: bson.D{{Key: "$dayOfMonth", Value: "$requested_at"}}},
			}},
			{Key: "totalSize", Value: bson.D{{Key: "$sum", Value: "$size"}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "totalSize", Value: -1}}}},
		{{Key: "$limit", Value: 1}},
	}

	maxCockDayPipelineCursor, err := collection.Aggregate(app.ctx, maxCockDayPipeline)
	if err != nil {
		log.E("Failed to aggregate max cock day", InnerError, err)
		return tgbotapi.InlineQueryResultArticle{}
	}

	var maxCockDayResult struct {
		ID struct {
			Year  int `bson:"year"`
			Month int `bson:"month"`
			Day   int `bson:"day"`
		} `bson:"_id"`
		TotalSize int `bson:"totalSize"`
	}

	if maxCockDayPipelineCursor.Next(app.ctx) {
		if err := maxCockDayPipelineCursor.Decode(&maxCockDayResult); err != nil {
			log.E("Failed to decode max cock day data", InnerError, err)
			return tgbotapi.InlineQueryResultArticle{}
		}
	} else {
		log.E("No max cock day data found")
		return tgbotapi.InlineQueryResultArticle{}
	}

	maxCockDate := time.Date(maxCockDayResult.ID.Year, time.Month(maxCockDayResult.ID.Month), maxCockDayResult.ID.Day, 0, 0, 0, 0, time.Local)
	maxCockSize := maxCockDayResult.TotalSize

	log.I("Successfully calculated max cock day", "Date", maxCockDate, "TotalSize", maxCockSize)

	var dominancePercent float64
	if totalResult.TotalCock > 0 {
		dominancePercent = (float64(userTotalCock) / float64(totalResult.TotalCock)) * 100
	}

	log.I("Successfully calculated user dominance percentage", "Dominance", dominancePercent)

	// Generate result text
	text := NewMsgCockDynamicsTemplate(
		// Общая динамика коков
		totalCock, uniqueUsersResult.Count, totalAvgCock, totalMedianCock,
		// Персональная динамика кока
		userTotalCock, userAvgCock, userIrk, userMaxCock, userMaxCockDate,
		// Кок-активы
		userYesterdayChangePercent, userYesterdayChangeCock,
		userDailyGrowth,
		bigCocksPercent, smallCocksPercent,
		maxCockDate, maxCockSize,
		dominancePercent,
	)

	return tgbotapi.NewInlineQueryResultArticleMarkdown(query.ID, "Динамика кока", text)*/

	collection := app.db.Database("dickbot_db").Collection("cocks")

	pipeline := mongo.Pipeline{
		{{Key: "$facet", Value: bson.D{
			{Key: "userStats", Value: bson.A{
				bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: query.From.ID}}}},
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: "$requested_at"},
					{Key: "totalSize", Value: bson.D{{Key: "$sum", Value: "$size"}}},
					{Key: "sizes", Value: bson.D{{Key: "$push", Value: "$size"}}},
					{Key: "avgSize", Value: bson.D{{Key: "$avg", Value: "$size"}}},
					{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
				}}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
			}},
			{Key: "globalStats", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "totalCock", Value: bson.D{{Key: "$sum", Value: "$size"}}},
					{Key: "avgSize", Value: bson.D{{Key: "$avg", Value: "$size"}}},
					{Key: "sizes", Value: bson.D{{Key: "$push", Value: "$size"}}},
				}}},
			}},
			{Key: "uniqueUsers", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$user_id"}}}},
				bson.D{{Key: "$count", Value: "count"}},
			}},
			{Key: "cockDistribution", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: nil},
					{Key: "bigCocks", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$gte", Value: bson.A{"$size", bigCockThreshold}}}, 1, 0,
					}}}}}},
					{Key: "smallCocks", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{
						bson.D{{Key: "$lt", Value: bson.A{"$size", bigCockThreshold}}}, 1, 0,
					}}}}}},
				}}},
			}},
			{Key: "maxCockDay", Value: bson.A{
				bson.D{{Key: "$group", Value: bson.D{
					{Key: "_id", Value: bson.D{
						{Key: "year", Value: bson.D{{Key: "$year", Value: "$requested_at"}}},
						{Key: "month", Value: bson.D{{Key: "$month", Value: "$requested_at"}}},
						{Key: "day", Value: bson.D{{Key: "$dayOfMonth", Value: "$requested_at"}}},
					}},
					{Key: "totalSize", Value: bson.D{{Key: "$sum", Value: "$size"}}},
				}}},
				bson.D{{Key: "$sort", Value: bson.D{{Key: "totalSize", Value: -1}}}},
				bson.D{{Key: "$limit", Value: 1}},
			}},
		}}},
	}

	cursor, err := collection.Aggregate(app.ctx, pipeline)
	if err != nil {
		log.E("Failed to aggregate data with facet pipeline", InnerError, err)
		return tgbotapi.InlineQueryResultArticle{}
	}

	var facetResults []struct {
		UserStats []struct {
			Date      time.Time `bson:"_id"`
			TotalSize int       `bson:"totalSize"`
			Sizes     []int     `bson:"sizes"`
			AvgSize   float64   `bson:"avgSize"`
			Count     int       `bson:"count"`
		} `bson:"userStats"`
		GlobalStats []struct {
			TotalCock int     `bson:"totalCock"`
			AvgSize   float64 `bson:"avgSize"`
			Sizes     []int   `bson:"sizes"`
		} `bson:"globalStats"`
		UniqueUsers []struct {
			Count int `bson:"count"`
		} `bson:"uniqueUsers"`
		CockDistribution []struct {
			BigCocks   int `bson:"bigCocks"`
			SmallCocks int `bson:"smallCocks"`
		} `bson:"cockDistribution"`
		MaxCockDay []struct {
			ID struct {
				Year  int `bson:"year"`
				Month int `bson:"month"`
				Day   int `bson:"day"`
			} `bson:"_id"`
			TotalSize int `bson:"totalSize"`
		} `bson:"maxCockDay"`
	}

	if err := cursor.All(app.ctx, &facetResults); err != nil {
		log.E("Failed to decode facet pipeline results", InnerError, err)
		return tgbotapi.InlineQueryResultArticle{}
	}

	s := facetResults[0]

	// Extract data from facet results
	userResults := s.UserStats
	globalStats := s.GlobalStats[0]
	uniqueUsers := s.UniqueUsers[0].Count
	cockDistribution := s.CockDistribution[0]
	maxCockDay := s.MaxCockDay[0]

	// Metrics initialization
	var userTotalCock, userAvgCock, userMaxCock, userYesterdayChangeCock int
	var userIrk, userYesterdayChangePercent, userDailyGrowth, bigCocksPercent, smallCocksPercent, dominancePercent float64
	var userMaxCockDate time.Time
	var totalCock, totalAvgCock, totalMedianCock int
	var maxCockDate time.Time
	var maxCockSize int

	// Process user stats
	for _, result := range userResults {
		userTotalCock += result.TotalSize

		// Track max cock
		for _, size := range result.Sizes {
			if size > userMaxCock {
				userMaxCock = size
				userMaxCockDate = result.Date
			}
		}
	}

	// Calculate global metrics
	totalCock = globalStats.TotalCock
	totalAvgCock = int(globalStats.AvgSize)
	totalMedianCock = median(globalStats.Sizes)

	// Calculate IRK
	var allCocks []int
	for _, result := range s.UserStats {
		allCocks = append(allCocks, result.Sizes...)
	}

	// Calculate IRK
	if totalCock > 0 && userTotalCock > 0 && len(allCocks) > 0 {
		// Нормализуем пользовательский общий кок относительно среднего общего размера
		normalizedUserCock := float64(userTotalCock) / float64(totalAvgCock)

		// Нормализуем количество записей пользователя
		normalizedUserRecords := float64(len(userResults)) / float64(len(allCocks))

		// Динамические веса (с ограничением)
		w1 := math.Max(1.0, math.Min(normalizedUserCock*2.0, 10.0))
		w2 := math.Max(1.0, math.Min(normalizedUserRecords*5.0, 10.0))

		// Вычисляем IRK
		rawIrk := normalizedUserCock / (1.0 + w1) * (normalizedUserRecords / (1.0 + w2))

		// Ограничиваем IRK в пределах [0.0, 1.0]
		userIrk = math.Max(0.0, math.Min(1.0, rawIrk))
	}

	// Calculate user's average cock size
	if len(userResults) > 0 {
		userAvgCock = int(float64(userTotalCock) / float64(len(userResults)))
	}

	// Calculate yesterday's change
	if len(userResults) > 1 {
		userYesterdayChangeCock = userResults[len(userResults)-1].TotalSize - userResults[len(userResults)-2].TotalSize
		if userResults[len(userResults)-2].TotalSize > 0 {
			userYesterdayChangePercent = float64(userYesterdayChangeCock) / float64(userResults[len(userResults)-2].TotalSize) * 100
		} else {
			userYesterdayChangePercent = 100
		}
	} else {
		userYesterdayChangePercent = 100
		userYesterdayChangeCock = 0
	}

	// Calculate daily growth
	var dailyGrowthSum float64
	for i := 1; i < len(userResults); i++ {
		dailyGrowthSum += float64(userResults[i].TotalSize - userResults[i-1].TotalSize)
	}
	if len(userResults) > 1 {
		userDailyGrowth = dailyGrowthSum / float64(len(userResults)-1)
	}

	// Calculate distribution percentages
	bigCocks := cockDistribution.BigCocks
	smallCocks := cockDistribution.SmallCocks
	totalCocks := bigCocks + smallCocks
	if totalCocks > 0 {
		bigCocksPercent = float64(bigCocks) / float64(totalCocks) * 100
		smallCocksPercent = float64(smallCocks) / float64(totalCocks) * 100
	}

	// Extract max cock day data
	maxCockDate = time.Date(maxCockDay.ID.Year, time.Month(maxCockDay.ID.Month), maxCockDay.ID.Day, 0, 0, 0, 0, time.Local)
	maxCockSize = maxCockDay.TotalSize

	// Calculate dominance percentage
	if totalCock > 0 {
		dominancePercent = (float64(userTotalCock) / float64(totalCock)) * 100
	}

	// Generate result text
	text := NewMsgCockDynamicsTemplate(
		// Общая динамика коков
		totalCock, uniqueUsers, totalAvgCock, totalMedianCock,
		// Персональная динамика кока
		userTotalCock, userAvgCock, userIrk, userMaxCock, userMaxCockDate,
		// Кок-активы
		userYesterdayChangePercent, userYesterdayChangeCock,
		userDailyGrowth,
		bigCocksPercent, smallCocksPercent,
		maxCockDate, maxCockSize,
		dominancePercent,
	)

	return tgbotapi.NewInlineQueryResultArticleMarkdown(query.ID, "Динамика кока", text)
}

func sum(data []int) int {
	total := 0
	for _, v := range data {
		total += v
	}
	return total
}

func median(data []int) int {
	n := len(data)
	if n == 0 {
		return 0
	}

	sort.Ints(data)
	if n%2 == 0 {
		return (data[n/2-1] + data[n/2]) / 2
	}
	return data[n/2]
}

func (app *Application) InlineQueryCockRaceImgStat(log *Logger, query *tgbotapi.InlineQuery) tgbotapi.InlineQueryResultPhoto {
	photo := tgbotapi.NewInlineQueryResultPhotoWithThumb(uuid.NewString(),
		"https://files.lynguard.com/raw/public/work-avatar.jpg",
		"https://files.lynguard.com/raw/public/work-avatar.jpg",
	)
	photo.Caption = "Тест photo.Caption"
	photo.Description = "Тест photo.Description"
	photo.Title = "Тест photo.Title"
	return photo
}

func (app *Application) InlineQueryCockRuler(log *Logger, query *tgbotapi.InlineQuery) tgbotapi.InlineQueryResultArticle {
	cocks := app.GetCockSizesFromCache(log)

	sort.Slice(cocks, func(i, j int) bool {
		return cocks[i].Size > cocks[j].Size
	})

	if len(cocks) > 13 {
		cocks = cocks[:13]
	}

	text := app.GenerateCockRulerText(log, query.From.ID, cocks)
	return InitializeInlineQuery("Линейка коков", strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(text, ".", "\\."), "-", "\\-"), "!", "\\!"))
}

func InitializeInlineQuery(title, message string) tgbotapi.InlineQueryResultArticle {
	return tgbotapi.NewInlineQueryResultArticleMarkdownV2(uuid.NewString(), title, message)
}
