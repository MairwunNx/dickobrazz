package app

import (
	"context"
	"dickobrazz/src/features/achievements"
	"dickobrazz/src/features/dynamics"
	"dickobrazz/src/features/help"
	"dickobrazz/src/features/ladder"
	"dickobrazz/src/features/privacy"
	"dickobrazz/src/features/race"
	"dickobrazz/src/features/ruler"
	"dickobrazz/src/features/seasons"
	"dickobrazz/src/features/size"
	"dickobrazz/src/shared/api"
	"dickobrazz/src/shared/localization"
	"dickobrazz/src/shared/logging"
	"dickobrazz/src/shared/metrics"
	"dickobrazz/src/shared/timings"
	"strings"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	ctx context.Context
	bot *tgbotapi.BotAPI
	log *logging.Logger
	loc *localization.LocalizationManager
	api *api.APIClient

	sizeHandler          *size.Handler
	rulerHandler         *ruler.Handler
	ladderHandler        *ladder.Handler
	raceHandler          *race.Handler
	dynamicsHandler      *dynamics.Handler
	seasonsHandler       *seasons.Handler
	seasonsCallback      *seasons.CallbackHandler
	achievementsHandler  *achievements.Handler
	achievementsCallback *achievements.CallbackHandler
	helpHandler          *help.Handler
	privacyHandler       *privacy.Handler
}

type RouterParams struct {
	Bot *tgbotapi.BotAPI
	Log *logging.Logger
	Loc *localization.LocalizationManager
	API *api.APIClient

	SizeHandler          *size.Handler
	RulerHandler         *ruler.Handler
	LadderHandler        *ladder.Handler
	RaceHandler          *race.Handler
	DynamicsHandler      *dynamics.Handler
	SeasonsHandler       *seasons.Handler
	SeasonsCallback      *seasons.CallbackHandler
	AchievementsHandler  *achievements.Handler
	AchievementsCallback *achievements.CallbackHandler
	HelpHandler          *help.Handler
	PrivacyHandler       *privacy.Handler
}

func NewRouter(p RouterParams) *Router {
	return &Router{
		bot:                  p.Bot,
		log:                  p.Log,
		loc:                  p.Loc,
		api:                  p.API,
		sizeHandler:          p.SizeHandler,
		rulerHandler:         p.RulerHandler,
		ladderHandler:        p.LadderHandler,
		raceHandler:          p.RaceHandler,
		dynamicsHandler:      p.DynamicsHandler,
		seasonsHandler:       p.SeasonsHandler,
		seasonsCallback:      p.SeasonsCallback,
		achievementsHandler:  p.AchievementsHandler,
		achievementsCallback: p.AchievementsCallback,
		helpHandler:          p.HelpHandler,
		privacyHandler:       p.PrivacyHandler,
	}
}

func (r *Router) SetContext(ctx context.Context) {
	r.ctx = ctx
}

func (r *Router) shouldShowDescription(log *logging.Logger, userID int64, username string) bool {
	profile, err := r.api.GetProfile(r.ctx, userID, username)
	if err != nil {
		log.E("Failed to get profile for shouldShowDescription", logging.InnerError, err)
		return true
	}

	if profile.CreatedAt == nil {
		return true
	}

	if time.Since(profile.CreatedAt.Time).Hours()/24 > 32 {
		return false
	}

	return true
}

func (r *Router) HandleUpdate(update tgbotapi.Update) {
	processingStarted := time.Now()
	handledKinds := 0
	updateKind := "ignored"
	if _, detectedLang := r.loc.LocalizerByUpdate(&update); detectedLang != "" {
		metrics.IncDetectedLanguage(detectedLang)
	}

	if msg := update.Message; msg != nil {
		user := update.SentFrom()
		log := r.log.With(
			logging.UserId, user.ID,
			logging.UserName, user.UserName,
			logging.ChatType, msg.Chat.Type,
			logging.ChatId, msg.Chat.ID,
		)
		log.I("Received message")
		metrics.IncMessagesHandled("message")
		handledKinds++
		updateKind = "message"

		if msg.IsCommand() {
			r.handleCommand(log, &update)
		}
	}

	if query := update.InlineQuery; query != nil {
		user := update.SentFrom()
		log := r.log.With(
			logging.UserId, user.ID,
			logging.UserName, user.UserName,
			logging.QueryId, query.ID,
			logging.ChatType, query.ChatType,
		)

		r.handleInlineQuery(log, &update)
		metrics.IncMessagesHandled("inline_query")
		handledKinds++
		updateKind = "inline_query"
	}

	if callback := update.CallbackQuery; callback != nil {
		user := update.SentFrom()
		log := r.log.With(
			logging.UserId, user.ID,
			logging.UserName, user.UserName,
			"callback_id", callback.ID,
			"callback_data", callback.Data,
		)

		r.handleCallbackQuery(log, &update)
		metrics.IncMessagesHandled("callback_query")
		handledKinds++
		updateKind = "callback_query"
	}

	if handledKinds == 0 {
		metrics.IncMessagesIgnored("unsupported_update")
		updateKind = "ignored"
	} else if handledKinds > 1 {
		updateKind = "multiple"
	}
	metrics.ObserveUpdateDuration(updateKind, time.Since(processingStarted))
}

