package achievements

import (
	"context"
	"dickobrazz/src/entities/achievement"
	"dickobrazz/src/shared/api"
	"dickobrazz/src/shared/localization"
	"dickobrazz/src/shared/logging"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type AchievementsResult struct {
	Text    string
	Buttons []tgbotapi.InlineKeyboardButton
}

type GetAction struct {
	api *api.APIClient
	loc *localization.LocalizationManager
}

func NewGetAction(apiClient *api.APIClient, loc *localization.LocalizationManager) *GetAction {
	return &GetAction{api: apiClient, loc: loc}
}

func (a *GetAction) Execute(ctx context.Context, log *logging.Logger, localizer *i18n.Localizer, userID int64, username string, page int) (*AchievementsResult, error) {
	type achResult struct {
		data *api.CockAchievementsData
		err  error
	}
	type respResult struct {
		data *api.RespectData
		err  error
	}

	achCh := make(chan achResult, 1)
	respCh := make(chan respResult, 1)

	go func() {
		data, err := a.api.GetCockAchievements(ctx, userID, username)
		achCh <- achResult{data, err}
	}()
	go func() {
		data, err := a.api.GetCockRespects(ctx, userID, username)
		respCh <- respResult{data, err}
	}()

	achRes := <-achCh
	respRes := <-respCh

	if achRes.err != nil {
		log.E("Failed to get achievements via API", logging.InnerError, achRes.err)
		return nil, achRes.err
	}
	achData := achRes.data

	achievementRespects := 0
	if respRes.err == nil && respRes.data != nil {
		achievementRespects = int(respRes.data.AchievementRespect)
	}

	achievementsList := GenerateAchievementsText(
		a.loc, localizer,
		achievement.AllAchievements,
		achData.Achievements,
		page, 10,
	)

	totalAchievements := achData.AchievementsTotal
	totalPages := (totalAchievements + 9) / 10

	var templateID string
	if page == 1 {
		templateID = localization.MsgCockAchievementsTemplate
	} else {
		templateID = localization.MsgCockAchievementsTemplateOtherPages
	}

	text := a.loc.Localize(localizer, templateID, map[string]any{
		"Completed":    achData.AchievementsDone,
		"Total":        totalAchievements,
		"Percent":      int(achData.AchievementsDonePercent),
		"Respects":     achievementRespects,
		"Achievements": achievementsList,
	})

	var buttons []tgbotapi.InlineKeyboardButton

	if page > 1 {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("◀️", fmt.Sprintf("ach_page:%d:%d", userID, page-1)))
	}

	buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%d/%d", page, totalPages), "ach_noop"))

	if page < totalPages {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("▶️", fmt.Sprintf("ach_page:%d:%d", userID, page+1)))
	}

	return &AchievementsResult{Text: text, Buttons: buttons}, nil
}
