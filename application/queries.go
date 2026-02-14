package application

import (
	"dickobrazz/application/api"
	"dickobrazz/application/datetime"
	"dickobrazz/application/geo"
	"dickobrazz/application/localization"
	"dickobrazz/application/logging"
	"dickobrazz/application/timings"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// shouldShowDescription проверяет, нужно ли показывать описания для пользователя.
// Описания НЕ показываются если пользователь зарегистрирован более 32 дней назад.
func (app *Application) shouldShowDescription(log *logging.Logger, userID int64, username string) bool {
	profile, err := app.api.GetProfile(app.ctx, userID, username)
	if err != nil {
		log.E("Failed to get profile for shouldShowDescription", logging.InnerError, err)
		return true
	}

	if profile.CreatedAt == nil {
		return true
	}

	createdAt, err := datetime.ParseUTC(*profile.CreatedAt)
	if err != nil {
		return true
	}

	if time.Since(createdAt).Hours()/24 > 32 {
		return false
	}

	return true
}

func (app *Application) HandleInlineQuery(log *logging.Logger, update *tgbotapi.Update) {
	query := update.InlineQuery
	if query == nil {
		return
	}
	var traceQueryCreated = func(l *logging.Logger) { l.I("Inline query successfully created") }

	// Синхронное выполнение CockSize (первым, так как может создавать данные)
	cockSizeResult := timings.ReportExecutionForResult(log.With(logging.QueryType, "CockSize"),
		func() tgbotapi.InlineQueryResultArticle { return app.InlineQueryCockSize(log, update) }, traceQueryCreated,
	)

	type queryResult struct {
		index  int
		result tgbotapi.InlineQueryResultArticle
	}

	parallelQueriesCount := 6
	resultsChan := make(chan queryResult, parallelQueriesCount)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		result := timings.ReportExecutionForResult(log.With(logging.QueryType, "CockRace"),
			func() tgbotapi.InlineQueryResultArticle { return app.InlineQueryCockRace(log, update) }, traceQueryCreated,
		)
		resultsChan <- queryResult{index: 1, result: result}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		result := timings.ReportExecutionForResult(log.With(logging.QueryType, "CockRuler"),
			func() tgbotapi.InlineQueryResultArticle { return app.InlineQueryCockRuler(log, update) }, traceQueryCreated,
		)
		resultsChan <- queryResult{index: 2, result: result}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		result := timings.ReportExecutionForResult(log.With(logging.QueryType, "CockLadder"),
			func() tgbotapi.InlineQueryResultArticle { return app.InlineQueryCockLadder(log, update) }, traceQueryCreated,
		)
		resultsChan <- queryResult{index: 3, result: result}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		result := timings.ReportExecutionForResult(log.With(logging.QueryType, "CockDynamic"),
			func() tgbotapi.InlineQueryResultArticle { return app.InlineQueryCockDynamic(log, update) }, traceQueryCreated,
		)
		resultsChan <- queryResult{index: 4, result: result}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		result := timings.ReportExecutionForResult(log.With(logging.QueryType, "CockSeason"),
			func() tgbotapi.InlineQueryResultArticle { return app.InlineQueryCockSeason(log, update) }, traceQueryCreated,
		)
		resultsChan <- queryResult{index: 5, result: result}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		result := timings.ReportExecutionForResult(log.With(logging.QueryType, "CockAchievements"),
			func() tgbotapi.InlineQueryResultArticle {
				return app.InlineQueryCockAchievements(log, update, 1)
			}, traceQueryCreated,
		)
		resultsChan <- queryResult{index: 6, result: result}
	}()

	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	parallelResults := make([]tgbotapi.InlineQueryResultArticle, parallelQueriesCount)
	for result := range resultsChan {
		parallelResults[result.index-1] = result.result
	}

	queries := make([]any, 0, parallelQueriesCount+1)
	queries = append(queries, cockSizeResult)
	for _, result := range parallelResults {
		queries = append(queries, result)
	}

	inlines := tgbotapi.InlineConfig{InlineQueryID: query.ID, IsPersonal: true, CacheTime: 60, Results: queries}

	if _, err := timings.ReportExecutionForResultError(log,
		func() (*tgbotapi.APIResponse, error) { return app.bot.Request(inlines) },
		func(l *logging.Logger) { l.I("Inline query successfully sent") },
	); err != nil {
		log.E("Failed to send inline query", logging.InnerError, err)
	}
}

