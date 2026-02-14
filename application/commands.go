package application

import (
	"dickobrazz/application/logging"
	"fmt"
	"math/big"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// LCG multiplier (Knuth), used for deterministic anonymous number from user ID.
const anonNumberMultiplier = 6364136223846793005

// generateAnonymousNumber produces a deterministic 4-digit anonymous number from user ID.
// Algorithm matches: ((id * multiplier) & 0xffffffff) % 10000.
func generateAnonymousNumber(userID int64) string {
	id := big.NewInt(userID)
	multiplier := big.NewInt(anonNumberMultiplier)
	mask := big.NewInt(0xffffffff)
	mod := big.NewInt(10000)

	var n big.Int
	n.Mul(id, multiplier)
	n.And(&n, mask)
	n.Mod(&n, mod)

	v := n.Int64()
	if v < 0 {
		v = -v
	}
	return fmt.Sprintf("%04d", v%10000)
}

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

	// Получаем профиль через API
	profile, err := app.api.GetProfile(app.ctx, msg.From.ID, msg.From.UserName)
	isHidden := false
	if err != nil {
		log.E("Failed to get user profile", logging.InnerError, err)
	} else {
		isHidden = profile.IsHidden
	}

	anonName := app.localization.Localize(localizer, AnonymousNameTemplate, map[string]any{"Number": generateAnonymousNumber(msg.From.ID)})
	realName := msg.From.UserName
	if realName == "" {
		anonName = app.localization.Localize(localizer, AnonymousNameTemplate, map[string]any{"Number": generateAnonymousNumber(msg.From.ID)})
		realName = anonName
	}

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
	anonName := app.localization.Localize(localizer, AnonymousNameTemplate, map[string]any{"Number": generateAnonymousNumber(user.ID)})
	realName := user.UserName
	if realName == "" {
		realName = anonName
	}

	if _, err := app.api.UpdatePrivacy(app.ctx, user.ID, user.UserName, isHidden); err != nil {
		log.E("Failed to update user privacy via API", logging.InnerError, err)
	}

	return anonName, realName
}

func buildHideCallbackData(userID int64, action string) string {
	return fmt.Sprintf("%s%d:%s", hideCallbackPrefix, userID, action)
}
