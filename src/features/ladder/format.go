package ladder

import (
	"dickobrazz/src/shared/api"
	"dickobrazz/src/shared/formatting"
	"dickobrazz/src/shared/localization"
	"fmt"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func GenerateCockLadderScoreboard(loc *localization.LocalizationManager, localizer *i18n.Localizer, userID int64, data *api.CockLadderData, showDescription bool) string {
	var winners []string
	var others []string
	isUserInScoreboard := false

	for index, entry := range data.Leaders {
		isCurrentUser := entry.UserID == userID
		placeEmoji := formatting.GetPlaceEmoji(index+1, isCurrentUser)

		if isCurrentUser {
			isUserInScoreboard = true
		}

		var scoreboardLine string
		if isCurrentUser {
			scoreboardLine = loc.Localize(localizer, localization.MsgCockLadderScoreboardSelected, map[string]any{
				"PlaceEmoji": placeEmoji,
				"Username":   formatting.EscapeMarkdownV2(entry.Nickname),
				"Size":       formatting.EscapeMarkdownV2(formatting.FormatDickSize(entry.TotalSize)),
			})
		} else {
			scoreboardLine = loc.Localize(localizer, localization.MsgCockLadderScoreboardDefault, map[string]any{
				"PlaceEmoji": placeEmoji,
				"Username":   formatting.EscapeMarkdownV2(entry.Nickname),
				"Size":       formatting.EscapeMarkdownV2(formatting.FormatDickSize(entry.TotalSize)),
			})
		}

		if index < 3 {
			winners = append(winners, scoreboardLine)
		} else {
			others = append(others, scoreboardLine)
		}
	}

	if !isUserInScoreboard && data.UserPosition != nil {
		neighborhood := data.Neighborhood
		var contextLines []string

		for _, neighbor := range neighborhood.Above {
			contextLines = append(contextLines, loc.Localize(localizer, localization.MsgCockLadderContextDefault, map[string]any{
				"PlaceEmoji": "🥀",
				"Username":   formatting.EscapeMarkdownV2(neighbor.Nickname),
				"Size":       formatting.EscapeMarkdownV2(formatting.FormatDickSize(neighbor.TotalSize)),
			}))
		}

		if neighborhood.Self != nil {
			self := neighborhood.Self
			contextLines = append(contextLines, loc.Localize(localizer, localization.MsgCockLadderContextSelected, map[string]any{
				"PlaceEmoji": fmt.Sprintf("🥀 *%d*\\.", *data.UserPosition),
				"Username":   formatting.EscapeMarkdownV2(self.Nickname),
				"Size":       formatting.EscapeMarkdownV2(formatting.FormatDickSize(self.TotalSize)),
			}))
		}

		for _, neighbor := range neighborhood.Below {
			contextLines = append(contextLines, loc.Localize(localizer, localization.MsgCockLadderContextDefault, map[string]any{
				"PlaceEmoji": "🥀",
				"Username":   formatting.EscapeMarkdownV2(neighbor.Nickname),
				"Size":       formatting.EscapeMarkdownV2(formatting.FormatDickSize(neighbor.TotalSize)),
			}))
		}

		if len(contextLines) > 0 {
			dots := loc.Localize(localizer, localization.CommonDots, nil)
			showTopDots := *data.UserPosition > len(data.Leaders)+1
			showBottomDots := *data.UserPosition < data.TotalParticipants

			var contextBlock string
			if showTopDots && showBottomDots {
				contextBlock = "\n" + dots + "\n" + strings.Join(contextLines, "\n") + "\n" + dots
			} else if showTopDots {
				contextBlock = "\n" + dots + "\n" + strings.Join(contextLines, "\n")
			} else if showBottomDots {
				contextBlock = "\n" + strings.Join(contextLines, "\n") + "\n" + dots
			} else {
				contextBlock = "\n" + strings.Join(contextLines, "\n")
			}
			others = append(others, contextBlock)
		}
	} else if !isUserInScoreboard {
		others = append(others, loc.Localize(localizer, localization.MsgCockScoreboardNotFound, nil))
	}

	if len(others) != 0 {
		template := localization.MsgCockLadderScoreboardTemplate
		if !showDescription {
			template = localization.MsgCockLadderScoreboardTemplateNoDesc
		}
		return loc.Localize(localizer, template, map[string]any{
			"Participants": data.TotalParticipants,
			"Winners":      strings.Join(winners, "\n"),
			"Others":       strings.Join(others, "\n"),
		})
	}
	template := localization.MsgCockLadderScoreboardWinnersTemplate
	if !showDescription {
		template = localization.MsgCockLadderScoreboardWinnersTemplateNoDesc
	}
	return loc.Localize(localizer, template, map[string]any{
		"Participants": data.TotalParticipants,
		"Winners":      strings.Join(winners, "\n"),
	})
}
