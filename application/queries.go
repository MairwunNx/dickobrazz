package application

import (
	"context"
	"dickobrazz/application/database"
	"dickobrazz/application/datetime"
	"dickobrazz/application/geo"
	"dickobrazz/application/logging"
	"dickobrazz/application/timings"
	"fmt"
	"sort"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

func (app *Application) HandleInlineQuery(log *logging.Logger, query *tgbotapi.InlineQuery) {
	var traceQueryCreated = func(l *logging.Logger) { l.I("Inline query successfully created") }

	queries := []any{
		timings.ReportExecutionForResult(log.With(logging.QueryType, "CockSize"),
			func() tgbotapi.InlineQueryResultArticle { return app.InlineQueryCockSize(log, query) }, traceQueryCreated,
		),
		timings.ReportExecutionForResult(log.With(logging.QueryType, "CockRace"),
			func() tgbotapi.InlineQueryResultArticle { return app.InlineQueryCockRace(log, query) }, traceQueryCreated,
		),
		timings.ReportExecutionForResult(log.With(logging.QueryType, "CockRuler"),
			func() tgbotapi.InlineQueryResultArticle { return app.InlineQueryCockRuler(log, query) }, traceQueryCreated,
		),
		timings.ReportExecutionForResult(log.With(logging.QueryType, "CockLadder"),
			func() tgbotapi.InlineQueryResultArticle { return app.InlineQueryCockLadder(log, query) }, traceQueryCreated,
		),
		timings.ReportExecutionForResult(log.With(logging.QueryType, "CockDynamic"),
			func() tgbotapi.InlineQueryResultArticle { return app.InlineQueryCockDynamic(log, query) }, traceQueryCreated,
		),
		timings.ReportExecutionForResult(log.With(logging.QueryType, "CockSeason"),
			func() tgbotapi.InlineQueryResultArticle { return app.InlineQueryCockSeason(log, query) }, traceQueryCreated,
		),
		timings.ReportExecutionForResult(log.With(logging.QueryType, "CockAchievements"),
			func() tgbotapi.InlineQueryResultArticle { 
				// Парсим номер страницы из query (если есть)
				page := 1
				// По умолчанию страница 1, можно расширить парсинг в будущем
				return app.InlineQueryCockAchievements(log, query, page) 
			}, traceQueryCreated,
		),
	}

	for i, q := range queries {
		if article, ok := q.(tgbotapi.InlineQueryResultArticle); ok {
			log.I("Inline query text preview", "index", i, "title", article.Title, "text_length", len(article.InputMessageContent.(tgbotapi.InputTextMessageContent).Text), "text", article.InputMessageContent.(tgbotapi.InputTextMessageContent).Text)
		}
	}

	inlines := tgbotapi.InlineConfig{InlineQueryID: query.ID, IsPersonal: true, CacheTime: 60, Results: queries}

	if _, err := timings.ReportExecutionForResultError(log,
		func() (*tgbotapi.APIResponse, error) { return app.bot.Request(inlines) },
		func(l *logging.Logger) { l.I("Inline query successfully sent") },
	); err != nil {
		log.E("Failed to send inline query", logging.InnerError, err)
	}
}

func (app *Application) InlineQueryCockSize(log *logging.Logger, query *tgbotapi.InlineQuery) tgbotapi.InlineQueryResultArticle {
	var size int

	if cached := app.GetCockSizeFromCache(log, query.From.ID); cached != nil {
		size = *cached
	} else {
		size = app.rnd.IntN(log, 60)

		// Нормализуем username (генерируем анонимное имя если пустой)
		normalizedUsername := NormalizeUsername(query.From.UserName, query.From.ID)

		cock := &Cock{
			ID:          uuid.NewString(),
			Size:        int32(size),
			Nickname:    normalizedUsername,
			UserID:      query.From.ID,
			RequestedAt: datetime.NowTime(),
		}

		app.SaveCockToCache(log, query.From.ID, normalizedUsername, size)
		app.SaveCockToMongo(log, cock)
	}

	emoji := EmojiFromSize(size)
	text := GenerateCockSizeText(size, emoji)
	subtext := geo.GetRegionBySize(size)

	text = text + "\n\n" + "_" + subtext + "_"

	return InitializeInlineQuery(
		"Размер кока",
		strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(text, ".", "\\."), "-", "\\-"), "!", "\\!"),
	)
}

