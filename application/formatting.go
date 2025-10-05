package application

import (
	"dickobrazz/application/datetime"
	"dickobrazz/application/logging"
	"fmt"
	"strings"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func GenerateCockSizeText(size int, emoji string) string {
	displaySize := size
	
	now := datetime.NowTime()
	if now.Month() == time.April && now.Day() == 1 {
		displaySize = -size
	}
	
	return fmt.Sprintf(MsgCockSize, displaySize, emoji)
}

func (app *Application) GenerateCockRulerText(log *logging.Logger, userID int64, cocks []UserCock) string {
	var winners []string
	var others []string
	isUserInScoreboard := false

	for index, cock := range cocks {
		isCurrentUser := cock.UserId == userID
		emoji := GetPlaceEmoji(index + 1)

		var line string
		if isCurrentUser {
			isUserInScoreboard = true
			line = fmt.Sprintf(MsgCockRulerScoreboardSelected, emoji, EscapeMarkdownV2(cock.UserName), cock.Size, EmojiFromSize(cock.Size))
		} else {
			line = fmt.Sprintf(MsgCockRulerScoreboardDefault, emoji, EscapeMarkdownV2(cock.UserName), cock.Size, EmojiFromSize(cock.Size))
		}

		if index < 3 {
			winners = append(winners, line)
		} else {
			others = append(others, line)
		}
	}

	if !isUserInScoreboard {
		if userCock := app.GetCockSizeFromCache(log, userID); userCock != nil {
			others = append(others, fmt.Sprintf(MsgCockRulerScoreboardOut, EscapeMarkdownV2(CommonDots), *userCock, EmojiFromSize(*userCock)))
		} else {
			others = append(others, MsgCockScoreboardNotFound)
		}
	}

	if len(others) != 0 {
		return fmt.Sprintf(
			MsgCockRulerScoreboardTemplate,
			strings.Join(winners, "\n"),
			strings.Join(others, "\n"),
		)
	} else {
		return fmt.Sprintf(
			MsgCockRulerScoreboardWinnersTemplate,
			strings.Join(winners, "\n"),
		)
	}
}

func (app *Application) GenerateCockRaceScoreboard(log *logging.Logger, userID int64, sizes []UserCockRace, seasonStart string) string {
	var winners []string
	var others []string
	isUserInScoreboard := false

	for index, user := range sizes {
		isCurrentUser := user.UserID == userID
		emoji := GetPlaceEmoji(index + 1)

		if isCurrentUser {
			isUserInScoreboard = true
		}

		var scoreboardLine string
		if isCurrentUser {
			scoreboardLine = fmt.Sprintf(MsgCockRaceScoreboardSelected, emoji, EscapeMarkdownV2(user.Nickname), FormatDickSize(int(user.TotalSize)))
		} else {
			scoreboardLine = fmt.Sprintf(MsgCockRaceScoreboardDefault, emoji, EscapeMarkdownV2(user.Nickname), FormatDickSize(int(user.TotalSize)))
		}

		if index < 3 {
			winners = append(winners, scoreboardLine)
		} else {
			others = append(others, scoreboardLine)
		}
	}

	if !isUserInScoreboard {
		if cock := app.GetUserAggregatedCock(log, userID); cock != nil {
			others = append(others, fmt.Sprintf(MsgCockRaceScoreboardOut, EscapeMarkdownV2(cock.Nickname), FormatDickSize(int(cock.TotalSize))))
		} else {
			others = append(others, MsgCockScoreboardNotFound)
		}
	}

	if len(others) != 0 {
		return fmt.Sprintf(
			MsgCockRaceScoreboardTemplate,
			strings.Join(winners, "\n"),
			strings.Join(others, "\n"),
			seasonStart,
		)
	} else {
		return fmt.Sprintf(
			MsgCockRaceScoreboardWinnersTemplate,
			strings.Join(winners, "\n"),
			seasonStart,
		)
	}
}

func GetPlaceEmoji(place int) string {
	switch place {
	case 1:
		return "🥇"
	case 2:
		return "🥈"
	case 3:
		return "🥉"
	default:
		return "🤧"
	}
}

func EscapeMarkdownV2(input string) string {
	var str strings.Builder
	escapable := "_*[]()~`>#+-=|{}.!\\"
	for _, char := range input {
		if strings.ContainsRune(escapable, char) {
			str.WriteRune('\\')
		}
		str.WriteRune(char)
	}
	return str.String()
}

var p = message.NewPrinter(language.Russian)

func FormatDickPercent(size float64) string {
	return p.Sprintf("%.1f", size)
}

func FormatDickSize(size int) string {
	return p.Sprintf("%d", size)
}

func FormatDickIkr(ikr float64) string {
	return p.Sprintf("%.3f", ikr)
}

func FormatLuckCoefficient(luck float64) string {
	return p.Sprintf("%.3f", luck)
}

func FormatVolatility(volatility float64) string {
	return p.Sprintf("%.1f", volatility)
}

func LuckEmoji(luck float64) string {
	switch {
	case luck >= 1.2:
		return "🍀✨"
	case luck >= 1.1:
		return "🍀"
	case luck >= 0.9:
		return "⚖️"
	case luck >= 0.8:
		return "😔"
	default:
		return "💀"
	}
}

func VolatilityEmoji(volatility float64) string {
	switch {
	case volatility < 10:
		return "📈"
	case volatility < 20:
		return "📊"
	default:
		return "🎢"
	}
}
