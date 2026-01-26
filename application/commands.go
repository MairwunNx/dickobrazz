package application

import (
	"dickobrazz/application/logging"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	hideCallbackPrefix = "hide_toggle:"
	hideActionHide     = "hide"
	hideActionShow     = "show"
)

func (app *Application) HandleCommand(log *logging.Logger, update *tgbotapi.Update) {
	msg := update.Message
	if msg == nil || !msg.IsCommand() {
		return
	}
	localizer, _ := app.localization.LocalizerByUpdate(update)

	switch msg.Command() {
	case "help":
		app.sendHelpMessage(log, localizer, msg)
	case "hide":
		app.sendHideMessage(log, localizer, msg)
	}
}

func (app *Application) ResolveUserNickname(log *logging.Logger, localizer *i18n.Localizer, user *tgbotapi.User) string {
	if user == nil {
		return ""
	}
	normalized := app.ResolveDisplayNickname(log, localizer, user.ID, user.UserName)
	app.UpsertUserProfile(log, user.ID, user.UserName, app.IsUserHidden(log, user.ID))
	return normalized
}

func (app *Application) IsUserHidden(log *logging.Logger, userID int64) bool {
	profile := app.GetUserProfile(log, userID)
	return profile != nil && profile.IsHidden
}

func (app *Application) ResolveDisplayNickname(log *logging.Logger, localizer *i18n.Localizer, userID int64, username string) string {
	if app.IsUserHidden(log, userID) {
		return GenerateAnonymousName(app.localization, localizer, userID)
	}
	return NormalizeUsername(app.localization, localizer, username, userID)
}

func (app *Application) sendHelpMessage(log *logging.Logger, localizer *i18n.Localizer, msg *tgbotapi.Message) {
	text := app.localization.Localize(localizer, MsgHelpText, nil)
	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = "MarkdownV2"
	if _, err := app.bot.Send(reply); err != nil {
		log.E("Failed to send /help message", logging.InnerError, err)
	}
}

func (app *Application) sendHideMessage(log *logging.Logger, localizer *i18n.Localizer, msg *tgbotapi.Message) {
	if msg.From == nil {
		return
	}
	isHidden := app.IsUserHidden(log, msg.From.ID)
	anonName := GenerateAnonymousName(app.localization, localizer, msg.From.ID)
	realName := NormalizeUsername(app.localization, localizer, msg.From.UserName, msg.From.ID)
	text, keyboard := app.buildHideMessage(localizer, isHidden, anonName, realName, msg.From.ID)

	reply := tgbotapi.NewMessage(msg.Chat.ID, text)
	reply.ParseMode = "MarkdownV2"
	if keyboard != nil {
		reply.ReplyMarkup = keyboard
	}
	if _, err := app.bot.Send(reply); err != nil {
		log.E("Failed to send /hide message", logging.InnerError, err)
	}
}

func (app *Application) buildHideMessage(localizer *i18n.Localizer, isHidden bool, anonName, realName string, userID int64) (string, *tgbotapi.InlineKeyboardMarkup) {
	msgKey := MsgHidePrompt
	buttonKey := MsgHideButtonHide
	action := hideActionHide

	if isHidden {
		msgKey = MsgHideStatusHidden
		buttonKey = MsgHideButtonShow
		action = hideActionShow
	}

	text := app.localization.Localize(localizer, msgKey, map[string]any{
		"Anon": EscapeMarkdownV2(anonName),
		"Real": EscapeMarkdownV2(realName),
	})
	button := tgbotapi.NewInlineKeyboardButtonData(
		app.localization.Localize(localizer, buttonKey, nil),
		buildHideCallbackData(userID, action),
	)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(button))
	return text, &keyboard
}

func (app *Application) setUserHiddenStatus(log *logging.Logger, localizer *i18n.Localizer, user *tgbotapi.User, isHidden bool) (string, string) {
	if user == nil {
		return "", ""
	}
	anonName := GenerateAnonymousName(app.localization, localizer, user.ID)
	realName := NormalizeUsername(app.localization, localizer, user.UserName, user.ID)

	app.UpsertUserProfile(log, user.ID, user.UserName, isHidden)

	return anonName, realName
}

func buildHideCallbackData(userID int64, action string) string {
	return fmt.Sprintf("%s%d:%s", hideCallbackPrefix, userID, action)
}
