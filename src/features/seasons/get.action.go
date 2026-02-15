package seasons

import (
	"context"
	"dickobrazz/src/shared/api"
	"dickobrazz/src/shared/localization"
	"dickobrazz/src/shared/logging"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type SeasonResult struct {
	Text      string
	SeasonNum int
	Buttons   []tgbotapi.InlineKeyboardButton
}

type GetAction struct {
	api *api.APIClient
	loc *localization.LocalizationManager
}

func NewGetAction(apiClient *api.APIClient, loc *localization.LocalizationManager) *GetAction {
	return &GetAction{api: apiClient, loc: loc}
}

func (a *GetAction) Execute(ctx context.Context, log *logging.Logger, localizer *i18n.Localizer, userID int64, username string, showDescription bool) (*SeasonResult, error) {
	seasonsData, err := a.api.GetCockSeasons(ctx, userID, username, 15, 1)
	if err != nil {
		log.E("Failed to get cock seasons via API", logging.InnerError, err)
		return nil, err
	}

	if len(seasonsData.Seasons) == 0 {
		return nil, nil
	}

	currentSeason := seasonsData.Seasons[0]
	text := generateSeasonPageText(a.loc, localizer, currentSeason, showDescription)

	var buttons []tgbotapi.InlineKeyboardButton
	if len(seasonsData.Seasons) > 1 {
		prevSeason := seasonsData.Seasons[1]
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(
			a.loc.Localize(localizer, localization.MsgSeasonButton, map[string]any{
				"Arrow":     "◀️",
				"SeasonNum": prevSeason.SeasonNum,
			}),
			fmt.Sprintf("season_page:%d", prevSeason.SeasonNum),
		))
	}

	return &SeasonResult{
		Text:      text,
		SeasonNum: currentSeason.SeasonNum,
		Buttons:   buttons,
	}, nil
}
