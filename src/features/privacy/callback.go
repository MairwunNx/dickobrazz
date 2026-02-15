package privacy

import (
	"context"
	"dickobrazz/src/shared/formatting"
	"dickobrazz/src/shared/localization"
	"dickobrazz/src/shared/logging"
	"dickobrazz/src/shared/telegram"
	"fmt"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func (h *Handler) HandleCallback(ctx context.Context, log *logging.Logger, localizer *i18n.Localizer, callback *tgbotapi.CallbackQuery) {
	data := callback.Data
	parts := strings.Split(strings.TrimPrefix(data, HideCallbackPrefix), ":")
	if len(parts) != 2 {
		log.E("Invalid hide callback data format", "data", data)
		callbackConfig := tgbotapi.NewCallback(callback.ID, h.loc.Localize(localizer, localization.MsgCallbackInvalidFormat, nil))
		if _, err := h.bot.Request(callbackConfig); err != nil {
			log.E("Failed to answer callback query", logging.InnerError, err)
		}
		return
	}

	targetUserID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		log.E("Failed to parse userID from hide callback", logging.InnerError, err)
		callbackConfig := tgbotapi.NewCallback(callback.ID, h.loc.Localize(localizer, localization.MsgCallbackParseError, nil))
		if _, err := h.bot.Request(callbackConfig); err != nil {
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
		callbackConfig := tgbotapi.NewCallback(callback.ID, h.loc.Localize(localizer, localization.MsgCallbackInvalidFormat, nil))
		if _, err := h.bot.Request(callbackConfig); err != nil {
			log.E("Failed to answer callback query", logging.InnerError, err)
		}
		return
	}

	if callback.From == nil || callback.From.ID != targetUserID {
		callbackConfig := tgbotapi.NewCallback(callback.ID, h.loc.Localize(localizer, "MsgCallbackNotForYou", nil))
		callbackConfig.ShowAlert = true
		if _, err := h.bot.Request(callbackConfig); err != nil {
			log.E("Failed to answer callback query", logging.InnerError, err)
		}
		return
	}

	toggleAction := NewToggleAction(h.api, h.loc)
	anonName, realName := toggleAction.Execute(ctx, log, localizer, callback.From, hide)

	msgKey := localization.MsgHidePrompt
	buttonKey := localization.MsgHideButtonHide
	btnAction := hideActionHide

	if hide {
		msgKey = localization.MsgHideStatusHidden
		buttonKey = localization.MsgHideButtonShow
		btnAction = hideActionShow
	}

	text := h.loc.Localize(localizer, msgKey, map[string]any{
		"Anon": formatting.EscapeMarkdownV2(anonName),
		"Real": formatting.EscapeMarkdownV2(realName),
	})
	button := tgbotapi.NewInlineKeyboardButtonData(
		h.loc.Localize(localizer, buttonKey, nil),
		fmt.Sprintf("%s%d:%s", HideCallbackPrefix, targetUserID, btnAction),
	)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))

	_, _ = h.bot.Request(tgbotapi.NewCallback(callback.ID, ""))

	telegram.EditCallbackMessage(log, h.bot, callback, text, &keyboard)
}
