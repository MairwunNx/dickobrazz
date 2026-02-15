package size

import (
	"context"
	"dickobrazz/src/shared/localization"
	"dickobrazz/src/shared/logging"
	"dickobrazz/src/shared/telegram"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	action *GenerateAction
	loc    *localization.LocalizationManager
}

func NewHandler(action *GenerateAction, loc *localization.LocalizationManager) *Handler {
	return &Handler{action: action, loc: loc}
}

func (h *Handler) HandleInlineQuery(ctx context.Context, log *logging.Logger, update *tgbotapi.Update) tgbotapi.InlineQueryResultArticle {
	query := update.InlineQuery
	if query == nil {
		return tgbotapi.InlineQueryResultArticle{}
	}
	localizer, _ := h.loc.LocalizerByUpdate(update)

	result, err := h.action.Execute(ctx, log, localizer, query.From.ID, query.From.UserName)
	if err != nil {
		return tgbotapi.InlineQueryResultArticle{}
	}

	text := result.Text + "\n\n" + "_" + result.Subtext + "_"

	return telegram.InitializeInlineQueryWithThumbAndDesc(
		h.loc.Localize(localizer, "InlineTitleCockSize", nil),
		strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(text, ".", "\\."), "-", "\\-"), "!", "\\!"),
		h.loc.Localize(localizer, "DescCockSize", nil),
		"https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_size.png",
	)
}
