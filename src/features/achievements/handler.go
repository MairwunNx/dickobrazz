package achievements

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

func (h *Handler) HandleInlineQuery(ctx context.Context, log *logging.Logger, update *tgbotapi.Update, page int) tgbotapi.InlineQueryResultArticle {
	query := update.InlineQuery
	if query == nil {
		return tgbotapi.InlineQueryResultArticle{}
	}
	localizer, _ := h.loc.LocalizerByUpdate(update)

	result, err := h.action.Execute(ctx, log, localizer, query.From.ID, query.From.UserName, page)
	if err != nil {
		return tgbotapi.InlineQueryResultArticle{}
	}

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(result.Buttons...),
	)

	article := tgbotapi.NewInlineQueryResultArticleMarkdownV2(
		fmt.Sprintf("ach_%d_%d", query.From.ID, page),
		h.loc.Localize(localizer, "InlineTitleCockAchievements", nil),
		result.Text,
	)
	article.ReplyMarkup = &kb
	article.ThumbURL = "https://files.mairwunnx.com/raw/public/dickobrazz%2Fico_achievements.png"
	article.Description = h.loc.Localize(localizer, "DescCockAchievements", nil)

	return article
}