func (app *Application) InlineQueryCockSize(log *logging.Logger, update *tgbotapi.Update) tgbotapi.InlineQueryResultArticle {
	query := update.InlineQuery
	if query == nil {
		return tgbotapi.InlineQueryResultArticle{}
	}
	localizer, _ := app.localization.LocalizerByUpdate(update)

	cockData, err := app.api.GenerateCockSize(app.ctx, query.From.ID, query.From.UserName)
	if err != nil {
		log.E("Failed to generate cock size via API", logging.InnerError, err)
		return tgbotapi.InlineQueryResultArticle{}
	}

	size := cockData.Size
	emoji := EmojiFromSize(size)
	text := GenerateCockSizeText(app.localization, localizer, size, emoji)
	subtext := geo.GetRegionBySize(size)
	subtext = app.localization.Localize(localizer, subtext, nil)

	text = text + "\n\n" + "_" + subtext + "_"

	return InitializeInlineQueryWithThumbAndDesc(
		app.localization.Localize(localizer, InlineTitleCockSize, nil),
		strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(text, ".", "\\."), "-", "\\-"), "!", "\\!"),
		app.localization.Localize(localizer, DescCockSize, nil),
		"https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_size.png",
	)
}

func (app *Application) InlineQueryCockLadder(log *logging.Logger, update *tgbotapi.Update) tgbotapi.InlineQueryResultArticle {
	query := update.InlineQuery
	if query == nil {
		return tgbotapi.InlineQueryResultArticle{}
	}
	localizer, _ := app.localization.LocalizerByUpdate(update)

	data, err := app.api.GetCockLadder(app.ctx, query.From.ID, query.From.UserName, 13, 1)
	if err != nil {
		log.E("Failed to get cock ladder via API", logging.InnerError, err)
		return tgbotapi.InlineQueryResultArticle{}
	}

	showDescription := app.shouldShowDescription(log, query.From.ID, query.From.UserName)
	text := app.GenerateCockLadderScoreboard(log, localizer, query.From.ID, data, showDescription)
	return InitializeInlineQueryWithThumbAndDesc(
		app.localization.Localize(localizer, InlineTitleCockLadder, nil),
		text,
		app.localization.Localize(localizer, DescCockLadder, nil),
		"https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_ladder.png",
	)
}

func (app *Application) InlineQueryCockRace(log *logging.Logger, update *tgbotapi.Update) tgbotapi.InlineQueryResultArticle {
	query := update.InlineQuery
	if query == nil {
		return tgbotapi.InlineQueryResultArticle{}
	}
	localizer, _ := app.localization.LocalizerByUpdate(update)

	data, err := app.api.GetCockRace(app.ctx, query.From.ID, query.From.UserName, 13, 1)
	if err != nil {
		log.E("Failed to get cock race via API", logging.InnerError, err)
		return tgbotapi.InlineQueryResultArticle{}
	}

	showDescription := app.shouldShowDescription(log, query.From.ID, query.From.UserName)
	text := app.GenerateCockRaceScoreboard(log, localizer, query.From.ID, data, showDescription)
	return InitializeInlineQueryWithThumbAndDesc(
		app.localization.Localize(localizer, InlineTitleCockRace, nil),
		text,
		app.localization.Localize(localizer, DescCockRace, nil),
		"https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_race.png",
	)
}

