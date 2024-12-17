package application

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"math/rand"
	"sort"
	"strings"
)

func (app *Application) HandleInlineQuery(log *Logger, query *tgbotapi.InlineQuery) {
	results := []any{
		app.InlineQueryCockSize(log, query),
		app.InlineQueryCockRace(log, query),
		app.InlineQueryCockRuler(log, query),
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
		size = rand.Intn(61)

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
