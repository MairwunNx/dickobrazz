package privacy

import (
	"context"
	"dickobrazz/src/shared/api"
	"dickobrazz/src/shared/localization"
	"dickobrazz/src/shared/logging"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type ToggleAction struct {
	api *api.APIClient
	loc *localization.LocalizationManager
}

func NewToggleAction(apiClient *api.APIClient, loc *localization.LocalizationManager) *ToggleAction {
	return &ToggleAction{api: apiClient, loc: loc}
}

func (a *ToggleAction) Execute(ctx context.Context, log *logging.Logger, localizer *i18n.Localizer, user *tgbotapi.User, isHidden bool) (anonName string, realName string) {
	if user == nil {
		return "", ""
	}
	anonName = a.loc.Localize(localizer, localization.AnonymousNameTemplate, map[string]any{"Number": generateAnonymousNumber(user.ID)})
	realName = user.UserName
	if realName == "" {
		realName = anonName
	}

	if _, err := a.api.UpdatePrivacy(ctx, user.ID, user.UserName, isHidden); err != nil {
		log.E("Failed to update user privacy via API", logging.InnerError, err)
	}

	return anonName, realName
}