func (app *Application) InlineQueryCockDynamic(log *logging.Logger, update *tgbotapi.Update) tgbotapi.InlineQueryResultArticle {
	query := update.InlineQuery
	if query == nil {
		return tgbotapi.InlineQueryResultArticle{}
	}
	localizer, _ := app.localization.LocalizerByUpdate(update)
	userID := query.From.ID
	username := query.From.UserName

	// Три параллельных запроса: global, personal, respects
	type globalResult struct {
		data *api.CockDynamicGlobalData
		err  error
	}
	type personalResult struct {
		data *api.CockDynamicPersonalData
		err  error
	}
	type respectsResult struct {
		data *api.RespectData
		err  error
	}

	globalCh := make(chan globalResult, 1)
	personalCh := make(chan personalResult, 1)
	respectsCh := make(chan respectsResult, 1)

	go func() {
		data, err := app.api.GetCockDynamicGlobal(app.ctx)
		globalCh <- globalResult{data, err}
	}()
	go func() {
		data, err := app.api.GetCockDynamicPersonal(app.ctx, userID, username)
		personalCh <- personalResult{data, err}
	}()
	go func() {
		data, err := app.api.GetCockRespects(app.ctx, userID, username)
		respectsCh <- respectsResult{data, err}
	}()

	globalRes := <-globalCh
	personalRes := <-personalCh
	respectsRes := <-respectsCh

	if globalRes.err != nil {
		log.E("Failed to get global dynamic", logging.InnerError, globalRes.err)
		return tgbotapi.InlineQueryResultArticle{}
	}
	if personalRes.err != nil {
		log.E("Failed to get personal dynamic", logging.InnerError, personalRes.err)
		text := app.localization.Localize(localizer, MsgCockDynamicNoData, nil)
		return InitializeInlineQueryWithThumbAndDesc(
			app.localization.Localize(localizer, InlineTitleCockDynamic, nil),
			text,
			app.localization.Localize(localizer, DescCockDynamic, nil),
			"https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_dynamic.png",
		)
	}

	global := globalRes.data
	personal := personalRes.data

	// Респекты
	userSeasonWins := 0
	userCockRespect := 0
	if respectsRes.err == nil && respectsRes.data != nil {
		userCockRespect = int(respectsRes.data.TotalRespect)
	}

	// Персональный рекорд
	individualRecordTotal := personal.Record.Size
	individualRecordDate := datetime.NowTime()
	if personal.Record.RequestedAt != nil {
		if t, err := datetime.ParseUTC(*personal.Record.RequestedAt); err == nil {
			individualRecordDate = t
		}
	}

	// Общий рекорд
	overallRecordTotal := global.Record.Total
	overallRecordDate := datetime.NowTime()
	if global.Record.RequestedAt != nil {
		if t, err := datetime.ParseUTC(*global.Record.RequestedAt); err == nil {
			overallRecordDate = t
		}
	}

	// Период дёргания кока
	var userPullingPeriod string
	if personal.FirstCockDate != nil {
		if firstDate, err := datetime.ParseUTC(*personal.FirstCockDate); err == nil {
			userPullingPeriod = FormatUserPullingPeriod(app.localization, localizer, firstDate, datetime.NowTime())
		} else {
			userPullingPeriod = app.localization.Localize(localizer, MsgUserPullingRecently, nil)
		}
	} else {
		userPullingPeriod = app.localization.Localize(localizer, MsgUserPullingRecently, nil)
	}

	text := NewMsgCockDynamicsTemplate(
		app.localization,
		localizer,
		/* Общая динамика коков */
		global.TotalSize,
		global.UniqueUsers,
		int(global.Recent.Average),
		int(global.Recent.Median),

		/* Персональная динамика кока */
		personal.TotalSize,
		int(personal.RecentAverage),
		personal.Irk,
		individualRecordTotal,
		individualRecordDate,

		/* Кок-активы */
		personal.DailyDynamics.YesterdayCockChangePercent,
		personal.DailyDynamics.YesterdayCockChange,
		personal.FiveCocksDynamics.FiveCocksChangePercent,
		personal.FiveCocksDynamics.FiveCocksChange,

		/* Соотношение коков */
		global.Distribution.HugePercent,
		global.Distribution.LittlePercent,

		/* Самый большой кок */
		overallRecordDate,
		overallRecordTotal,

		/* % доминирование */
		personal.Dominance,

		/* Сезонные достижения */
		userSeasonWins,
		userCockRespect,

		/* Всего дёрнуто коков */
		global.TotalCocksCount,
		personal.CocksCount,

		/* Коэффициент везения и волатильность */
		personal.LuckCoefficient,
		personal.Volatility,

		/* Средняя скорость прироста */
		personal.GrowthSpeed,

		/* Скорость роста общей статистики */
		global.GrowthSpeed,

		/* Период дергания кока пользователем */
		userPullingPeriod,
	)

	article := tgbotapi.NewInlineQueryResultArticleMarkdown(query.ID, app.localization.Localize(localizer, InlineTitleCockDynamic, nil), text)
	article.ThumbURL = "https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_dynamic.png"
	article.Description = app.localization.Localize(localizer, DescCockDynamic, nil)
	return article
}

