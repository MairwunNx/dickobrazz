package dynamics

import (
	"context"
	"dickobrazz/src/shared/localization"
	"dickobrazz/src/shared/logging"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler struct {
	action *GetAction
	loc    *localization.LocalizationManager
}

func NewHandler(action *GetAction, loc *localization.LocalizationManager) *Handler {
	return &Handler{action: action, loc: loc}
}

func (h *Handler) HandleInlineQuery(ctx context.Context, log *logging.Logger, update *tgbotapi.Update) tgbotapi.InlineQueryResultArticle {
	query := update.InlineQuery
	if query == nil {
		return tgbotapi.InlineQueryResultArticle{}
	}
	localizer, _ := h.loc.LocalizerByUpdate(update)

	text, err := h.action.Execute(ctx, log, localizer, query.From.ID, query.From.UserName)
	if err != nil {
		return tgbotapi.InlineQueryResultArticle{}
	}

	article := tgbotapi.NewInlineQueryResultArticleMarkdown(query.ID, h.loc.Localize(localizer, localization.InlineTitleCockDynamic, nil), text)
	article.ThumbURL = "https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_dynamic.png"
	article.Description = h.loc.Localize(localizer, localization.DescCockDynamic, nil)
	return article
}
