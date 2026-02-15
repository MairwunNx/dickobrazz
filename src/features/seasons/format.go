package seasons

import (
	"dickobrazz/src/shared/api"
	"dickobrazz/src/shared/formatting"
	"dickobrazz/src/shared/localization"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func generateSeasonPageText(locMgr *localization.LocalizationManager, localizer *i18n.Localizer, season api.SeasonWithWinners, showDescription bool) string {
	startDate := formatting.EscapeMarkdownV2(season.StartDate.FormatDateMSK())
	endDate := formatting.EscapeMarkdownV2(season.EndDate.FormatDateMSK())

	var winnerLines []string
	for _, winner := range season.Winners {
		medal := formatting.GetMedalByPosition(winner.Place - 1)
		line := NewMsgCockSeasonWinnerTemplate(
			locMgr,
			localizer,
			medal,
			winner.Nickname,
			formatting.FormatDickSize(winner.TotalSize),
		)
		winnerLines = append(winnerLines, line)
	}

	winnersText := strings.Join(winnerLines, "\n")

	if season.IsActive {
		seasonBlock := NewMsgCockSeasonTemplate(locMgr, localizer, winnersText, startDate, endDate, season.SeasonNum)
		if showDescription {
			footer := NewMsgCockSeasonTemplateFooter(locMgr, localizer)
			return seasonBlock + "\n\n" + footer
		}
		return seasonBlock
	}

	return NewMsgCockSeasonWithWinnersTemplate(locMgr, localizer, winnersText, startDate, endDate, season.SeasonNum)
}