func (app *Application) InlineQueryCockSeason(log *logging.Logger, update *tgbotapi.Update) tgbotapi.InlineQueryResultArticle {
	query := update.InlineQuery
	if query == nil {
		return tgbotapi.InlineQueryResultArticle{}
	}
	localizer, _ := app.localization.LocalizerByUpdate(update)

	seasonsData, err := app.api.GetCockSeasons(app.ctx, query.From.ID, query.From.UserName, 15, 1)
	if err != nil {
		log.E("Failed to get cock seasons via API", logging.InnerError, err)
		text := NewMsgCockSeasonNoSeasonsTemplate(app.localization, localizer)
		return InitializeInlineQueryWithThumbAndDesc(
			app.localization.Localize(localizer, InlineTitleCockSeason, nil),
			text,
			app.localization.Localize(localizer, DescCockSeason, nil),
			"https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_seasons.png",
		)
	}

	if len(seasonsData.Seasons) == 0 {
		text := NewMsgCockSeasonNoSeasonsTemplate(app.localization, localizer)
		return InitializeInlineQueryWithThumbAndDesc(
			app.localization.Localize(localizer, InlineTitleCockSeason, nil),
			text,
			app.localization.Localize(localizer, DescCockSeason, nil),
			"https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_seasons.png",
		)
	}

	// Первый элемент = самый новый (текущий) сезон
	currentSeason := seasonsData.Seasons[0]
	showDescription := app.shouldShowDescription(log, query.From.ID, query.From.UserName)
	text := generateSeasonPageText(app.localization, localizer, currentSeason, showDescription)

	// Кнопки навигации
	var buttons []tgbotapi.InlineKeyboardButton

	// Кнопка "предыдущий сезон" (более старый)
	if len(seasonsData.Seasons) > 1 {
		prevSeason := seasonsData.Seasons[1]
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
			app.localization.Localize(localizer, MsgSeasonButton, map[string]any{
				"Arrow":     "◀️",
				"SeasonNum": prevSeason.SeasonNum,
			}),
			fmt.Sprintf("season_page:%d", prevSeason.SeasonNum),
		))
	}

	article := tgbotapi.NewInlineQueryResultArticleMarkdownV2(
		fmt.Sprintf("season_%d", currentSeason.SeasonNum),
		app.localization.Localize(localizer, InlineTitleCockSeason, nil),
		text,
	)

	if len(buttons) > 0 {
		kb := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(buttons...),
		)
		article.ReplyMarkup = &kb
	}
	article.ThumbURL = "https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_seasons.png"
	article.Description = app.localization.Localize(localizer, DescCockSeason, nil)

	return article
}

func (app *Application) InlineQueryCockRuler(log *logging.Logger, update *tgbotapi.Update) tgbotapi.InlineQueryResultArticle {
	query := update.InlineQuery
	if query == nil {
		return tgbotapi.InlineQueryResultArticle{}
	}
	localizer, _ := app.localization.LocalizerByUpdate(update)

	data, err := app.api.GetCockRuler(app.ctx, query.From.ID, query.From.UserName, 13, 1)
	if err != nil {
		log.E("Failed to get cock ruler via API", logging.InnerError, err)
		return tgbotapi.InlineQueryResultArticle{}
	}

	showDescription := app.shouldShowDescription(log, query.From.ID, query.From.UserName)
	text := app.GenerateCockRulerText(log, localizer, query.From.ID, data, showDescription)
	return InitializeInlineQueryWithThumbAndDesc(
		app.localization.Localize(localizer, InlineTitleCockRuler, nil),
		text,
		app.localization.Localize(localizer, DescCockRuler, nil),
		"https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_ruler.png",
	)
}