func (app *Application) InlineQueryCockLadder(log *logging.Logger, query *tgbotapi.InlineQuery) tgbotapi.InlineQueryResultArticle {
	cocks := app.AggregateCockSizes(log)
	totalParticipants := app.GetTotalCockersCount(log)
	text := app.GenerateCockLadderScoreboard(log, query.From.ID, cocks, totalParticipants)
	return InitializeInlineQuery("Ладдер коков", text)
}

func (app *Application) InlineQueryCockRace(log *logging.Logger, query *tgbotapi.InlineQuery) tgbotapi.InlineQueryResultArticle {
	currentSeason := app.GetCurrentSeason(log)
	
	var cocks []UserCockRace
	var seasonStartDate string
	var totalParticipants int
	
	if currentSeason != nil {
		cocks = app.AggregateCockSizesForSeason(log, *currentSeason)
		totalParticipants = app.GetSeasonCockersCount(log, *currentSeason)
		seasonStartDate = EscapeMarkdownV2(currentSeason.StartDate.Format("02.01.2006"))
	} else {
		cocks = app.AggregateCockSizes(log)
		totalParticipants = app.GetTotalCockersCount(log)
		seasonStartDate = "хуй знает когда" // Заглушка для случая если нет активного сезона (чего в целом быть не может, я в это верю.)
	}
	
	text := app.GenerateCockRaceScoreboard(log, query.From.ID, cocks, seasonStartDate, totalParticipants, currentSeason)
	return InitializeInlineQuery("Гонка коков", text)
}

func (app *Application) InlineQueryCockDynamic(log *logging.Logger, query *tgbotapi.InlineQuery) tgbotapi.InlineQueryResultArticle {
	collection := timings.ReportExecutionForResult(log,
		func() *mongo.Collection { return database.CollectionCocks(app.db) },
		func(l *logging.Logger) { l.I("Collection successfully fetched") },
	)

	pipeline := timings.ReportExecutionForResult(log,
		func() mongo.Pipeline { return database.PipelineDynamic(query.From.ID) },
		func(l *logging.Logger) { l.I("Cock dynamic pipeline has successfully built") },
	)

	cursor, err := timings.ReportExecutionForResultError(log,
		func() (*mongo.Cursor, error) {
			return collection.Aggregate(app.ctx, pipeline)
		},
		func(l *logging.Logger) { l.I("Cock dynamic pipeline has successfully aggregated") },
	)

	if err != nil {
		log.E("Aggregation failed", logging.InnerError, err)
		return tgbotapi.InlineQueryResultArticle{}
	}

	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			log.E("Failed to close mongo cursor", logging.InnerError, err)
		}
	}(cursor, app.ctx)

	var result *database.DocumentCockDynamic

	if err := timings.ReportExecutionForResult(log,
		func() error { cursor.Next(app.ctx); return cursor.Decode(&result) },
		func(l *logging.Logger) { l.I("Cock dynamic pipeline has successfully decoded") },
	); err != nil {
		log.E("Failed to decode aggregation results", logging.InnerError, err)
		return tgbotapi.InlineQueryResultArticle{}
	}

	individualCockTotal := result.IndividualCockTotal[0]
	individualCockRecent := result.IndividualCockRecent[0]
	individualRecord := result.IndividualRecord[0]
	individualIrk := result.IndividualIrk[0]
	individualDominance := result.IndividualDominance[0]
	individualDailyDynamics := result.IndividualDailyDynamics[0]
	individualFiveCocksDynamics := result.IndividualFiveCocksDynamics[0]
	individualGrowthSpeed := result.IndividualGrowthSpeed[0]

	overall := result.Overall[0]
	overallRecent := result.OverallRecent[0]
	overallCockers := result.Uniques[0].Count
	overallDistribution := result.Distribution[0]
	overallRecord := result.Record[0]
	
	totalCocksCount := result.TotalCocksCount[0].TotalCount
	userCocksCount := result.IndividualCocksCount[0].UserCount
	
	userLuckCoefficient := result.IndividualLuck[0].LuckCoefficient
	userVolatility := result.IndividualVolatility[0].Volatility

	userSeasonWins := app.GetUserSeasonWins(log, query.From.ID)
	userCockRespect := app.GetUserCockRespect(log, query.From.ID)

	text := NewMsgCockDynamicsTemplate(
		/* Общая динамика коков */
		overall.Size,
		overallCockers,
		overallRecent.Average,
		overallRecent.Median,

		/* Персональная динамика кока */
		individualCockTotal.Total,
		individualCockRecent.Average,
		individualIrk.Irk,
		individualRecord.Total,
		individualRecord.RequestedAt,

		/* Кок-активы */
		individualDailyDynamics.YesterdayCockChangePercent,
		individualDailyDynamics.YesterdayCockChange,
		individualFiveCocksDynamics.FiveCocksChangePercent,
		individualFiveCocksDynamics.FiveCocksChange,

		/* Соотношение коков */
		overallDistribution.HugePercent,
		overallDistribution.LittlePercent,

		/* Самый большой кок */
		overallRecord.RequestedAt,
		overallRecord.Total,

		/* % доминирование */
		individualDominance.Dominance,

		/* Сезонные достижения */
		userSeasonWins,
		userCockRespect,

		/* Всего дёрнуто коков */
		totalCocksCount,
		userCocksCount,

		/* Коэффициент везения и волатильность */
		userLuckCoefficient,
		userVolatility,
		
		/* Средняя скорость прироста */
		individualGrowthSpeed.GrowthSpeed,
	)

	return tgbotapi.NewInlineQueryResultArticleMarkdown(query.ID, "Динамика кока", text)
}

