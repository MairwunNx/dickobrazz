package achievements

import (
	"context"
	"dickobrazz/src/shared/localization"
	"dickobrazz/src/shared/logging"
	"dickobrazz/src/shared/telegram"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type CallbackHandler struct {
	action *GetAction
	loc    *localization.LocalizationManager
	bot    *tgbotapi.BotAPI
}

func NewCallbackHandler(action *GetAction, loc *localization.LocalizationManager, bot *tgbotapi.BotAPI) *CallbackHandler {
	return &CallbackHandler{action: action, loc: loc, bot: bot}
}

func (h *CallbackHandler) HandleCallback(ctx context.Context, log *logging.Logger, localizer *i18n.Localizer, callback *tgbotapi.CallbackQuery) {
	parts := strings.Split(strings.TrimPrefix(callback.Data, "ach_page:"), ":")
	if len(parts) != 2 {
		log.E("Invalid ach_page callback data format", "data", callback.Data)
		callbackConfig := tgbotapi.NewCallback(callback.ID, h.loc.Localize(localizer, localization.MsgCallbackInvalidFormat, nil))
		if _, err := h.bot.Request(callbackConfig); err != nil {
			log.E("Failed to answer callback query", logging.InnerError, err)
		}
		return
	}

	userID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		log.E("Failed to parse userID from callback", logging.InnerError, err)
		callbackConfig := tgbotapi.NewCallback(callback.ID, h.loc.Localize(localizer, localization.MsgCallbackParseError, nil))
		if _, err := h.bot.Request(callbackConfig); err != nil {
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

	result, err := h.action.Execute(ctx, log, localizer, userID, username, page)
	if err != nil {
		_, _ = h.bot.Request(tgbotapi.NewCallback(callback.ID, ""))
		return
	}

	_, _ = h.bot.Request(tgbotapi.NewCallback(callback.ID, ""))

	kb := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(result.Buttons...))
	telegram.EditCallbackMessage(log, h.bot, callback, result.Text, &kb)
}

func HandleAchNoop(bot *tgbotapi.BotAPI, log *logging.Logger, callback *tgbotapi.CallbackQuery) {
	callbackConfig := tgbotapi.NewCallback(callback.ID, "")
	if _, err := bot.Request(callbackConfig); err != nil {
		log.E("Failed to answer callback query", logging.InnerError, err)
	}
}