func (app *Application) InlineQueryCockAchievements(log *logging.Logger, update *tgbotapi.Update, page int) tgbotapi.InlineQueryResultArticle {
	query := update.InlineQuery
	if query == nil {
		return tgbotapi.InlineQueryResultArticle{}
	}
	localizer, _ := app.localization.LocalizerByUpdate(update)
	userID := query.From.ID

	// Параллельно получаем ачивки и респекты
	type achResult struct {
		data *api.CockAchievementsData
		err  error
	}
	type respResult struct {
		data *api.RespectData
		err  error
	}

	achCh := make(chan achResult, 1)
	respCh := make(chan respResult, 1)

	go func() {
		data, err := app.api.GetCockAchievements(app.ctx, userID, query.From.UserName)
		achCh <- achResult{data, err}
	}()
	go func() {
		data, err := app.api.GetCockRespects(app.ctx, userID, query.From.UserName)
		respCh <- respResult{data, err}
	}()

	achRes := <-achCh
	respRes := <-respCh

	if achRes.err != nil {
		log.E("Failed to get achievements via API", logging.InnerError, achRes.err)
		return tgbotapi.InlineQueryResultArticle{}
	}
	achData := achRes.data

	achievementRespects := 0
	if respRes.err == nil && respRes.data != nil {
		achievementRespects = int(respRes.data.AchievementRespect)
	}

	achievementsList := GenerateAchievementsText(
		app.localization,
		localizer,
		AllAchievements,
		achData.Achievements,
		page,
		10,
	)

	totalAchievements := achData.AchievementsTotal
	totalPages := (totalAchievements + 9) / 10

	var templateID string
	if page == 1 {
		templateID = MsgCockAchievementsTemplate
	} else {
		templateID = MsgCockAchievementsTemplateOtherPages
	}

	text := app.localization.Localize(localizer, templateID, map[string]any{
		"Completed":    achData.AchievementsDone,
		"Total":        totalAchievements,
		"Percent":      int(achData.AchievementsDonePercent),
		"Respects":     achievementRespects,
		"Achievements": achievementsList,
	})

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

	article := tgbotapi.NewInlineQueryResultArticleMarkdownV2(
		fmt.Sprintf("ach_%d_%d", userID, page),
		app.localization.Localize(localizer, InlineTitleCockAchievements, nil),
		text,
	)
	article.ReplyMarkup = &kb
	article.ThumbURL = "https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_achievements.png"
	article.Description = app.localization.Localize(localizer, DescCockAchievements, nil)

	return article
}

func InitializeInlineQueryWithThumbAndDesc(title, message, description, thumbURL string) tgbotapi.InlineQueryResultArticle {
	article := tgbotapi.NewInlineQueryResultArticleMarkdownV2(fmt.Sprintf("q_%d", time.Now().UnixNano()), title, message)
	article.ThumbURL = thumbURL
	article.Description = description
	return article
}

func (app *Application) HandleCallbackQuery(log *logging.Logger, update *tgbotapi.Update) {
	callback := update.CallbackQuery
	if callback == nil {
		return
	}
	localizer, _ := app.localization.LocalizerByUpdate(update)
	data := callback.Data

	if strings.HasPrefix(data, hideCallbackPrefix) {
		app.handleHideCallback(log, localizer, callback)
		return
	}

	if strings.HasPrefix(data, "season_page:") {
		app.handleSeasonPageCallback(log, localizer, callback)
		return
	}

	if strings.HasPrefix(data, "ach_page:") {
		app.handleAchPageCallback(log, localizer, callback)
		return
	}

	if data == "ach_noop" {
		callbackConfig := tgbotapi.NewCallback(callback.ID, "")
		if _, err := app.bot.Request(callbackConfig); err != nil {
			log.E("Failed to answer callback query", logging.InnerError, err)
		}
	}
}

