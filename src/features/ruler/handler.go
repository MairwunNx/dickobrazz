package ruler

import (
	"context"
	"dickobrazz/src/shared/localization"
	"dickobrazz/src/shared/logging"
	"dickobrazz/src/shared/telegram"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	action *GetAction
	loc    *localization.LocalizationManager
}

func NewHandler(action *GetAction, loc *localization.LocalizationManager) *Handler {
	return &Handler{action: action, loc: loc}
}

func (h *Handler) HandleInlineQuery(ctx context.Context, log *logging.Logger, update *tgbotapi.Update, showDescription bool) tgbotapi.InlineQueryResultArticle {
	query := update.InlineQuery
	if query == nil {
		return tgbotapi.InlineQueryResultArticle{}
	}
	localizer, _ := h.loc.LocalizerByUpdate(update)

	text, err := h.action.Execute(ctx, log, localizer, query.From.ID, query.From.UserName, showDescription)
	if err != nil {
		return tgbotapi.InlineQueryResultArticle{}
	}

	return telegram.InitializeInlineQueryWithThumbAndDesc(
		h.loc.Localize(localizer, "InlineTitleCockRuler", nil),
		text,
		h.loc.Localize(localizer, "DescCockRuler", nil),
		"https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_ruler.png",
	)
}
