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
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

// shouldShowDescription проверяет, нужно ли показывать описания для пользователя
// Описания НЕ показываются если: userCocksCount > 32 И username != "mairwunnx"
func (app *Application) shouldShowDescription(log *logging.Logger, userID int64, username string) bool {
	if username == "mairwunnx0" {
		return true
	}
	
	cocksCount := app.GetUserCocksCount(log, userID)
	
	// Если больше 32 коков, не показываем описания, очевидно юзер уже не новичок
	if cocksCount > 32 {
		return false
	}
	
	return true
}

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

	return InitializeInlineQueryWithThumb(
		"Размер кока",
		strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(text, ".", "\\."), "-", "\\-"), "!", "\\!"),
		"https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_size.png",
	)
}

func (app *Application) InlineQueryCockLadder(log *logging.Logger, query *tgbotapi.InlineQuery) tgbotapi.InlineQueryResultArticle {
	cocks := app.AggregateCockSizes(log)
	totalParticipants := app.GetTotalCockersCount(log)
	showDescription := app.shouldShowDescription(log, query.From.ID, query.From.UserName)
	text := app.GenerateCockLadderScoreboard(log, query.From.ID, cocks, totalParticipants, showDescription)
	return InitializeInlineQueryWithThumb(
		"Ладдер коков",
		text,
		"https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_ladder.png",
	)
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
	
	showDescription := app.shouldShowDescription(log, query.From.ID, query.From.UserName)
	text := app.GenerateCockRaceScoreboard(log, query.From.ID, cocks, seasonStartDate, totalParticipants, currentSeason, showDescription)
	return InitializeInlineQueryWithThumb(
		"Гонка коков",
		text,
		"https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_race.png",
	)
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

	// Проверяем, что у пользователя есть данные
	if len(result.IndividualCockTotal) == 0 || len(result.Overall) == 0 {
		log.E("User has no cock data yet")
		text := "🤔 *Недостаточно данных для динамики*\n\n_Сначала дёрни кок хотя бы раз\\!_"
		return InitializeInlineQueryWithThumb(
			"Динамика кока",
			text,
			"https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_dynamic.png",
		)
	}

	individualCockTotal := result.IndividualCockTotal[0]
	
	// Проверяем наличие данных по среднему коку (требует минимум 5 коков)
	var individualCockRecentAverage int
	if len(result.IndividualCockRecent) > 0 {
		individualCockRecentAverage = result.IndividualCockRecent[0].Average
	}
	
	// Проверяем наличие данных по рекорду
	var individualRecordTotal int
	var individualRecordDate time.Time
	if len(result.IndividualRecord) > 0 {
		individualRecordTotal = result.IndividualRecord[0].Total
		individualRecordDate = result.IndividualRecord[0].RequestedAt
	} else {
		// Если нет рекорда, используем данные из общего
		individualRecordTotal = individualCockTotal.Total
		individualRecordDate = datetime.NowTime()
	}
	
	individualIrk := result.IndividualIrk[0]
	individualDominance := result.IndividualDominance[0]
	
	// Получаем дату первого кока пользователя
	var userFirstCockDate time.Time
	var userPullingPeriod string
	if len(result.IndividualFirstCockDate) > 0 {
		userFirstCockDate = result.IndividualFirstCockDate[0].FirstDate
		userPullingPeriod = FormatUserPullingPeriod(userFirstCockDate, datetime.NowTime())
	} else {
		userPullingPeriod = "недавно"
	}
	
	// Проверяем наличие данных для дневной динамики (может отсутствовать у новых пользователей)
	var yesterdayCockChange int
	var yesterdayCockChangePercent float64
	if len(result.IndividualDailyDynamics) > 0 {
		yesterdayCockChange = result.IndividualDailyDynamics[0].YesterdayCockChange
		yesterdayCockChangePercent = result.IndividualDailyDynamics[0].YesterdayCockChangePercent
	}
	
	// Проверяем наличие данных для динамики за 5 коков (требует минимум 5 коков)
	var fiveCocksChange int
	var fiveCocksChangePercent float64
	if len(result.IndividualFiveCocksDynamics) > 0 {
		fiveCocksChange = result.IndividualFiveCocksDynamics[0].FiveCocksChange
		fiveCocksChangePercent = result.IndividualFiveCocksDynamics[0].FiveCocksChangePercent
	}
	
	// Проверяем наличие данных для скорости роста (требует минимум 5 коков)
	var growthSpeed float64
	if len(result.IndividualGrowthSpeed) > 0 {
		growthSpeed = result.IndividualGrowthSpeed[0].GrowthSpeed
	}

	overall := result.Overall[0]
	overallRecent := result.OverallRecent[0]
	overallCockers := result.Uniques[0].Count
	overallDistribution := result.Distribution[0]
	overallRecord := result.Record[0]
	
	totalCocksCount := result.TotalCocksCount[0].TotalCount
	
	// Получаем скорость роста общей статистики
	var overallGrowthSpeed float64
	if len(result.OverallGrowthSpeed) > 0 {
		overallGrowthSpeed = result.OverallGrowthSpeed[0].GrowthSpeed
	}
	
	// Проверяем наличие данных по количеству коков пользователя
	var userCocksCount int
	if len(result.IndividualCocksCount) > 0 {
		userCocksCount = result.IndividualCocksCount[0].UserCount
	}
	
	// Проверяем наличие данных для коэффициента везения (требует минимум 5 коков)
	var userLuckCoefficient float64
	if len(result.IndividualLuck) > 0 {
		userLuckCoefficient = result.IndividualLuck[0].LuckCoefficient
	} else {
		userLuckCoefficient = 1.0 // Нейтральное значение
	}
	
	// Проверяем наличие данных для волатильности (требует минимум 5 коков)
	var userVolatility float64
	if len(result.IndividualVolatility) > 0 {
		userVolatility = result.IndividualVolatility[0].Volatility
	}

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
		individualCockRecentAverage,
		individualIrk.Irk,
		individualRecordTotal,
		individualRecordDate,

		/* Кок-активы */
		yesterdayCockChangePercent,
		yesterdayCockChange,
		fiveCocksChangePercent,
		fiveCocksChange,

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
		growthSpeed,
		
		/* Скорость роста общей статистики */
		overallGrowthSpeed,
		
		/* Период дергания кока пользователем */
		userPullingPeriod,
	)

	article := tgbotapi.NewInlineQueryResultArticleMarkdown(query.ID, "Динамика кока", text)
	article.ThumbURL = "https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_dynamic.png"
	return article
}