func (app *Application) handleHideCallback(log *logging.Logger, localizer *i18n.Localizer, callback *tgbotapi.CallbackQuery) {
	data := callback.Data
	parts := strings.Split(strings.TrimPrefix(data, hideCallbackPrefix), ":")
	if len(parts) != 2 {
		log.E("Invalid hide callback data format", "data", data)
		callbackConfig := tgbotapi.NewCallback(callback.ID, app.localization.Localize(localizer, MsgCallbackInvalidFormat, nil))
		if _, err := app.bot.Request(callbackConfig); err != nil {
			log.E("Failed to answer callback query", logging.InnerError, err)
		}
		return
	}

	targetUserID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		log.E("Failed to parse userID from hide callback", logging.InnerError, err)
		callbackConfig := tgbotapi.NewCallback(callback.ID, app.localization.Localize(localizer, MsgCallbackParseError, nil))
		if _, err := app.bot.Request(callbackConfig); err != nil {
			log.E("Failed to answer callback query", logging.InnerError, err)
		}
		return
	}

	action := parts[1]
	hide := false
	switch action {
	case hideActionHide:
		hide = true
	case hideActionShow:
		hide = false
	default:
		log.E("Invalid hide callback action", "action", action)
		callbackConfig := tgbotapi.NewCallback(callback.ID, app.localization.Localize(localizer, MsgCallbackInvalidFormat, nil))
		if _, err := app.bot.Request(callbackConfig); err != nil {
			log.E("Failed to answer callback query", logging.InnerError, err)
		}
		return
	}

	if callback.From == nil || callback.From.ID != targetUserID {
		callbackConfig := tgbotapi.NewCallback(callback.ID, app.localization.Localize(localizer, MsgCallbackNotForYou, nil))
		callbackConfig.ShowAlert = true
		if _, err := app.bot.Request(callbackConfig); err != nil {
			log.E("Failed to answer callback query", logging.InnerError, err)
		}
		return
	}

	anonName, realName := app.setUserHiddenStatus(log, localizer, callback.From, hide)
	text, keyboard := app.buildHideMessage(localizer, hide, anonName, realName, targetUserID)

	_, _ = app.bot.Request(tgbotapi.NewCallback(callback.ID, ""))

	app.editCallbackMessage(log, callback, text, keyboard)
}

