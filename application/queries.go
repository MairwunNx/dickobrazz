package application

import (
	"fmt"
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

	log.I("**** TEST", "NowTime", NowTime())

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

func (app *Application) InlineQueryCockDynamic(log *Logger, query *tgbotapi.InlineQuery) tgbotapi.InlineQueryResultArticle {
	collection := TraceTimeExecutionForResult(log, TraceKindGatherCollection, func() *mongo.Collection {
		return app.db.Database("dickbot_db").Collection("cocks")
	})

	pipeline := TraceTimeExecutionForResult(log, TraceKindCreatePipeline, func() mongo.Pipeline {
		return mongo.Pipeline{
			{{Key: "$facet", Value: bson.D{
				{Key: "individual", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: query.From.ID}}}},
					bson.D{{Key: "$group", Value: bson.D{
						{Key: "_id", Value: "$requested_at"},
						{Key: "total", Value: bson.D{{Key: "$sum", Value: "$size"}}},
						{Key: "sizes", Value: bson.D{{Key: "$push", Value: "$size"}}},
						{Key: "average", Value: bson.D{{Key: "$avg", Value: "$size"}}},
					}}},
					bson.D{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
					bson.D{{Key: "$project", Value: bson.D{
						{Key: "_id", Value: 1},
						{Key: "total", Value: 1},
						{Key: "sizes", Value: 1},
						{Key: "average", Value: bson.D{{Key: "$round", Value: bson.A{"$average", 0}}}},
					}}},
				}},
				{Key: "individual_cock", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: query.From.ID}}}},
					bson.D{{Key: "$group", Value: bson.D{
						{Key: "_id", Value: nil},
						{Key: "total", Value: bson.D{{Key: "$sum", Value: "$size"}}},
						{Key: "avg_val", Value: bson.D{{Key: "$avg", Value: "$size"}}},
					}}},
					bson.D{{Key: "$project", Value: bson.D{
						{Key: "_id", Value: 0},
						{Key: "total", Value: 1},
						{Key: "average", Value: bson.D{{Key: "$round", Value: bson.A{"$avg_val", 0}}}},
					}}},
				}},
				{Key: "overall", Value: bson.A{
					bson.D{{Key: "$group", Value: bson.D{
						{Key: "_id", Value: nil},
						{Key: "size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
						{Key: "median", Value: bson.D{{Key: "$median", Value: bson.D{
							{Key: "input", Value: "$size"},
							{Key: "method", Value: "approximate"},
						}}}},
						{Key: "average", Value: bson.D{{Key: "$avg", Value: "$size"}}},
					}}},
					bson.D{{Key: "$project", Value: bson.D{
						{Key: "_id", Value: nil},
						{Key: "size", Value: 1},
						{Key: "median", Value: 1},
						{Key: "average", Value: bson.D{{Key: "$round", Value: bson.A{"$average", 0}}}},
					}}},
					bson.D{{Key: "$limit", Value: 1}},
				}},
				{Key: "uniques", Value: bson.A{
					bson.D{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$user_id"}}}},
					bson.D{{Key: "$count", Value: "count"}},
				}},
				{Key: "distribution", Value: bson.A{
					bson.D{{Key: "$group", Value: bson.D{
						{Key: "_id", Value: nil},
						{Key: "huge", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{
							bson.D{{Key: "$gte", Value: bson.A{"$size", 19}}}, 1, 0,
						}}}}}},
						{Key: "little", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{
							bson.D{{Key: "$lt", Value: bson.A{"$size", 19}}}, 1, 0,
						}}}}}},
					}}},
					bson.D{{Key: "$addFields", Value: bson.D{
						{Key: "total", Value: bson.D{{Key: "$add", Value: bson.A{"$huge", "$little"}}}},
					}}},
					bson.D{{Key: "$project", Value: bson.D{
						{Key: "_id", Value: nil},
						{Key: "huge", Value: bson.D{{Key: "$cond", Value: bson.A{
							bson.D{{Key: "$eq", Value: bson.A{"$total", 0}}},
							0,
							bson.D{{Key: "$multiply", Value: bson.A{bson.D{{Key: "$divide", Value: bson.A{"$huge", "$total"}}}, 100}}},
						}}}},
						{Key: "little", Value: bson.D{{Key: "$cond", Value: bson.A{
							bson.D{{Key: "$eq", Value: bson.A{"$total", 0}}},
							0,
							bson.D{{Key: "$multiply", Value: bson.A{bson.D{{Key: "$divide", Value: bson.A{"$little", "$total"}}}, 100}}},
						}}}},
					}}},
				}},
				{Key: "record", Value: bson.A{
					bson.D{{Key: "$group", Value: bson.D{
						{Key: "_id", Value: bson.D{
							{Key: "year", Value: bson.D{{Key: "$year", Value: "$requested_at"}}},
							{Key: "month", Value: bson.D{{Key: "$month", Value: "$requested_at"}}},
							{Key: "day", Value: bson.D{{Key: "$dayOfMonth", Value: "$requested_at"}}},
						}},
						{Key: "requested_at", Value: bson.D{{Key: "$first", Value: "$requested_at"}}},
						{Key: "total", Value: bson.D{{Key: "$sum", Value: "$size"}}},
					}}},
					bson.D{{Key: "$sort", Value: bson.D{{Key: "total", Value: -1}}}},
					bson.D{{Key: "$limit", Value: 1}},
				}},
				{Key: "individual_record", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: query.From.ID}}}},
					bson.D{{Key: "$group", Value: bson.D{
						{Key: "_id", Value: "$requested_at"},
						{Key: "requested_at", Value: bson.D{{Key: "$first", Value: "$requested_at"}}},
						{Key: "total", Value: bson.D{{Key: "$first", Value: "$size"}}},
					}}},
					bson.D{{Key: "$sort", Value: bson.D{{Key: "total", Value: -1}}}},
					bson.D{{Key: "$limit", Value: 1}},
				}},
			}}}}
	})

	cursor, err := TraceTimeExecutionForResultError(log, TraceKindAggregatePipeline, func() (*mongo.Cursor, error) {
		return collection.Aggregate(app.ctx, pipeline)
	})

	if err != nil {
		log.E("Aggregation failed", InnerError, err)
		return tgbotapi.InlineQueryResultArticle{}
	}

	var result struct {
		Individual []struct {
			Date    time.Time `bson:"_id"`
			Total   int       `bson:"total"`
			Sizes   []int     `bson:"sizes"`
			Average int       `bson:"average"`
		} `bson:"individual"`

		IndividualCock []struct {
			Total   int `bson:"total"`
			Average int `bson:"average"`
		} `bson:"individual_cock"`

		IndividualRecord []struct {
			RequestedAt time.Time `bson:"requested_at"`
			Total       int       `bson:"total"`
		} `bson:"individual_record"`

		Overall []struct {
			Size    int `bson:"size"`
			Average int `bson:"average"`
			Median  int `bson:"median"`
		} `bson:"overall"`

		Uniques []struct {
			Count int `bson:"count"`
		} `bson:"uniques"`

		Distribution []struct {
			HugePercent   float64 `bson:"huge"`
			LittlePercent float64 `bson:"little"`
		} `bson:"distribution"`

		Record []struct {
			RequestedAt time.Time `bson:"requested_at"`
			Total       int       `bson:"total"`
		} `bson:"record"`
	}

	if err := TraceTimeExecutionForResult(log, TraceKindInflatePipeline, func() error {
		if cursor.Next(app.ctx) {
			return cursor.Decode(&result)
		}
		return fmt.Errorf("no results found in aggregation")
	}); err != nil {
		log.E("Failed to decode aggregation results", InnerError, err)
		return tgbotapi.InlineQueryResultArticle{}
	}

	log.I("Aggregation completed successfully")

	log.I("***** IndividualRecord", "IndividualRecord", result.IndividualRecord)

	individual := result.Individual
	overall := result.Overall[0]

	overallCockers := result.Uniques[0].Count
	distribution := result.Distribution[0]
	record := result.Record[0]
	individualRecord := result.IndividualRecord[0]
	individualCock := result.IndividualCock[0]

	// Metrics initialization
	var totalUserCock, yesterdayCockChange int
	var irk, yesterdayChangePercent, dailyGrowth, dominancePercent float64
	var totalCock, avgCock, medianCock int

	// Calculate global metrics
	totalCock = overall.Size
	avgCock = overall.Average
	medianCock = overall.Median

	// Gather all individual cocks
	var userCocks []int
	for _, stat := range individual {
		userCocks = append(userCocks, stat.Sizes...)
	}

	// todo: remove

	totalUserCock = individualCock.Total

	// Calculate IRK
	if totalCock > 0 && totalUserCock > 0 && len(userCocks) > 0 {
		normalizedCock := float64(totalUserCock) / float64(avgCock)
		normalizedRecords := float64(len(individual)) / float64(len(userCocks))

		w1 := math.Max(1.0, math.Min(normalizedCock*2.0, 10.0))
		w2 := math.Max(1.0, math.Min(normalizedRecords*5.0, 10.0))

		rawIrk := normalizedCock / (1.0 + w1) * (normalizedRecords / (1.0 + w2))
		irk = math.Max(0.0, math.Min(1.0, rawIrk))
	}

	// Calculate yesterday's change
	if len(individual) > 1 {
		yesterdayCockChange = individual[len(individual)-1].Total - individual[len(individual)-2].Total
		if individual[len(individual)-2].Total > 0 {
			yesterdayChangePercent = float64(yesterdayCockChange) / float64(individual[len(individual)-2].Total) * 100
		} else {
			yesterdayChangePercent = 100
		}
	} else {
		yesterdayChangePercent = 100
		yesterdayCockChange = 0
	}

	// Calculate daily growth
	var growthSum float64
	for i := 1; i < len(individual); i++ {
		growthSum += float64(individual[i].Total - individual[i-1].Total)
	}
	if len(individual) > 1 {
		dailyGrowth = growthSum / float64(len(individual)-1)
	}

	// Calculate dominance percentage
	if totalCock > 0 {
		dominancePercent = (float64(totalUserCock) / float64(totalCock)) * 100
	}

	// Generate result text
	text := NewMsgCockDynamicsTemplate(
		// Общая динамика коков
		totalCock, overallCockers, avgCock, medianCock,
		// Персональная динамика кока
		totalUserCock, individualCock.Average, irk, individualRecord.Total, individualRecord.RequestedAt,
		// Кок-активы
		yesterdayChangePercent, yesterdayCockChange,
		dailyGrowth,
		distribution.HugePercent, distribution.LittlePercent,
		record.RequestedAt, record.Total,
		dominancePercent,
	)

	return tgbotapi.NewInlineQueryResultArticleMarkdown(query.ID, "Динамика кока", text)
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
