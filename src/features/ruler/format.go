package ruler

import (
	"dickobrazz/src/shared/api"
	"dickobrazz/src/shared/emoji"
	"dickobrazz/src/shared/formatting"
	"dickobrazz/src/shared/localization"
	"fmt"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func GenerateCockRulerText(loc *localization.LocalizationManager, localizer *i18n.Localizer, userID int64, data *api.CockRulerData, showDescription bool) string {
	var winners []string
	var others []string
	isUserInScoreboard := false

	for index, entry := range data.Leaders {
		isCurrentUser := entry.UserID == userID
		placeEmoji := formatting.GetPlaceEmoji(index+1, isCurrentUser)
		formattedSize := formatting.FormatCockSizeForDate(entry.Size)

		var line string
		if isCurrentUser {
			isUserInScoreboard = true
			line = loc.Localize(localizer, localization.MsgCockRulerScoreboardSelected, map[string]any{
				"PlaceEmoji": placeEmoji,
				"Username":   formatting.EscapeMarkdownV2(entry.Nickname),
				"Size":       formattedSize,
				"SizeEmoji":  emoji.EmojiFromSize(entry.Size),
			})
		} else {
			line = loc.Localize(localizer, localization.MsgCockRulerScoreboardDefault, map[string]any{
				"PlaceEmoji": placeEmoji,
				"Username":   formatting.EscapeMarkdownV2(entry.Nickname),
				"Size":       formattedSize,
				"SizeEmoji":  emoji.EmojiFromSize(entry.Size),
			})
		}

		if index < 3 {
			winners = append(winners, line)
		} else {
			others = append(others, line)
		}
	}

	if !isUserInScoreboard && data.UserPosition != nil {
		neighborhood := data.Neighborhood
		var contextLines []string

		for _, neighbor := range neighborhood.Above {
			contextLines = append(contextLines, loc.Localize(localizer, localization.MsgCockRulerContextDefault, map[string]any{
				"PlaceEmoji": "🥀",
				"Username":   formatting.EscapeMarkdownV2(neighbor.Nickname),
				"Size":       formatting.EscapeMarkdownV2(formatting.FormatCockSizeForDate(neighbor.Size)),
				"SizeEmoji":  emoji.EmojiFromSize(neighbor.Size),
			}))
		}

		if neighborhood.Self != nil {
			self := neighborhood.Self
			contextLines = append(contextLines, loc.Localize(localizer, localization.MsgCockRulerContextSelected, map[string]any{
				"PlaceEmoji": fmt.Sprintf("🥀 *%d*\\.", *data.UserPosition),
				"Username":   formatting.EscapeMarkdownV2(self.Nickname),
				"Size":       formatting.EscapeMarkdownV2(formatting.FormatCockSizeForDate(self.Size)),
				"SizeEmoji":  emoji.EmojiFromSize(self.Size),
			}))
		}

		for _, neighbor := range neighborhood.Below {
			contextLines = append(contextLines, loc.Localize(localizer, localization.MsgCockRulerContextDefault, map[string]any{
				"PlaceEmoji": "🥀",
				"Username":   formatting.EscapeMarkdownV2(neighbor.Nickname),
				"Size":       formatting.EscapeMarkdownV2(formatting.FormatCockSizeForDate(neighbor.Size)),
				"SizeEmoji":  emoji.EmojiFromSize(neighbor.Size),
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
		template := localization.MsgCockRulerScoreboardTemplate
		if !showDescription {
			template = localization.MsgCockRulerScoreboardTemplateNoDesc
		}
		return loc.Localize(localizer, template, map[string]any{
			"Participants": data.TotalParticipants,
			"Winners":      strings.Join(winners, "\n"),
			"Others":       strings.Join(others, "\n"),
		})
	}
	template := localization.MsgCockRulerScoreboardWinnersTemplate
	if !showDescription {
		template = localization.MsgCockRulerScoreboardWinnersTemplateNoDesc
	}
	return loc.Localize(localizer, template, map[string]any{
		"Participants": data.TotalParticipants,
		"Winners":      strings.Join(winners, "\n"),
	})
}
