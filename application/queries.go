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

const (
	k1 = 0.1 // Коэффициент влияния крупного размера
	k2 = 0.5 // Коэффициент влияния количества записей
)

func (app *Application) InlineQueryCockDynamic(log *Logger, query *tgbotapi.InlineQuery) tgbotapi.InlineQueryResultArticle {
	collection := app.db.Database("dickbot_db").Collection("cocks")

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
			if size > userMaxCock {
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

	// Generate result text
	text := NewMsgCockDynamicsTemplate(
		// Общая динамика коков
		totalCock, len(allCocks), totalAvgCock, totalMedianCock,
		// Персональная динамика кока
		userTotalCock, userAvgCock, userIrk, userMaxCock, userMaxCockDate,
		// Кок-активы
		userYesterdayChangePercent, userYesterdayChangeCock,
		userDailyGrowth,
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