func (app *Application) InlineQueryCockSeason(log *logging.Logger, query *tgbotapi.InlineQuery) tgbotapi.InlineQueryResultArticle {
	seasons := app.GetAllSeasons(log)
	totalSeasonsCount := app.GetAllSeasonsCount(log)
	
	getSeasonWinners := func(season CockSeason) []SeasonWinner {
		return app.GetSeasonWinners(log, season)
	}
	
	text := NewMsgCockSeasonsFullText(seasons, totalSeasonsCount, getSeasonWinners)
	return InitializeInlineQuery("Сезоны коков", text)
}

func (app *Application) InlineQueryCockRuler(log *logging.Logger, query *tgbotapi.InlineQuery) tgbotapi.InlineQueryResultArticle {
	cocks := app.GetCockSizesFromCache(log)
	totalParticipants := len(cocks)

	sort.Slice(cocks, func(i, j int) bool {
		return cocks[i].Size > cocks[j].Size
	})

	if len(cocks) > 13 {
		cocks = cocks[:13]
	}

	text := app.GenerateCockRulerText(log, query.From.ID, cocks, totalParticipants)
	return InitializeInlineQuery("Линейка коков", text)
}

func (app *Application) InlineQueryCockAchievements(log *logging.Logger, query *tgbotapi.InlineQuery, page int) tgbotapi.InlineQueryResultArticle {
	userID := query.From.ID
	
	// Проверка только для тестового пользователя
	if userID != 362695653 {
		text := "🔒 *Кок\\-ачивки временно доступны только для тестирования*\n\n_Скоро будут доступны для всех\\!_"
		return InitializeInlineQuery("Кок-ачивки", text)
	}
	
	// Проверяем и обновляем достижения (только для mairwunnx, раз в сутки)
	app.CheckAndUpdateAchievements(log, userID)
	
	// Получаем достижения пользователя
	userAchievements := app.GetUserAchievements(log, userID)
	
	// Генерируем текст с пагинацией (10 ачивок на страницу)
	achievementsList, completedCount, totalRespects, percentComplete := GenerateAchievementsText(
		AllAchievements,
		userAchievements,
		page,
		10,
	)
	
	totalAchievements := len(AllAchievements)
	totalPages := (totalAchievements + 9) / 10
	
	text := fmt.Sprintf(
		MsgCockAchievementsTemplate,
		completedCount,
		totalAchievements,
		percentComplete,
		totalRespects,
		achievementsList,
	)
	
	// Создаем кнопки пагинации
	var buttons []tgbotapi.InlineKeyboardButton
	
	if page > 1 {
		// Кнопка "предыдущая страница"
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("◀️", fmt.Sprintf("ach_page:%d", page-1)))
	}
	
	// Кнопка "текущая страница / всего страниц"
	buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d/%d", page, totalPages), "ach_noop"))
	
	if page < totalPages {
		// Кнопка "следующая страница"
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("▶️", fmt.Sprintf("ach_page:%d", page+1)))
	}
	
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttons...),
	)
	
	article := tgbotapi.NewInlineQueryResultArticleMarkdownV2(
		uuid.NewString(),
		"Кок-ачивки",
		text,
	)
	article.ReplyMarkup = &kb
	
	return article
}

