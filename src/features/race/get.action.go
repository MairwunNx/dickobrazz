package race

import (
	"context"
	"dickobrazz/src/shared/api"
	"dickobrazz/src/shared/localization"
	"dickobrazz/src/shared/logging"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type GetAction struct {
	api *api.APIClient
	loc *localization.LocalizationManager
}

func NewGetAction(apiClient *api.APIClient, loc *localization.LocalizationManager) *GetAction {
	return &GetAction{api: apiClient, loc: loc}
}

func (a *GetAction) Execute(ctx context.Context, log *logging.Logger, localizer *i18n.Localizer, userID int64, username string, showDescription bool) (string, error) {
	data, err := a.api.GetCockRace(ctx, userID, username, 13, 1)
	if err != nil {
		log.E("Failed to get cock race via API", logging.InnerError, err)
		return "", err
	}

	text := GenerateCockRaceScoreboard(a.loc, localizer, userID, data, showDescription)
	return text, nil
}