func (app *Application) handleSeasonPageCallback(log *logging.Logger, localizer *i18n.Localizer, callback *tgbotapi.CallbackQuery) {
	seasonNumStr := strings.TrimPrefix(callback.Data, "season_page:")
	seasonNum := 1
	if parsedSeasonNum, err := strconv.Atoi(seasonNumStr); err != nil {
		log.E("Failed to parse season number", logging.InnerError, err)
	} else {
		seasonNum = parsedSeasonNum
	}

	userID := int64(0)
	username := ""
	if callback.From != nil {
		userID = callback.From.ID
		username = callback.From.UserName
	}

	// Запрашиваем сезоны с бэкэнда
	seasonsData, err := app.api.GetCockSeasons(app.ctx, userID, username, 15, 1)
	if err != nil {
		log.E("Failed to get seasons via API", logging.InnerError, err)
		callbackConfig := tgbotapi.NewCallback(callback.ID, app.localization.Localize(localizer, MsgSeasonNotFound, nil))
		if _, err := app.bot.Request(callbackConfig); err != nil {
			log.E("Failed to answer callback query", logging.InnerError, err)
		}
		return
	}

	// Ищем нужный сезон в полученных данных
	var targetSeason *api.SeasonWithWinners
	var targetIdx int
	for idx, s := range seasonsData.Seasons {
		if s.SeasonNum == seasonNum {
			targetSeason = &seasonsData.Seasons[idx]
			targetIdx = idx
			break
		}
	}

	// Если не нашли на первой странице, пробуем вычислить нужную API-страницу
	if targetSeason == nil && seasonsData.Page.TotalPages > 1 {
		// Сезоны от новых к старым: сезон с наибольшим номером на первой странице
		// Вычисляем apiPage для нужного сезона
		for apiPage := 2; apiPage <= seasonsData.Page.TotalPages; apiPage++ {
			pageData, err := app.api.GetCockSeasons(app.ctx, userID, username, 15, apiPage)
			if err != nil {
				break
			}
			for idx, s := range pageData.Seasons {
				if s.SeasonNum == seasonNum {
					targetSeason = &pageData.Seasons[idx]
					targetIdx = idx
					seasonsData = pageData
					break
				}
			}
			if targetSeason != nil {
				break
			}
		}
	}

	if targetSeason == nil {
		log.E("Season not found", "season_num", seasonNum)
		callbackConfig := tgbotapi.NewCallback(callback.ID, app.localization.Localize(localizer, MsgSeasonNotFound, nil))
		if _, err := app.bot.Request(callbackConfig); err != nil {
			log.E("Failed to answer callback query", logging.InnerError, err)
		}
		return
	}

	showDescription := app.shouldShowDescription(log, userID, username)
	text := generateSeasonPageText(app.localization, localizer, *targetSeason, showDescription)

	// Кнопки навигации
	var buttons []tgbotapi.InlineKeyboardButton

	// Кнопка "предыдущий сезон" (более старый)
	if targetIdx < len(seasonsData.Seasons)-1 {
		prevSeason := seasonsData.Seasons[targetIdx+1]
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
			app.localization.Localize(localizer, MsgSeasonButton, map[string]any{
				"Arrow":     "◀️",
				"SeasonNum": prevSeason.SeasonNum,
			}),
			fmt.Sprintf("season_page:%d", prevSeason.SeasonNum),
		))
	} else if targetSeason.SeasonNum > 1 {
		// Есть более старые сезоны, но на другой странице API
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
			app.localization.Localize(localizer, MsgSeasonButton, map[string]any{
				"Arrow":     "◀️",
				"SeasonNum": targetSeason.SeasonNum - 1,
			}),
			fmt.Sprintf("season_page:%d", targetSeason.SeasonNum-1),
		))
	}

	// Кнопка "следующий сезон" (более новый)
	if targetIdx > 0 {
		nextSeason := seasonsData.Seasons[targetIdx-1]
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
			app.localization.Localize(localizer, MsgSeasonButton, map[string]any{
				"Arrow":     "▶️",
				"SeasonNum": nextSeason.SeasonNum,
			}),
			fmt.Sprintf("season_page:%d", nextSeason.SeasonNum),
		))
	}

	_, _ = app.bot.Request(tgbotapi.NewCallback(callback.ID, ""))

	var keyboard *tgbotapi.InlineKeyboardMarkup
	if len(buttons) > 0 {
		kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
		keyboard = &kb
	}
	app.editCallbackMessage(log, callback, text, keyboard)
}

func (app *Application) handleAchPageCallback(log *logging.Logger, localizer *i18n.Localizer, callback *tgbotapi.CallbackQuery) {
	parts := strings.Split(strings.TrimPrefix(callback.Data, "ach_page:"), ":")
	if len(parts) != 2 {
		log.E("Invalid ach_page callback data format", "data", callback.Data)
		callbackConfig := tgbotapi.NewCallback(callback.ID, app.localization.Localize(localizer, MsgCallbackInvalidFormat, nil))
		if _, err := app.bot.Request(callbackConfig); err != nil {
			log.E("Failed to answer callback query", logging.InnerError, err)
		}
		return
	}

	userID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		log.E("Failed to parse userID from callback", logging.InnerError, err)
		callbackConfig := tgbotapi.NewCallback(callback.ID, app.localization.Localize(localizer, MsgCallbackParseError, nil))
		if _, err := app.bot.Request(callbackConfig); err != nil {
			log.E("Failed to answer callback query", logging.InnerError, err)
		}
		return
	}

	page := 1
	if parsedPage, err := strconv.Atoi(parts[1]); err != nil {
		log.E("Failed to parse page number", logging.InnerError, err)
	} else {
		page = parsedPage
	}

	username := ""
	if callback.From != nil {
		username = callback.From.UserName
	}

	// Параллельно получаем ачивки и респекты
	type achResult struct {
		data *api.CockAchievementsData
		err  error
	}
	type respResult struct {
		data *api.RespectData
		err  error
	}

	achCh := make(chan achResult, 1)
	respCh := make(chan respResult, 1)

	go func() {
		data, err := app.api.GetCockAchievements(app.ctx, userID, username)
		achCh <- achResult{data, err}
	}()
	go func() {
		data, err := app.api.GetCockRespects(app.ctx, userID, username)
		respCh <- respResult{data, err}
	}()

	achRes := <-achCh
	respRes := <-respCh

	if achRes.err != nil {
		log.E("Failed to get achievements via API", logging.InnerError, achRes.err)
		_, _ = app.bot.Request(tgbotapi.NewCallback(callback.ID, ""))
		return
	}
	achData := achRes.data

	achievementRespects := 0
	if respRes.err == nil && respRes.data != nil {
		achievementRespects = int(respRes.data.AchievementRespect)
	}

	achievementsList := GenerateAchievementsText(
		app.localization,
		localizer,
		AllAchievements,
		achData.Achievements,
		page,
		10,
	)

	totalAchievements := achData.AchievementsTotal
	totalPages := (totalAchievements + 9) / 10

	var templateID string
	if page == 1 {
		templateID = MsgCockAchievementsTemplate
	} else {
		templateID = MsgCockAchievementsTemplateOtherPages
	}

	text := app.localization.Localize(localizer, templateID, map[string]any{
		"Completed":    achData.AchievementsDone,
		"Total":        totalAchievements,
		"Percent":      int(achData.AchievementsDonePercent),
		"Respects":     achievementRespects,
		"Achievements": achievementsList,
	})

	var buttons []tgbotapi.InlineKeyboardButton

	if page > 1 {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("◀️", fmt.Sprintf("ach_page:%d:%d", userID, page-1)))
	}

	buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d/%d", page, totalPages), "ach_noop"))

	if page < totalPages {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("▶️", fmt.Sprintf("ach_page:%d:%d", userID, page+1)))
	}

	_, _ = app.bot.Request(tgbotapi.NewCallback(callback.ID, ""))

	kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
	app.editCallbackMessage(log, callback, text, &kb)
}