func InitializeInlineQuery(title, message string) tgbotapi.InlineQueryResultArticle {
	return tgbotapi.NewInlineQueryResultArticleMarkdownV2(uuid.NewString(), title, message)
}

func (app *Application) HandleCallbackQuery(log *logging.Logger, callback *tgbotapi.CallbackQuery) {
	// Парсим callback data
	data := callback.Data
	
	// Обрабатываем пагинацию ачивок
	if strings.HasPrefix(data, "ach_page:") {
		// Парсим номер страницы
		pageStr := strings.TrimPrefix(data, "ach_page:")
		page := 1
		if parsedPage, err := strconv.Atoi(pageStr); err != nil {
			log.E("Failed to parse page number", logging.InnerError, err)
			page = 1
		} else {
			page = parsedPage
		}
		
		// Проверяем, что callback от тестового пользователя
		userID := callback.From.ID
		if userID != 362695653 {
			// Отвечаем на callback и выходим
			callbackConfig := tgbotapi.NewCallback(callback.ID, "Ачивки доступны только для тестирования")
			if _, err := app.bot.Request(callbackConfig); err != nil {
				log.E("Failed to answer callback query", logging.InnerError, err)
			}
			return
		}
		
		// Проверяем и обновляем достижения (раз в сутки)
		app.CheckAndUpdateAchievements(log, userID)
		
		// Получаем достижения пользователя
		userAchievements := app.GetUserAchievements(log, userID)
		
		// Генерируем текст для запрошенной страницы
		achievementsList, completedCount, totalRespects, percentComplete := GenerateAchievementsText(
			AllAchievements,
			userAchievements,
			page,
			10,
		)
		
		totalAchievements := len(AllAchievements)
		totalPages := (totalAchievements + 9) / 10
		
		text := fmt.Sprintf(
			MsgCockAchievementsTemplate,
			completedCount,
			totalAchievements,
			percentComplete,
			totalRespects,
			achievementsList,
		)
		
		// Создаем кнопки пагинации для новой страницы
		var buttons []tgbotapi.InlineKeyboardButton
		
		if page > 1 {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("◀️", fmt.Sprintf("ach_page:%d", page-1)))
		}
		
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d/%d", page, totalPages), "ach_noop"))
		
		if page < totalPages {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("▶️", fmt.Sprintf("ach_page:%d", page+1)))
		}
		
		kb := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(buttons...),
		)
		
		// Редактируем существующее сообщение
		editMsg := tgbotapi.NewEditMessageTextAndMarkup(
			callback.Message.Chat.ID,
			callback.Message.MessageID,
			text,
			kb,
		)
		editMsg.ParseMode = "MarkdownV2"
		
		if _, err := app.bot.Send(editMsg); err != nil {
			log.E("Failed to edit message", logging.InnerError, err)
		}
		
		// Отвечаем на callback (убираем "часики" на кнопке)
		callbackConfig := tgbotapi.NewCallback(callback.ID, "")
		if _, err := app.bot.Request(callbackConfig); err != nil {
			log.E("Failed to answer callback query", logging.InnerError, err)
		}
		
		log.I("Successfully handled achievements pagination callback", "page", page)
	} else if data == "ach_noop" {
		// Просто отвечаем на callback (для кнопки с текущей страницей)
		callbackConfig := tgbotapi.NewCallback(callback.ID, "")
		if _, err := app.bot.Request(callbackConfig); err != nil {
			log.E("Failed to answer callback query", logging.InnerError, err)
		}
	}
}