func (app *Application) InlineQueryCockSeason(log *logging.Logger, query *tgbotapi.InlineQuery) tgbotapi.InlineQueryResultArticle {
	allSeasons := app.GetAllSeasonsForStats(log)
	
	if len(allSeasons) == 0 {
		text := NewMsgCockSeasonNoSeasonsTemplate()
		return InitializeInlineQueryWithThumb(
			"Сезоны коков",
			text,
			"https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_seasons.png",
		)
	}
	
	// Начинаем с последнего (самого нового) сезона
	currentSeasonIdx := len(allSeasons) - 1
	currentSeason := allSeasons[currentSeasonIdx]
	
	getSeasonWinners := func(season CockSeason) []SeasonWinner {
		return app.GetSeasonWinners(log, season)
	}
	
	showDescription := app.shouldShowDescription(log, query.From.ID, query.From.UserName)
	text := NewMsgCockSeasonSinglePage(currentSeason, getSeasonWinners, showDescription)
	
	// Создаем кнопки навигации
	var buttons []tgbotapi.InlineKeyboardButton
	
	// Кнопка "предыдущий сезон" (более старый, влево)
	if currentSeasonIdx > 0 {
		prevSeason := allSeasons[currentSeasonIdx-1]
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("◀️ Сезон %d", prevSeason.SeasonNum),
			fmt.Sprintf("season_page:%d", prevSeason.SeasonNum),
		))
	}
	
	// Кнопка "следующий сезон" (более новый, вправо) - только если есть более новый
	// (на самом деле, если мы на последнем сезоне, следующего нет)
	// Но для будущих сезонов это может быть полезно
	
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(buttons...),
	)
	
	article := tgbotapi.NewInlineQueryResultArticleMarkdownV2(
		uuid.NewString(),
		"Сезоны коков",
		text,
	)
	article.ReplyMarkup = &kb
	article.ThumbURL = "https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_seasons.png"
	
	return article
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

	showDescription := app.shouldShowDescription(log, query.From.ID, query.From.UserName)
	text := app.GenerateCockRulerText(log, query.From.ID, cocks, totalParticipants, showDescription)
	return InitializeInlineQueryWithThumb(
		"Линейка коков",
		text,
		"https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_ruler.png",
	)
}

