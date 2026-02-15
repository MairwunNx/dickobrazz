package seasons

import (
	"context"
	"dickobrazz/src/shared/api"
	"dickobrazz/src/shared/localization"
	"dickobrazz/src/shared/logging"
	"dickobrazz/src/shared/telegram"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type CallbackHandler struct {
	api *api.APIClient
	loc *localization.LocalizationManager
	bot *tgbotapi.BotAPI
}

func NewCallbackHandler(apiClient *api.APIClient, loc *localization.LocalizationManager, bot *tgbotapi.BotAPI) *CallbackHandler {
	return &CallbackHandler{api: apiClient, loc: loc, bot: bot}
}

func (h *CallbackHandler) HandleCallback(ctx context.Context, log *logging.Logger, localizer *i18n.Localizer, callback *tgbotapi.CallbackQuery, showDescription bool) {
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

	seasonsData, err := h.api.GetCockSeasons(ctx, userID, username, 15, 1)
	if err != nil {
		log.E("Failed to get seasons via API", logging.InnerError, err)
		callbackConfig := tgbotapi.NewCallback(callback.ID, h.loc.Localize(localizer, localization.MsgSeasonNotFound, nil))
		if _, err := h.bot.Request(callbackConfig); err != nil {
			log.E("Failed to answer callback query", logging.InnerError, err)
		}
		return
	}

	var targetSeason *api.SeasonWithWinners
	var targetIdx int
	for idx, s := range seasonsData.Seasons {
		if s.SeasonNum == seasonNum {
			targetSeason = &seasonsData.Seasons[idx]
			targetIdx = idx
			break
		}
	}

	if targetSeason == nil && seasonsData.Page.TotalPages > 1 {
		for apiPage := 2; apiPage <= seasonsData.Page.TotalPages; apiPage++ {
			pageData, err := h.api.GetCockSeasons(ctx, userID, username, 15, apiPage)
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
		callbackConfig := tgbotapi.NewCallback(callback.ID, h.loc.Localize(localizer, localization.MsgSeasonNotFound, nil))
		if _, err := h.bot.Request(callbackConfig); err != nil {
			log.E("Failed to answer callback query", logging.InnerError, err)
		}
		return
	}

	text := generateSeasonPageText(h.loc, localizer, *targetSeason, showDescription)

	var buttons []tgbotapi.InlineKeyboardButton

	if targetIdx < len(seasonsData.Seasons)-1 {
		prevSeason := seasonsData.Seasons[targetIdx+1]
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
			h.loc.Localize(localizer, localization.MsgSeasonButton, map[string]any{
				"Arrow":     "◀️",
				"SeasonNum": prevSeason.SeasonNum,
			}),
			fmt.Sprintf("season_page:%d", prevSeason.SeasonNum),
		))
	} else if targetSeason.SeasonNum > 1 {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
			h.loc.Localize(localizer, localization.MsgSeasonButton, map[string]any{
				"Arrow":     "◀️",
				"SeasonNum": targetSeason.SeasonNum - 1,
			}),
			fmt.Sprintf("season_page:%d", targetSeason.SeasonNum-1),
		))
	}

	if targetIdx > 0 {
		nextSeason := seasonsData.Seasons[targetIdx-1]
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
			h.loc.Localize(localizer, localization.MsgSeasonButton, map[string]any{
				"Arrow":     "▶️",
				"SeasonNum": nextSeason.SeasonNum,
			}),
			fmt.Sprintf("season_page:%d", nextSeason.SeasonNum),
		))
	}

	_, _ = h.bot.Request(tgbotapi.NewCallback(callback.ID, ""))

	var keyboard *tgbotapi.InlineKeyboardMarkup
	if len(buttons) > 0 {
		kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
		keyboard = &kb
	}
	telegram.EditCallbackMessage(log, h.bot, callback, text, keyboard)
}
