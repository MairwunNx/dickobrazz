package help

import (
	"dickobrazz/src/shared/localization"
	"dickobrazz/src/shared/logging"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Handler struct {
	loc *localization.LocalizationManager
	bot *tgbotapi.BotAPI
}

func NewHandler(loc *localization.LocalizationManager, bot *tgbotapi.BotAPI) *Handler {
	return &Handler{loc: loc, bot: bot}
}

func (h *Handler) HandleCommand(log *logging.Logger, localizer *i18n.Localizer, msg *tgbotapi.Message) {
	text := h.loc.Localize(localizer, localization.MsgHelpText, nil)
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = "MarkdownV2"
	if _, err := h.bot.Send(reply); err != nil {
		log.E("Failed to send /help message", logging.InnerError, err)
	}
}