func (app *Application) InlineQueryCockAchievements(log *logging.Logger, query *tgbotapi.InlineQuery, page int) tgbotapi.InlineQueryResultArticle {
	userID := query.From.ID
	
	// Проверка только для тестового пользователя
	// if userID != 362695653 {
	// 	text := "🔒 *Кок\\-ачивки временно доступны только для тестирования*\n\n_Скоро будут доступны для всех\\!_"
	// 	return InitializeInlineQueryWithThumb(
	// 		"Кок-ачивки",
	// 		text,
	// 		"https://files.mairwunnx.com/raw/public/dickobrazz%2FGemini_Generated_Image_qkh4tfqkh4tfqkh4.png",
	// 	)
	// }
	
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
	
	// Выбираем шаблон в зависимости от страницы
	var template string
	if page == 1 {
		template = MsgCockAchievementsTemplate
	} else {
		template = MsgCockAchievementsTemplateOtherPages
	}
	
	text := fmt.Sprintf(
		template,
		completedCount,
		totalAchievements,
		percentComplete,
		totalRespects,
		achievementsList,
	)
	
	// Создаем кнопки пагинации (с userID владельца)
	var buttons []tgbotapi.InlineKeyboardButton
	
	if page > 1 {
		// Кнопка "предыдущая страница"
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("◀️", fmt.Sprintf("ach_page:%d:%d", userID, page-1)))
	}
	
	// Кнопка "текущая страница / всего страниц"
	buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d/%d", page, totalPages), "ach_noop"))
	
	if page < totalPages {
		// Кнопка "следующая страница"
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("▶️", fmt.Sprintf("ach_page:%d:%d", userID, page+1)))
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
	article.ThumbURL = "https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_achievements.png"
	
	return article
}

func InitializeInlineQuery(title, message string) tgbotapi.InlineQueryResultArticle {
	return tgbotapi.NewInlineQueryResultArticleMarkdownV2(uuid.NewString(), title, message)
}

func InitializeInlineQueryWithThumb(title, message, thumbURL string) tgbotapi.InlineQueryResultArticle {
	article := tgbotapi.NewInlineQueryResultArticleMarkdownV2(uuid.NewString(), title, message)
	article.ThumbURL = thumbURL
	return article
}