// editCallbackMessage — вспомогательная функция для редактирования сообщения из callback
func (app *Application) editCallbackMessage(log *logging.Logger, callback *tgbotapi.CallbackQuery, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	if callback.InlineMessageID != "" {
		edit := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				InlineMessageID: callback.InlineMessageID,
			},
			Text:      text,
			ParseMode: "MarkdownV2",
		}
		if keyboard != nil {
			edit.ReplyMarkup = keyboard
		}
		if _, err := app.bot.Request(edit); err != nil {
			log.E("Failed to edit inline message", logging.InnerError, err)
		}
	} else if callback.Message != nil {
		edit := tgbotapi.NewEditMessageText(
			callback.Message.Chat.ID,
			callback.Message.MessageID,
			text,
		)
		edit.ParseMode = "MarkdownV2"
		if keyboard != nil {
			edit.ReplyMarkup = keyboard
		}
		if _, err := app.bot.Request(edit); err != nil {
			log.E("Failed to edit chat message", logging.InnerError, err)
		}
	} else {
		log.E("CallbackQuery has neither Message nor InlineMessageID")
	}
}

// generateSeasonPageText генерирует текст для одной страницы сезона из API-данных
func generateSeasonPageText(locMgr *localization.LocalizationManager, localizer *i18n.Localizer, season api.SeasonWithWinners, showDescription bool) string {
	startDate := EscapeMarkdownV2(datetime.FormatDateMSK(season.StartDate))
	endDate := EscapeMarkdownV2(datetime.FormatDateMSK(season.EndDate))

	var winnerLines []string
	for _, winner := range season.Winners {
		medal := GetMedalByPosition(winner.Place - 1)
		line := NewMsgCockSeasonWinnerTemplate(
			locMgr,
			localizer,
			medal,
			winner.Nickname,
			FormatDickSize(winner.TotalSize),
		)
		winnerLines = append(winnerLines, line)
	}

	winnersText := strings.Join(winnerLines, "\n")

	if season.IsActive {
		seasonBlock := NewMsgCockSeasonTemplate(locMgr, localizer, winnersText, startDate, endDate, season.SeasonNum)
		if showDescription {
			footer := NewMsgCockSeasonTemplateFooter(locMgr, localizer)
			return seasonBlock + "\n\n" + footer
		}
		return seasonBlock
	}

	return NewMsgCockSeasonWithWinnersTemplate(locMgr, localizer, winnersText, startDate, endDate, season.SeasonNum)
}