func (r *Router) handleCommand(log *logging.Logger, update *tgbotapi.Update) {
	msg := update.Message
	if msg == nil || !msg.IsCommand() {
		return
	}
	localizer, _ := r.loc.LocalizerByUpdate(update)

	switch msg.Command() {
	case "help":
		r.helpHandler.HandleCommand(log, localizer, msg)
	case "hide":
		r.privacyHandler.HandleCommand(r.ctx, log, localizer, msg)
	}
}

func (r *Router) handleInlineQuery(log *logging.Logger, update *tgbotapi.Update) {
	query := update.InlineQuery
	if query == nil {
		return
	}
	var traceQueryCreated = func(l *logging.Logger) { l.I("Inline query successfully created") }

	cockSizeResult := timings.ReportExecutionForResult(log.With(logging.QueryType, "CockSize"),
		func() tgbotapi.InlineQueryResultArticle {
			return r.sizeHandler.HandleInlineQuery(r.ctx, log, update)
		}, traceQueryCreated,
	)

	type queryResult struct {
		index  int
		result tgbotapi.InlineQueryResultArticle
	}

	showDescription := r.shouldShowDescription(log, query.From.ID, query.From.UserName)

	parallelQueriesCount := 6
	resultsChan := make(chan queryResult, parallelQueriesCount)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		result := timings.ReportExecutionForResult(log.With(logging.QueryType, "CockRace"),
			func() tgbotapi.InlineQueryResultArticle {
				return r.raceHandler.HandleInlineQuery(r.ctx, log, update, showDescription)
			}, traceQueryCreated,
		)
		resultsChan <- queryResult{index: 1, result: result}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		result := timings.ReportExecutionForResult(log.With(logging.QueryType, "CockRuler"),
			func() tgbotapi.InlineQueryResultArticle {
				return r.rulerHandler.HandleInlineQuery(r.ctx, log, update, showDescription)
			}, traceQueryCreated,
		)
		resultsChan <- queryResult{index: 2, result: result}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		result := timings.ReportExecutionForResult(log.With(logging.QueryType, "CockLadder"),
			func() tgbotapi.InlineQueryResultArticle {
				return r.ladderHandler.HandleInlineQuery(r.ctx, log, update, showDescription)
			}, traceQueryCreated,
		)
		resultsChan <- queryResult{index: 3, result: result}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		result := timings.ReportExecutionForResult(log.With(logging.QueryType, "CockDynamic"),
			func() tgbotapi.InlineQueryResultArticle {
				return r.dynamicsHandler.HandleInlineQuery(r.ctx, log, update)
			}, traceQueryCreated,
		)
		resultsChan <- queryResult{index: 4, result: result}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		result := timings.ReportExecutionForResult(log.With(logging.QueryType, "CockSeason"),
			func() tgbotapi.InlineQueryResultArticle {
				return r.seasonsHandler.HandleInlineQuery(r.ctx, log, update, showDescription)
			}, traceQueryCreated,
		)
		resultsChan <- queryResult{index: 5, result: result}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		result := timings.ReportExecutionForResult(log.With(logging.QueryType, "CockAchievements"),
			func() tgbotapi.InlineQueryResultArticle {
				return r.achievementsHandler.HandleInlineQuery(r.ctx, log, update, 1)
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
		func() (*tgbotapi.APIResponse, error) { return r.bot.Request(inlines) },
		func(l *logging.Logger) { l.I("Inline query successfully sent") },
	); err != nil {
		log.E("Failed to send inline query", logging.InnerError, err)
	}
}

func (r *Router) handleCallbackQuery(log *logging.Logger, update *tgbotapi.Update) {
	callback := update.CallbackQuery
	if callback == nil {
		return
	}
	localizer, _ := r.loc.LocalizerByUpdate(update)
	data := callback.Data

	if strings.HasPrefix(data, privacy.HideCallbackPrefix) {
		r.privacyHandler.HandleCallback(r.ctx, log, localizer, callback)
		return
	}

	if strings.HasPrefix(data, "season_page:") {
		showDescription := r.shouldShowDescription(log, callback.From.ID, callback.From.UserName)
		r.seasonsCallback.HandleCallback(r.ctx, log, localizer, callback, showDescription)
		return
	}

	if strings.HasPrefix(data, "ach_page:") {
		r.achievementsCallback.HandleCallback(r.ctx, log, localizer, callback)
		return
	}

	if data == "ach_noop" {
		achievements.HandleAchNoop(r.bot, log, callback)
	}
}
