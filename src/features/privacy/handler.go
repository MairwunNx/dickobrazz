package privacy

import (
	"context"
	"dickobrazz/src/shared/api"
	"dickobrazz/src/shared/formatting"
	"dickobrazz/src/shared/localization"
	"dickobrazz/src/shared/logging"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	HideCallbackPrefix = "hide_toggle:"
	hideActionHide     = "hide"
	hideActionShow     = "show"
)

type Handler struct {
	loc *localization.LocalizationManager
	api *api.APIClient
	bot *tgbotapi.BotAPI
}

func NewHandler(loc *localization.LocalizationManager, apiClient *api.APIClient, bot *tgbotapi.BotAPI) *Handler {
	return &Handler{loc: loc, api: apiClient, bot: bot}
}

func (h *Handler) HandleCommand(ctx context.Context, log *logging.Logger, localizer *i18n.Localizer, msg *tgbotapi.Message) {
	if msg.From == nil {
		return
	}

	profile, err := h.api.GetProfile(ctx, msg.From.ID, msg.From.UserName)
	isHidden := false
	if err != nil {
		log.E("Failed to get user profile", logging.InnerError, err)
	} else {
		isHidden = profile.IsHidden
	}

	anonName := h.loc.Localize(localizer, localization.AnonymousNameTemplate, map[string]any{"Number": generateAnonymousNumber(msg.From.ID)})
	realName := msg.From.UserName
	if realName == "" {
		realName = anonName
	}

	text, keyboard := h.buildHideMessage(localizer, isHidden, anonName, realName, msg.From.ID)

	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = "MarkdownV2"
	if keyboard != nil {
		reply.ReplyMarkup = keyboard
	}
	if _, err := h.bot.Send(reply); err != nil {
		log.E("Failed to send /hide message", logging.InnerError, err)
	}
}

func (h *Handler) buildHideMessage(localizer *i18n.Localizer, isHidden bool, anonName, realName string, userID int64) (string, *tgbotapi.InlineKeyboardMarkup) {
	msgKey := localization.MsgHidePrompt
	buttonKey := localization.MsgHideButtonHide
	action := hideActionHide

	if isHidden {
		msgKey = localization.MsgHideStatusHidden
		buttonKey = localization.MsgHideButtonShow
		action = hideActionShow
	}

	text := h.loc.Localize(localizer, msgKey, map[string]any{
		"Anon": formatting.EscapeMarkdownV2(anonName),
		"Real": formatting.EscapeMarkdownV2(realName),
	})
	button := tgbotapi.NewInlineKeyboardButtonData(
		h.loc.Localize(localizer, buttonKey, nil),
		fmt.Sprintf("%s%d:%s", HideCallbackPrefix, userID, action),
	)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
	return text, &keyboard
}
