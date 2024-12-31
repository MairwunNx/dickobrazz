package application

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
				{Key: "individual_irk", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: query.From.ID}}}},
					bson.D{{Key: "$group", Value: bson.D{
						{Key: "_id", Value: nil},
						{Key: "user_total_size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
					}}},
					bson.D{{Key: "$lookup", Value: bson.D{
						{Key: "from", Value: "cocks"},
						{Key: "pipeline", Value: bson.A{
							bson.D{{Key: "$group", Value: bson.D{
								{Key: "_id", Value: nil},
								{Key: "global_total_size", Value: bson.D{{Key: "$sum", Value: "$size"}}},
							}}},
						}},
						{Key: "as", Value: "global_data"},
					}}},
					bson.D{{Key: "$unwind", Value: "$global_data"}},
					bson.D{{Key: "$project", Value: bson.D{
						{Key: "_id", Value: nil},
						{Key: "irk", Value: bson.D{{Key: "$round", Value: bson.A{
							bson.D{{Key: "$cond", Value: bson.A{
								bson.D{{Key: "$lte", Value: bson.A{"$global_data.global_total_size", 0}}},
								0,
								bson.D{{Key: "$divide", Value: bson.A{
									bson.D{{Key: "$log10", Value: bson.D{{Key: "$add", Value: bson.A{1, "$user_total_size"}}}}},
									bson.D{{Key: "$log10", Value: bson.D{{Key: "$add", Value: bson.A{1, "$global_data.global_total_size"}}}}},
								}}},
							}}},
							3,
						}}}},
					}}},
				}},
				{Key: "individual_dominance", Value: bson.A{
					bson.D{{Key: "$group", Value: bson.D{
						{Key: "_id", Value: nil},
						{Key: "total_cock", Value: bson.D{{Key: "$sum", Value: "$size"}}},
						{Key: "total_user_cock", Value: bson.D{{Key: "$sum", Value: bson.D{{Key: "$cond", Value: bson.A{
							bson.D{{Key: "$eq", Value: bson.A{"$user_id", query.From.ID}}},
							"$size",
							0,
						}}}}}},
					}}},
					bson.D{{Key: "$project", Value: bson.D{
						{Key: "_id", Value: nil},
						{Key: "dominance", Value: bson.D{{Key: "$round", Value: bson.A{
							bson.D{{Key: "$multiply", Value: bson.A{
								bson.D{{Key: "$cond", Value: bson.A{
									bson.D{{Key: "$eq", Value: bson.A{"$total_cock", 0}}},
									0,
									bson.D{{Key: "$divide", Value: bson.A{"$total_user_cock", "$total_cock"}}},
								}}},
								100,
							}}},
							1,
						}}}},
					}}},
				}},
				{Key: "individual_daily_growth", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: query.From.ID}}}},
					bson.D{{Key: "$sort", Value: bson.D{{Key: "requested_at", Value: -1}}}},
					bson.D{{Key: "$setWindowFields", Value: bson.D{
						{Key: "partitionBy", Value: "$user_id"},
						{Key: "sortBy", Value: bson.D{{Key: "requested_at", Value: -1}}},
						{Key: "output", Value: bson.D{
							{Key: "prev_size", Value: bson.D{{Key: "$shift", Value: bson.D{
								{Key: "output", Value: "$size"},
								{Key: "by", Value: 1},
							}}}},
						}},
					}}},
					bson.D{{Key: "$set", Value: bson.D{
						{Key: "growth", Value: bson.D{{Key: "$round", Value: bson.A{
							bson.D{{Key: "$subtract", Value: bson.A{"$size", "$prev_size"}}},
							1,
						}}}},
					}}},
					bson.D{{Key: "$group", Value: bson.D{
						{Key: "_id", Value: "$user_id"},
						{Key: "average_daily_growth", Value: bson.D{{Key: "$avg", Value: "$growth"}}},
					}}},
					bson.D{{Key: "$project", Value: bson.D{
						{Key: "_id", Value: nil},
						{Key: "average", Value: bson.D{{Key: "$round", Value: bson.A{"$average_daily_growth", 1}}}},
					}}},
				}},
				{Key: "individual_daily_dynamics", Value: bson.A{
					bson.D{{Key: "$match", Value: bson.D{{Key: "user_id", Value: query.From.ID}}}},
					bson.D{{Key: "$project", Value: bson.D{
						{Key: "requested_at", Value: 1},
						{Key: "size", Value: 1},
					}}},
					bson.D{{Key: "$sort", Value: bson.D{{Key: "requested_at", Value: -1}}}},
					bson.D{{Key: "$limit", Value: 2}},
					bson.D{{Key: "$group", Value: bson.D{
						{Key: "_id", Value: nil},
						{Key: "curr_cock", Value: bson.D{{Key: "$first", Value: "$size"}}},
						{Key: "prev_cock", Value: bson.D{{Key: "$last", Value: "$size"}}},
					}}},
					bson.D{{Key: "$project", Value: bson.D{
						{Key: "_id", Value: 0},
						{Key: "yesterday_cock_change", Value: bson.D{{Key: "$subtract", Value: bson.A{"$curr_cock", "$prev_cock"}}}},
						{Key: "yesterday_cock_change_percent", Value: bson.D{{Key: "$round", Value: bson.A{
							bson.D{{Key: "$multiply", Value: bson.A{
								bson.D{{Key: "$divide", Value: bson.A{
									bson.D{{Key: "$subtract", Value: bson.A{"$curr_cock", "$prev_cock"}}},
									"$prev_cock",
								}}},
								100,
							}}},
							1,
						}}}},
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
						{Key: "_id", Value: nil},
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
		IndividualCock []struct {
			Total   int `bson:"total"`
			Average int `bson:"average"`
		} `bson:"individual_cock"`

		IndividualIrk []struct {
			Irk float64 `bson:"irk"`
		}

		IndividualRecord []struct {
			RequestedAt time.Time `bson:"requested_at"`
			Total       int       `bson:"total"`
		} `bson:"individual_record"`

		IndividualDominance []struct {
			Dominance float64 `bson:"dominance"`
		} `bson:"individual_dominance"`

		IndividualDailyGrowth []struct {
			Average float64 `bson:"average"`
		} `bson:"individual_daily_growth"`

		IndividualDailyDynamics []struct {
			YesterdayCockChange        int     `bson:"yesterday_cock_change"`
			YesterdayCockChangePercent float64 `bson:"yesterday_cock_change_percent"`
		} `bson:"individual_daily_dynamics"`

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

	log.I("Aggregation completed successfully", "AggregationResult", result)

	individualCock := result.IndividualCock[0]
	individualRecord := result.IndividualRecord[0]
	individualIrk := result.IndividualIrk[0]
	individualDominance := result.IndividualDominance[0]
	individualDailyGrowth := result.IndividualDailyGrowth[0]
	individualDailyDynamics := result.IndividualDailyDynamics[0]

	overall := result.Overall[0]
	overallCockers := result.Uniques[0].Count
	overallDistribution := result.Distribution[0]
	overallRecord := result.Record[0]

	text := NewMsgCockDynamicsTemplate(
		/* Общая динамика коков */
		overall.Size,
		overallCockers,
		overall.Average,
		overall.Median,

		/* Персональная динамика кока */
		individualCock.Total,
		individualCock.Average,
		individualIrk.Irk,
		individualRecord.Total,
		individualRecord.RequestedAt,

		/* Кок-активы */
		individualDailyDynamics.YesterdayCockChangePercent,
		individualDailyDynamics.YesterdayCockChange,
		individualDailyGrowth.Average,

		/* Соотношение коков */
		overallDistribution.HugePercent,
		overallDistribution.LittlePercent,

		/* Самый большой кок */
		overallRecord.RequestedAt,
		overallRecord.Total,

		/* % доминирование */
		individualDominance.Dominance,
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
