package seasons

import (
	"dickobrazz/src/shared/formatting"
	"dickobrazz/src/shared/localization"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func NewMsgCockSeasonTemplate(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, pretenders string, startDate, endDate string, seasonNum int) string {
	return localizationManager.Localize(localizer, "MsgCockSeasonTemplate", map[string]any{
		"Pretenders": pretenders,
		"StartDate":  startDate,
		"EndDate":    endDate,
		"SeasonNum":  seasonNum,
	})
}

func NewMsgCockSeasonWithWinnersTemplate(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, winners string, startDate, endDate string, seasonNum int) string {
	return localizationManager.Localize(localizer, "MsgCockSeasonWithWinnersTemplate", map[string]any{
		"Winners":   winners,
		"StartDate": startDate,
		"EndDate":   endDate,
		"SeasonNum": seasonNum,
	})
}

func NewMsgCockSeasonWinnerTemplate(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, medal, nickname, totalSize string) string {
	return localizationManager.Localize(localizer, "MsgCockSeasonWinnerTemplate", map[string]any{
		"Medal":    medal,
		"Nickname": formatting.EscapeMarkdownV2(nickname),
		"Size":     formatting.EscapeMarkdownV2(totalSize),
	})
}

func NewMsgCockSeasonTemplateFooter(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer) string {
	return localizationManager.Localize(localizer, "MsgCockSeasonTemplateFooter", nil)
}

func NewMsgCockSeasonNoSeasonsTemplate(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer) string {
	return localizationManager.Localize(localizer, "MsgCockSeasonNoSeasonsTemplate", nil)
}
