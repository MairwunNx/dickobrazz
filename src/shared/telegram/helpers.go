package telegram

import (
	"dickobrazz/src/shared/logging"
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func InitializeInlineQueryWithThumbAndDesc(title, message, description, thumbURL string) tgbotapi.InlineQueryResultArticle {
	article := tgbotapi.NewInlineQueryResultArticleMarkdownV2(fmt.Sprintf("q_%d", time.Now().UnixNano()), title, message)
	article.ThumbURL = thumbURL
	article.Description = description
	return article
}

func EditCallbackMessage(log *logging.Logger, bot *tgbotapi.BotAPI, callback *tgbotapi.CallbackQuery, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	if callback.InlineMessageID != "" {
		edit := tgbotapi.EditMessageTextConfig{
			BaseEdit: tgbotapi.BaseEdit{
				InlineMessageID: callback.InlineMessageID,
			},
			Text:      text,
			ParseMode: "MarkdownV2",
		}
		if keyboard != nil {
			edit.ReplyMarkup = keyboard
		}
		if _, err := bot.Request(edit); err != nil {
			log.E("Failed to edit inline message", logging.InnerError, err)
		}
	} else if callback.Message != nil {
		edit := tgbotapi.NewEditMessageText(
			callback.Message.Chat.ID,
			callback.Message.MessageID,
			text,
		)
		edit.ParseMode = "MarkdownV2"
		if keyboard != nil {
			edit.ReplyMarkup = keyboard
		}
		if _, err := bot.Request(edit); err != nil {
			log.E("Failed to edit chat message", logging.InnerError, err)
		}
	} else {
		log.E("CallbackQuery has neither Message nor InlineMessageID")
	}
}
