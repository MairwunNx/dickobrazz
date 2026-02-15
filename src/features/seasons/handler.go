package seasons

import (
	"context"
	"dickobrazz/src/shared/localization"
	"dickobrazz/src/shared/logging"
	"fmt"

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

	result, err := h.action.Execute(ctx, log, localizer, query.From.ID, query.From.UserName, showDescription)
	if err != nil || result == nil {
		text := NewMsgCockSeasonNoSeasonsTemplate(h.loc, localizer)
		article := tgbotapi.NewInlineQueryResultArticleMarkdownV2(
			"season_err",
			h.loc.Localize(localizer, localization.InlineTitleCockSeason, nil),
			text,
		)
		article.ThumbURL = "https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_seasons.png"
		article.Description = h.loc.Localize(localizer, localization.DescCockSeason, nil)
		return article
	}

	article := tgbotapi.NewInlineQueryResultArticleMarkdownV2(
		fmt.Sprintf("season_%d", result.SeasonNum),
		h.loc.Localize(localizer, localization.InlineTitleCockSeason, nil),
		result.Text,
	)

	if len(result.Buttons) > 0 {
		kb := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(result.Buttons...),
		)
		article.ReplyMarkup = &kb
	}
	article.ThumbURL = "https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_seasons.png"
	article.Description = h.loc.Localize(localizer, localization.DescCockSeason, nil)

	return article
}
