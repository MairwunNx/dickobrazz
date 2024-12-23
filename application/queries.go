package application

import (
	"crypto/rand"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math"
	"math/big"
	"os"
	"sort"
	"strings"
	"time"
)

func (app *Application) HandleInlineQuery(log *Logger, query *tgbotapi.InlineQuery) {
	var results []any
	if query.From.UserName == "mairwunnx" {
		results = []any{
			app.InlineQueryCockSize(log, query),
			app.InlineQueryCockRace(log, query),
			app.InlineQueryCockRaceImgStat(log, query),
			app.InlineQueryCockRuler(log, query),
		}
	} else {
		results = []any{
			app.InlineQueryCockSize(log, query),
			app.InlineQueryCockRace(log, query),
			app.InlineQueryCockRuler(log, query),
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
		random, _ := rand.Int(rand.Reader, big.NewInt(int64(61)))
		size = int(random.Int64())

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

type defaultColorPalette struct{}

func (dp defaultColorPalette) BackgroundColor() drawing.Color {
	return drawing.ColorFromHex("232323")
}

func (dp defaultColorPalette) BackgroundStrokeColor() drawing.Color {
	return drawing.ColorFromHex("00ff00")
}

func (dp defaultColorPalette) CanvasColor() drawing.Color {
	return drawing.ColorFromHex("232323")
}

func (dp defaultColorPalette) CanvasStrokeColor() drawing.Color {
	return drawing.ColorFromHex("ff0000")
}

func (dp defaultColorPalette) AxisStrokeColor() drawing.Color {
	return drawing.ColorFromHex("00ff00")
}

func (dp defaultColorPalette) TextColor() drawing.Color {
	return drawing.ColorFromHex("e9e9e9")
}

func (dp defaultColorPalette) GetSeriesColor(index int) drawing.Color {
	return drawing.ColorFromHex("c6c6c6")
}

func (app *Application) InlineQueryCockRaceImgStat(log *Logger, query *tgbotapi.InlineQuery) tgbotapi.InlineQueryResultPhoto {
	collection := app.db.Database("dickbot_db").Collection("cocks")

	pipeline := mongo.Pipeline{
		{
			{"$match", bson.D{{"user_id", query.From.ID}}},
		},
		{
			{"$group", bson.D{
				{"_id", "$requested_at"},
				{"size", bson.D{{"$avg", "$size"}}},
			}},
		},
		{
			{"$sort", bson.D{{"_id", 1}}},
		},
	}

	cursor, err := collection.Aggregate(app.ctx, pipeline)
	if err != nil {
		log.E("Failed to aggregate cock sizes", InnerError, err)
		return tgbotapi.InlineQueryResultPhoto{}
	}

	var results []struct {
		ID   time.Time `bson:"_id"`
		Size int32     `bson:"size"`
	}
	if err := cursor.All(app.ctx, &results); err != nil {
		log.E("Failed to decode cock sizes", InnerError, err)
		return tgbotapi.InlineQueryResultPhoto{}
	}

	log.I("Successfully aggregated cock sizes")

	if len(results) == 0 {
		log.I("No cock sizes found for user")
		return tgbotapi.InlineQueryResultPhoto{}
	}

	startDate := results[0].ID
	endDate := results[len(results)-1].ID
	startDateFormatted := startDate.Format("02.01.06")
	endDateFormatted := endDate.Format("02.01.06")

	// Interpolating to 14 points
	timestamps := make([]float64, len(results))
	sizes := make([]float64, len(results))
	for i, result := range results {
		timestamps[i] = float64(result.ID.Unix())
		sizes[i] = float64(result.Size)
	}

	interpolatedX, interpolatedY := interpolatePoints(timestamps, sizes, 14)
	annotations := spaceAroundAnnotations(interpolatedX, interpolatedY, 3)

	graph := chart.Chart{
		Title:      fmt.Sprintf("Развитие размера кока, текущий мой кок: %s", FormatDickSize(int32(sizes[len(sizes)-1]))),
		TitleStyle: chart.Style{FontSize: 14.0, Show: true},
		XAxis: chart.XAxis{
			GridMajorStyle: chart.Style{
				Show:        true,
				StrokeWidth: 0.5,
				StrokeColor: drawing.ColorFromHex("7b7b7b"),
			},
			GridMinorStyle: chart.Style{
				Show:        true,
				StrokeWidth: 0.25,
				StrokeColor: drawing.ColorFromHex("9b9b9b"),
			},
			Name:         fmt.Sprintf("Статистика моих коков с %s по %s", startDateFormatted, endDateFormatted),
			NameStyle:    chart.Style{Show: true},
			Style:        chart.Style{Show: true},
			TickPosition: chart.TickPositionBetweenTicks,
			Ticks: func() []chart.Tick {
				var ticks []chart.Tick
				for _, t := range results {
					ticks = append(ticks, chart.Tick{
						Value: float64(t.ID.Unix()),
						Label: t.ID.Format("02.01.06"),
					})
				}
				return ticks
			}(),
		},
		YAxis: chart.YAxis{
			GridMajorStyle: chart.Style{
				Show:        true,
				StrokeWidth: 0.5,
				StrokeColor: drawing.ColorFromHex("7b7b7b"),
			},
			GridMinorStyle: chart.Style{
				Show:        true,
				StrokeWidth: 0.25,
				StrokeColor: drawing.ColorFromHex("9b9b9b"),
			},
			Name:      "Размер (см)",
			NameStyle: chart.Style{Show: true},
			Style:     chart.Style{Show: true},
			ValueFormatter: func(v interface{}) string {
				value := v.(float64)
				return FormatDickSize(int32(value))
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "Размер кока",
				XValues: interpolatedX,
				YValues: interpolatedY,
				Style: chart.Style{
					Show:        true,
					StrokeWidth: 2.0,
				},
			},
			chart.AnnotationSeries{
				Annotations: annotations,
			},
		},
		ColorPalette: defaultColorPalette{},
	}

	filePath := fmt.Sprintf("%s_graph.png", query.ID)
	outputFile, err := os.Create(filePath)
	if err != nil {
		log.E("Failed to create graph image file", InnerError, err)
		return tgbotapi.InlineQueryResultPhoto{}
	}
	defer outputFile.Close()

	if err := graph.Render(chart.PNG, outputFile); err != nil {
		log.E("Failed to render graph image", InnerError, err)
		return tgbotapi.InlineQueryResultPhoto{}
	}

	log.I("Successfully created graph image")

	return tgbotapi.NewInlineQueryResultPhoto(
		query.ID,
		fmt.Sprintf("file://%s", filePath),
		//fmt.Sprintf("Развитие размера кока"),
	)
}

func interpolatePoints(x, y []float64, resolution int) ([]float64, []float64) {
	if len(x) != len(y) || len(x) < 2 {
		return x, y
	}

	interpolatedX := make([]float64, resolution)
	interpolatedY := make([]float64, resolution)

	xMin, xMax := x[0], x[len(x)-1]
	step := (xMax - xMin) / float64(resolution-1)

	for i := 0; i < resolution; i++ {
		interpolatedX[i] = xMin + float64(i)*step
		interpolatedY[i] = interpolateY(interpolatedX[i], x, y)
	}

	return interpolatedX, interpolatedY
}

func spaceAroundAnnotations(xValues, yValues []float64, count int) []chart.Value2 {
	n := len(xValues)
	if n < count {
		// Если точек меньше числа аннотаций, возвращаем все точки
		annotations := make([]chart.Value2, n)
		for i := range xValues {
			annotations[i] = chart.Value2{
				XValue: xValues[i],
				YValue: yValues[i],
				Label:  fmt.Sprintf("%.1f см", yValues[i]),
			}
		}
		return annotations
	}

	// Рассчитываем равномерные отступы (space around)
	step := float64(n-1) / float64(count+1)
	annotations := make([]chart.Value2, count)

	for i := 0; i < count; i++ {
		index := int(math.Round(step * float64(i+1))) // Распределяем равномерно между первой и последней точкой
		annotations[i] = chart.Value2{
			XValue: xValues[index],
			YValue: yValues[index],
			Label:  fmt.Sprintf("%.1f см", yValues[index]),
		}
	}

	return annotations
}

func interpolateY(x float64, xPoints, yPoints []float64) float64 {
	for i := 1; i < len(xPoints); i++ {
		if x <= xPoints[i] {
			x1, x2 := xPoints[i-1], xPoints[i]
			y1, y2 := yPoints[i-1], yPoints[i]
			return y1 + (y2-y1)*(x-x1)/(x2-x1)
		}
	}
	return math.NaN()
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