func (app *Application) HandleCallbackQuery(log *logging.Logger, callback *tgbotapi.CallbackQuery) {
	// Парсим callback data
	data := callback.Data
	
	// Обрабатываем пагинацию сезонов
	if strings.HasPrefix(data, "season_page:") {
		// Парсим номер сезона
		seasonNumStr := strings.TrimPrefix(data, "season_page:")
		seasonNum := 1
		if parsedSeasonNum, err := strconv.Atoi(seasonNumStr); err != nil {
			log.E("Failed to parse season number", logging.InnerError, err)
			seasonNum = 1
		} else {
			seasonNum = parsedSeasonNum
		}
		
		// Получаем все сезоны
		allSeasons := app.GetAllSeasonsForStats(log)
		
		// Находим нужный сезон
		var targetSeason *CockSeason
		var targetIdx int
		for idx, s := range allSeasons {
			if s.SeasonNum == seasonNum {
				targetSeason = &s
				targetIdx = idx
				break
			}
		}
		
		if targetSeason == nil {
			log.E("Season not found", "season_num", seasonNum)
			callbackConfig := tgbotapi.NewCallback(callback.ID, "Сезон не найден")
			if _, err := app.bot.Request(callbackConfig); err != nil {
				log.E("Failed to answer callback query", logging.InnerError, err)
			}
			return
		}
		
		getSeasonWinners := func(season CockSeason) []SeasonWinner {
			return app.GetSeasonWinners(log, season)
		}
		
		showDescription := app.shouldShowDescription(log, callback.From.ID, callback.From.UserName)
		text := NewMsgCockSeasonSinglePage(*targetSeason, getSeasonWinners, showDescription)
		
		// Создаем кнопки навигации
		var buttons []tgbotapi.InlineKeyboardButton
		
		// Кнопка "предыдущий сезон" (более старый, влево)
		if targetIdx > 0 {
			prevSeason := allSeasons[targetIdx-1]
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("◀️ Сезон %d", prevSeason.SeasonNum),
				fmt.Sprintf("season_page:%d", prevSeason.SeasonNum),
			))
		}
		
		// Кнопка "следующий сезон" (более новый, вправо)
		if targetIdx < len(allSeasons)-1 {
			nextSeason := allSeasons[targetIdx+1]
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
				fmt.Sprintf("▶️ Сезон %d", nextSeason.SeasonNum),
				fmt.Sprintf("season_page:%d", nextSeason.SeasonNum),
			))
		}
		
		kb := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(buttons...),
		)
		
		// Отвечаем на callback
		_, _ = app.bot.Request(tgbotapi.NewCallback(callback.ID, ""))
		
		// Редактируем существующее сообщение
		if callback.InlineMessageID != "" {
			edit := tgbotapi.EditMessageTextConfig{
				BaseEdit: tgbotapi.BaseEdit{
					InlineMessageID: callback.InlineMessageID,
					ReplyMarkup:     &kb,
				},
				Text:      text,
				ParseMode: "MarkdownV2",
			}
			if _, err := app.bot.Request(edit); err != nil {
				log.E("Failed to edit inline message", logging.InnerError, err)
			} else {
				log.I("Successfully edited inline message", "season_num", seasonNum)
			}
		} else if callback.Message != nil {
			edit := tgbotapi.NewEditMessageTextAndMarkup(
				callback.Message.Chat.ID,
				callback.Message.MessageID,
				text,
				kb,
			)
			edit.ParseMode = "MarkdownV2"
			
			if _, err := app.bot.Request(edit); err != nil {
				log.E("Failed to edit chat message", logging.InnerError, err)
			} else {
				log.I("Successfully edited chat message", "season_num", seasonNum)
			}
		} else {
			log.E("CallbackQuery has neither Message nor InlineMessageID")
		}
	} else if strings.HasPrefix(data, "ach_page:") {
		// Парсим userID и номер страницы из формата "ach_page:userID:page"
		parts := strings.Split(strings.TrimPrefix(data, "ach_page:"), ":")
		if len(parts) != 2 {
			log.E("Invalid ach_page callback data format", "data", data)
			callbackConfig := tgbotapi.NewCallback(callback.ID, "Неверный формат данных")
			if _, err := app.bot.Request(callbackConfig); err != nil {
				log.E("Failed to answer callback query", logging.InnerError, err)
			}
			return
		}
		
		userID, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			log.E("Failed to parse userID from callback", logging.InnerError, err)
			callbackConfig := tgbotapi.NewCallback(callback.ID, "Ошибка парсинга данных")
			if _, err := app.bot.Request(callbackConfig); err != nil {
				log.E("Failed to answer callback query", logging.InnerError, err)
			}
			return
		}
		
		page := 1
		if parsedPage, err := strconv.Atoi(parts[1]); err != nil {
			log.E("Failed to parse page number", logging.InnerError, err)
			page = 1
		} else {
			page = parsedPage
		}
		// if userID != 362695653 {
		// 	// Отвечаем на callback и выходим
		// 	callbackConfig := tgbotapi.NewCallback(callback.ID, "Ачивки доступны только для тестирования")
		// 	if _, err := app.bot.Request(callbackConfig); err != nil {
		// 		log.E("Failed to answer callback query", logging.InnerError, err)
		// 	}
		// 	return
		// }
		
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
		
		// Выбираем шаблон в зависимости от страницы
		var template string
		if page == 1 {
			template = MsgCockAchievementsTemplate
		} else {
			template = MsgCockAchievementsTemplateOtherPages
		}
		
		text := fmt.Sprintf(
			template,
			completedCount,
			totalAchievements,
			percentComplete,
			totalRespects,
			achievementsList,
		)
		
		// Создаем кнопки пагинации для новой страницы (с userID владельца)
		var buttons []tgbotapi.InlineKeyboardButton
		
		if page > 1 {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("◀️", fmt.Sprintf("ach_page:%d:%d", userID, page-1)))
		}
		
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d/%d", page, totalPages), "ach_noop"))
		
		if page < totalPages {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("▶️", fmt.Sprintf("ach_page:%d:%d", userID, page+1)))
		}
		
		kb := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(buttons...),
		)
		
		// Отвечаем на callback (убираем "часики" на кнопке)
		_, _ = app.bot.Request(tgbotapi.NewCallback(callback.ID, ""))
		
		// Редактируем существующее сообщение
		if callback.InlineMessageID != "" {
			// INLINE message: редактируем по InlineMessageID
			edit := tgbotapi.EditMessageTextConfig{
				BaseEdit: tgbotapi.BaseEdit{
					InlineMessageID: callback.InlineMessageID,
					ReplyMarkup:     &kb,
				},
				Text:      text,
				ParseMode: "MarkdownV2",
			}
			if _, err := app.bot.Request(edit); err != nil {
				log.E("Failed to edit inline message", logging.InnerError, err)
			} else {
				log.I("Successfully edited inline message", "page", page)
			}
		} else if callback.Message != nil {
			// Обычное сообщение в чате: редактируем по chat_id/message_id
			edit := tgbotapi.NewEditMessageTextAndMarkup(
				callback.Message.Chat.ID,
				callback.Message.MessageID,
				text,
				kb,
			)
			edit.ParseMode = "MarkdownV2"
			
			if _, err := app.bot.Request(edit); err != nil {
				log.E("Failed to edit chat message", logging.InnerError, err)
			} else {
				log.I("Successfully edited chat message", "page", page)
			}
		} else {
			// Крайний случай — некуда редактировать
			log.E("CallbackQuery has neither Message nor InlineMessageID")
		}
	} else if data == "ach_noop" {
		// Просто отвечаем на callback (для кнопки с текущей страницей)
		callbackConfig := tgbotapi.NewCallback(callback.ID, "")
		if _, err := app.bot.Request(callbackConfig); err != nil {
			log.E("Failed to answer callback query", logging.InnerError, err)
		}
	}
}
