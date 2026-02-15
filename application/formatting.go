package application

import (
	"dickobrazz/application/api"
	"dickobrazz/application/datetime"
	"dickobrazz/application/localization"
	"dickobrazz/application/logging"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"math"
	"math/rand"
	"sort"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var glitchMarks = []rune{
	'\u0335', '\u0336', '\u0337', '\u0338',
	'\u0300', '\u0301', '\u0302', '\u0303',
	'\u0304', '\u0305', '\u0306', '\u0307',
	'\u0308', '\u0309', '\u030A', '\u030B',
	'\u0310', '\u0311', '\u0312', '\u0313',
	'\u0334', '\u034F', '\u0350', '\u0351',
	'\u0352', '\u0353', '\u0354', '\u0355', '\u0356',
}

var mathFancy = map[int]string{
	0: "sin(0)", 1: "0!", 2: "C(2,1)", 3: "1! + 2!", 4: "2²", 5: "√25",
	6: "3!", 7: "3! + 1", 8: "2³", 9: "3²", 10: "C(5,2)", 11: "(1011)₂",
	12: "4! / 2", 13: "F₇", 14: "Cat₄", 15: "C(6,2)", 16: "2⁴", 17: "√289",
	18: "3! · 3", 19: "3³ − 2³", 20: "5! / 6", 21: "F₈", 22: "⌊π^e⌋", 23: "⌈π^e⌉",
	24: "4!", 25: "5²", 26: "4! + 2!", 27: "3³", 28: "T₇ = 7·8/2", 29: "2⁵ − 3",
	30: "2 · 5!!", 31: "2⁵ − 1", 32: "2⁵", 33: "4! + 3! + 2! + 0!", 34: "F₉",
	35: "C(7,3)", 36: "6²", 37: "⌊12π⌋", 38: "(100110)₂", 39: "3³ + 2·3!",
	40: "5! / 3", 41: "n² + n + 41 |_{n=0}", 42: "Cat₅", 43: "⌊14π⌋", 44: "⌊√2000⌋",
	45: "C(10,2)", 46: "4! + 4! − 2!", 47: "⌊15π⌋", 48: "4! · 2", 49: "7²",
	50: "⌊16π⌋", 51: "4! + 3³", 52: "6!! + 2²", 53: "⌊17π⌋", 54: "3³ + 3³",
	55: "F₁₀", 56: "C(8,3)", 57: "4! + 3! + 3³", 58: "6!! + C(5,2)", 59: "⌊19π⌋",
	60: "5! / 2", 61: "√3721",
}

var (
	rnd   = rand.New(rand.NewSource(time.Now().UnixNano()))
	rndMu sync.Mutex
)

func isMathDay(t time.Time) bool {
	return t.Month() == time.March && t.Day() == 14
}

func isProgrammersDay(t time.Time) bool {
	return t.YearDay() == 256
}

func toProgrammersNotation(n int) string {
	rndMu.Lock()
	useBinary := rnd.Intn(2) == 0
	rndMu.Unlock()

	if useBinary {
		if n < 0 {
			return "-0b" + strconv.FormatUint(uint64(-n), 2)
		}
		return "0b" + strconv.FormatUint(uint64(n), 2)
	}
	if n < 0 {
		return fmt.Sprintf("-0x%X", -n)
	}
	return fmt.Sprintf("0x%X", n)
}

func glitchify(s string) string {
	var sb strings.Builder
	for _, ch := range s {
		sb.WriteRune(ch)
		rndMu.Lock()
		count := rnd.Intn(3) + 1
		marks := make([]rune, count)
		for i := 0; i < count; i++ {
			marks[i] = glitchMarks[rnd.Intn(len(glitchMarks))]
		}
		rndMu.Unlock()

		for _, mark := range marks {
			sb.WriteRune(mark)
		}
	}
	return sb.String()
}

func fancyMathOrDefault(n int) string {
	if s, ok := mathFancy[n]; ok {
		return s
	}
	return strconv.Itoa(n)
}

// FormatCockSizeForDate форматирует размер в зависимости от текущей даты
func FormatCockSizeForDate(size int) string {
	displaySize := size
	now := datetime.NowTime()

	if now.Month() == time.April && now.Day() == 1 {
		displaySize = -size
	}

	if isMathDay(now) {
		return fancyMathOrDefault(displaySize)
	}

	if isProgrammersDay(now) {
		return toProgrammersNotation(displaySize)
	}

	if now.Month() == time.October && now.Day() == 31 {
		return glitchify(strconv.Itoa(displaySize))
	}

	return strconv.Itoa(displaySize)
}

func GenerateCockSizeText(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, size int, emoji string) string {
	formattedSize := FormatCockSizeForDate(size)
	return localizationManager.Localize(localizer, MsgCockSize, map[string]any{
		"Size":  formattedSize,
		"Emoji": emoji,
	})
}

func (app *Application) GenerateCockRulerText(log *logging.Logger, localizer *i18n.Localizer, userID int64, data *api.CockRulerData, showDescription bool) string {
	var winners []string
	var others []string
	isUserInScoreboard := false

	for index, entry := range data.Leaders {
		isCurrentUser := entry.UserID == userID
		emoji := GetPlaceEmoji(index+1, isCurrentUser)
		formattedSize := FormatCockSizeForDate(entry.Size)

		var line string
		if isCurrentUser {
			isUserInScoreboard = true
			line = app.localization.Localize(localizer, MsgCockRulerScoreboardSelected, map[string]any{
				"PlaceEmoji": emoji,
				"Username":   EscapeMarkdownV2(entry.Nickname),
				"Size":       formattedSize,
				"SizeEmoji":  EmojiFromSize(entry.Size),
			})
		} else {
			line = app.localization.Localize(localizer, MsgCockRulerScoreboardDefault, map[string]any{
				"PlaceEmoji": emoji,
				"Username":   EscapeMarkdownV2(entry.Nickname),
				"Size":       formattedSize,
				"SizeEmoji":  EmojiFromSize(entry.Size),
			})
		}

		if index < 3 {
			winners = append(winners, line)
		} else {
			others = append(others, line)
		}
	}

	// Показываем секцию соседей если пользователя нет в leaders
	if !isUserInScoreboard && data.UserPosition != nil {
		neighborhood := data.Neighborhood
		var contextLines []string

		for _, neighbor := range neighborhood.Above {
			contextLines = append(contextLines, app.localization.Localize(localizer, MsgCockRulerContextDefault, map[string]any{
				"PlaceEmoji": "🥀",
				"Username":   EscapeMarkdownV2(neighbor.Nickname),
				"Size":       EscapeMarkdownV2(FormatCockSizeForDate(neighbor.Size)),
				"SizeEmoji":  EmojiFromSize(neighbor.Size),
			}))
		}

		if neighborhood.Self != nil {
			self := neighborhood.Self
			contextLines = append(contextLines, app.localization.Localize(localizer, MsgCockRulerContextSelected, map[string]any{
				"PlaceEmoji": fmt.Sprintf("🥀 *%d*\\.", *data.UserPosition),
				"Username":   EscapeMarkdownV2(self.Nickname),
				"Size":       EscapeMarkdownV2(FormatCockSizeForDate(self.Size)),
				"SizeEmoji":  EmojiFromSize(self.Size),
			}))
		}

		for _, neighbor := range neighborhood.Below {
			contextLines = append(contextLines, app.localization.Localize(localizer, MsgCockRulerContextDefault, map[string]any{
				"PlaceEmoji": "🥀",
				"Username":   EscapeMarkdownV2(neighbor.Nickname),
				"Size":       EscapeMarkdownV2(FormatCockSizeForDate(neighbor.Size)),
				"SizeEmoji":  EmojiFromSize(neighbor.Size),
			}))
		}

		if len(contextLines) > 0 {
			dots := app.localization.Localize(localizer, CommonDots, nil)
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
		others = append(others, app.localization.Localize(localizer, MsgCockScoreboardNotFound, nil))
	}

	if len(others) != 0 {
		template := MsgCockRulerScoreboardTemplate
		if !showDescription {
			template = MsgCockRulerScoreboardTemplateNoDesc
		}
		return app.localization.Localize(localizer, template, map[string]any{
			"Participants": data.TotalParticipants,
			"Winners":      strings.Join(winners, "\n"),
			"Others":       strings.Join(others, "\n"),
		})
	}
	template := MsgCockRulerScoreboardWinnersTemplate
	if !showDescription {
		template = MsgCockRulerScoreboardWinnersTemplateNoDesc
	}
	return app.localization.Localize(localizer, template, map[string]any{
		"Participants": data.TotalParticipants,
		"Winners":      strings.Join(winners, "\n"),
	})
}

func (app *Application) GenerateCockRaceScoreboard(log *logging.Logger, localizer *i18n.Localizer, userID int64, data *api.CockRaceData, showDescription bool) string {
	var winners []string
	var others []string
	isUserInScoreboard := false

	for index, entry := range data.Leaders {
		isCurrentUser := entry.UserID == userID
		emoji := GetPlaceEmoji(index+1, isCurrentUser)

		if isCurrentUser {
			isUserInScoreboard = true
		}

		var scoreboardLine string
		if isCurrentUser {
			scoreboardLine = app.localization.Localize(localizer, MsgCockRaceScoreboardSelected, map[string]any{
				"PlaceEmoji": emoji,
				"Username":   EscapeMarkdownV2(entry.Nickname),
				"Size":       EscapeMarkdownV2(FormatDickSize(entry.TotalSize)),
			})
		} else {
			scoreboardLine = app.localization.Localize(localizer, MsgCockRaceScoreboardDefault, map[string]any{
				"PlaceEmoji": emoji,
				"Username":   EscapeMarkdownV2(entry.Nickname),
				"Size":       EscapeMarkdownV2(FormatDickSize(entry.TotalSize)),
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
			contextLines = append(contextLines, app.localization.Localize(localizer, MsgCockRaceContextDefault, map[string]any{
				"PlaceEmoji": "🥀",
				"Username":   EscapeMarkdownV2(neighbor.Nickname),
				"Size":       EscapeMarkdownV2(FormatDickSize(neighbor.TotalSize)),
			}))
		}

		if neighborhood.Self != nil {
			self := neighborhood.Self
			contextLines = append(contextLines, app.localization.Localize(localizer, MsgCockRaceContextSelected, map[string]any{
				"PlaceEmoji": fmt.Sprintf("🥀 *%d*\\.", *data.UserPosition),
				"Username":   EscapeMarkdownV2(self.Nickname),
				"Size":       EscapeMarkdownV2(FormatDickSize(self.TotalSize)),
			}))
		}

		for _, neighbor := range neighborhood.Below {
			contextLines = append(contextLines, app.localization.Localize(localizer, MsgCockRaceContextDefault, map[string]any{
				"PlaceEmoji": "🥀",
				"Username":   EscapeMarkdownV2(neighbor.Nickname),
				"Size":       EscapeMarkdownV2(FormatDickSize(neighbor.TotalSize)),
			}))
		}

		if len(contextLines) > 0 {
			dots := app.localization.Localize(localizer, CommonDots, nil)
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
		others = append(others, app.localization.Localize(localizer, MsgCockScoreboardNotFound, nil))
	}

	// Footer
	seasonStart := EscapeMarkdownV2(data.Season.StartDate.FormatDateMSK())
	seasonEnd := EscapeMarkdownV2(data.Season.EndDate.FormatDateMSK())
	seasonNum := data.Season.SeasonNum
	seasonWord := app.localization.Localize(localizer, UnitSeasonGenitive, map[string]any{"Count": seasonNum})

	now := datetime.NowTime()
	timeRemaining := FormatTimeRemaining(app.localization, localizer, data.Season.EndDate.Time, now)

	footerLine := app.localization.Localize(localizer, MsgCockRaceFooterActiveSeason, map[string]any{
		"SeasonNum": seasonNum,
		"StartDate": seasonStart,
		"EndDate":   seasonEnd,
		"Remaining": EscapeMarkdownV2(timeRemaining),
	})

	if len(others) != 0 {
		template := MsgCockRaceScoreboardTemplate
		if !showDescription {
			template = MsgCockRaceScoreboardTemplateNoDesc
		}
		return app.localization.Localize(localizer, template, map[string]any{
			"Participants": data.TotalParticipants,
			"Winners":      strings.Join(winners, "\n"),
			"Others":       strings.Join(others, "\n"),
			"Footer":       footerLine,
			"SeasonNum":    seasonNum,
			"SeasonWord":   seasonWord,
		})
	}
	template := MsgCockRaceScoreboardWinnersTemplate
	if !showDescription {
		template = MsgCockRaceScoreboardWinnersTemplateNoDesc
	}
	return app.localization.Localize(localizer, template, map[string]any{
		"Participants": data.TotalParticipants,
		"Winners":      strings.Join(winners, "\n"),
		"Footer":       footerLine,
		"SeasonNum":    seasonNum,
		"SeasonWord":   seasonWord,
	})
}

func (app *Application) GenerateCockLadderScoreboard(log *logging.Logger, localizer *i18n.Localizer, userID int64, data *api.CockLadderData, showDescription bool) string {
	var winners []string
	var others []string
	isUserInScoreboard := false

	for index, entry := range data.Leaders {
		isCurrentUser := entry.UserID == userID
		emoji := GetPlaceEmoji(index+1, isCurrentUser)

		if isCurrentUser {
			isUserInScoreboard = true
		}

		var scoreboardLine string
		if isCurrentUser {
			scoreboardLine = app.localization.Localize(localizer, MsgCockLadderScoreboardSelected, map[string]any{
				"PlaceEmoji": emoji,
				"Username":   EscapeMarkdownV2(entry.Nickname),
				"Size":       EscapeMarkdownV2(FormatDickSize(entry.TotalSize)),
			})
		} else {
			scoreboardLine = app.localization.Localize(localizer, MsgCockLadderScoreboardDefault, map[string]any{
				"PlaceEmoji": emoji,
				"Username":   EscapeMarkdownV2(entry.Nickname),
				"Size":       EscapeMarkdownV2(FormatDickSize(entry.TotalSize)),
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
			contextLines = append(contextLines, app.localization.Localize(localizer, MsgCockLadderContextDefault, map[string]any{
				"PlaceEmoji": "🥀",
				"Username":   EscapeMarkdownV2(neighbor.Nickname),
				"Size":       EscapeMarkdownV2(FormatDickSize(neighbor.TotalSize)),
			}))
		}

		if neighborhood.Self != nil {
			self := neighborhood.Self
			contextLines = append(contextLines, app.localization.Localize(localizer, MsgCockLadderContextSelected, map[string]any{
				"PlaceEmoji": fmt.Sprintf("🥀 *%d*\\.", *data.UserPosition),
				"Username":   EscapeMarkdownV2(self.Nickname),
				"Size":       EscapeMarkdownV2(FormatDickSize(self.TotalSize)),
			}))
		}

		for _, neighbor := range neighborhood.Below {
			contextLines = append(contextLines, app.localization.Localize(localizer, MsgCockLadderContextDefault, map[string]any{
				"PlaceEmoji": "🥀",
				"Username":   EscapeMarkdownV2(neighbor.Nickname),
				"Size":       EscapeMarkdownV2(FormatDickSize(neighbor.TotalSize)),
			}))
		}

		if len(contextLines) > 0 {
			dots := app.localization.Localize(localizer, CommonDots, nil)
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
		others = append(others, app.localization.Localize(localizer, MsgCockScoreboardNotFound, nil))
	}

	if len(others) != 0 {
		template := MsgCockLadderScoreboardTemplate
		if !showDescription {
			template = MsgCockLadderScoreboardTemplateNoDesc
		}
		return app.localization.Localize(localizer, template, map[string]any{
			"Participants": data.TotalParticipants,
			"Winners":      strings.Join(winners, "\n"),
			"Others":       strings.Join(others, "\n"),
		})
	}
	template := MsgCockLadderScoreboardWinnersTemplate
	if !showDescription {
		template = MsgCockLadderScoreboardWinnersTemplateNoDesc
	}
	return app.localization.Localize(localizer, template, map[string]any{
		"Participants": data.TotalParticipants,
		"Winners":      strings.Join(winners, "\n"),
	})
}

func GetPlaceEmoji(place int, isCurrentUser bool) string {
	switch place {
	case 1:
		return "🥇"
	case 2:
		return "🥈"
	case 3:
		return "🥉"
	default:
		now := datetime.NowTime()
		month := now.Month()

		var emoji string
		switch month {
		case time.March, time.April, time.May:
			emoji = "🫠"
		case time.June, time.July, time.August:
			emoji = "🥵"
		case time.September, time.October, time.November:
			emoji = "🤧"
		default:
			emoji = "🥶"
		}

		if isCurrentUser {
			return fmt.Sprintf("%s *%d*\\.", emoji, place)
		}
		return fmt.Sprintf("%s %d\\.", emoji, place)
	}
}

func GetPlaceEmojiForContext(place int, bold bool) string {
	if bold {
		return fmt.Sprintf("🥀 *%d*\\.", place)
	}
	return fmt.Sprintf("🥀 %d\\.", place)
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
	case luck >= 1.98:
		return "👑🌌🌈🦄🍀🤩"
	case luck >= 1.92:
		return "🌌🌈🦄🍀🤩"
	case luck >= 1.833:
		return "🌈🦄🍀🤩"
	case luck >= 1.7:
		return "🍀🤩"
	case luck >= 1.5:
		return "🤩"
	case luck >= 1.2:
		return "🍀✨"
	case luck >= 1.1:
		return "🍀"
	case luck >= 0.9:
		return "⚖️"
	case luck >= 0.7:
		return "😕"
	case luck >= 0.5:
		return "😔"
	case luck >= 0.3:
		return "💀"
	case luck >= 0.2:
		return "☠️"
	default:
		return "🔥☠️🔥"
	}
}

func LuckLabel(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, luck float64) string {
	switch {
	case luck >= 1.98:
		return localizationManager.Localize(localizer, LuckLabelGodRandom, nil)
	case luck >= 1.92:
		return localizationManager.Localize(localizer, LuckLabelCosmicLuck, nil)
	case luck >= 1.833:
		return localizationManager.Localize(localizer, LuckLabelFairyLuck, nil)
	case luck >= 1.7:
		return localizationManager.Localize(localizer, LuckLabelSuperLuck, nil)
	case luck >= 1.5:
		return localizationManager.Localize(localizer, LuckLabelIncredibleLuck, nil)
	case luck >= 1.2:
		return localizationManager.Localize(localizer, LuckLabelVeryLucky, nil)
	case luck >= 1.1:
		return localizationManager.Localize(localizer, LuckLabelLucky, nil)
	case luck >= 0.9:
		return localizationManager.Localize(localizer, LuckLabelBalanced, nil)
	case luck >= 0.7:
		return localizationManager.Localize(localizer, LuckLabelUnlucky, nil)
	case luck >= 0.5:
		return localizationManager.Localize(localizer, LuckLabelBad, nil)
	case luck >= 0.3:
		return localizationManager.Localize(localizer, LuckLabelGloom, nil)
	case luck >= 0.2:
		return localizationManager.Localize(localizer, LuckLabelHellTilt, nil)
	default:
		return localizationManager.Localize(localizer, LuckLabelBurningInHell, nil)
	}
}

func LuckDisplay(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, luck float64) string {
	return fmt.Sprintf("%s _(%s)_", LuckEmoji(luck), LuckLabel(localizationManager, localizer, luck))
}

func VolatilityEmoji(volatility float64) string {
	switch {
	case volatility < 1:
		return "🧱"
	case volatility < 3:
		return "🧊"
	case volatility < 6:
		return "📈"
	case volatility < 10:
		return "📉📈"
	case volatility < 15:
		return "🎢"
	case volatility < 25:
		return "🎢🌪️"
	default:
		return "🌪️💥"
	}
}

func VolatilityLabel(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, volatility float64) string {
	switch {
	case volatility < 1:
		return localizationManager.Localize(localizer, VolatilityLabelStone, nil)
	case volatility < 3:
		return localizationManager.Localize(localizer, VolatilityLabelStable, nil)
	case volatility < 6:
		return localizationManager.Localize(localizer, VolatilityLabelModerate, nil)
	case volatility < 10:
		return localizationManager.Localize(localizer, VolatilityLabelLivelySpread, nil)
	case volatility < 15:
		return localizationManager.Localize(localizer, VolatilityLabelUneven, nil)
	case volatility < 25:
		return localizationManager.Localize(localizer, VolatilityLabelChaotic, nil)
	default:
		return localizationManager.Localize(localizer, VolatilityLabelRandom, nil)
	}
}

func VolatilityDisplay(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, volatility float64) string {
	return fmt.Sprintf("%s _(%s)_", VolatilityEmoji(volatility), VolatilityLabel(localizationManager, localizer, volatility))
}

func clamp01(x float64) float64 {
	if math.IsNaN(x) {
		return 0
	}
	if x < 0 {
		return 0
	}
	if x > 1 {
		return 1
	}
	return x
}

func IrkLabel(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, irk float64) string {
	irk = clamp01(irk)

	bucket := int(math.Floor(irk * 10))
	if irk >= 1.0 {
		bucket = 10
	}

	labels := [...]string{
		IrkLabelZero, IrkLabelMinimal, IrkLabelVerySmall, IrkLabelSmall, IrkLabelReduced,
		IrkLabelAverage, IrkLabelIncreased, IrkLabelLarge, IrkLabelVeryLarge, IrkLabelMaximum, IrkLabelUltimate,
	}

	return localizationManager.Localize(localizer, labels[bucket], nil)
}

func GrowthSpeedLabel(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, speed float64) string {
	switch {
	case speed >= 50:
		return localizationManager.Localize(localizer, GrowthSpeedLabelCosmic, nil)
	case speed >= 40:
		return localizationManager.Localize(localizer, GrowthSpeedLabelExtreme, nil)
	case speed >= 30:
		return localizationManager.Localize(localizer, GrowthSpeedLabelVeryFast, nil)
	case speed >= 20:
		return localizationManager.Localize(localizer, GrowthSpeedLabelFast, nil)
	case speed >= 15:
		return localizationManager.Localize(localizer, GrowthSpeedLabelModerate, nil)
	case speed >= 10:
		return localizationManager.Localize(localizer, GrowthSpeedLabelAverage, nil)
	case speed >= 5:
		return localizationManager.Localize(localizer, GrowthSpeedLabelSlow, nil)
	case speed >= 2:
		return localizationManager.Localize(localizer, GrowthSpeedLabelVerySlow, nil)
	case speed >= 0.5:
		return localizationManager.Localize(localizer, GrowthSpeedLabelTurtle, nil)
	default:
		return localizationManager.Localize(localizer, GrowthSpeedLabelStalled, nil)
	}
}

func GrowthSpeedEmoji(speed float64) string {
	switch {
	case speed >= 50:
		return "👑🌌🚀💫"
	case speed >= 40:
		return "🚀🔥⚡"
	case speed >= 30:
		return "⚡💨🏎️"
	case speed >= 20:
		return "🏃💨"
	case speed >= 15:
		return "🚶‍♂️⏱️"
	case speed >= 10:
		return "🚶"
	case speed >= 5:
		return "🐢⏳"
	case speed >= 2:
		return "🐌🕰️"
	case speed >= 0.5:
		return "🐢🌿"
	default:
		return "🗿⛔"
	}
}

func GrowthSpeedDisplay(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, speed float64) string {
	emoji := GrowthSpeedEmoji(speed)
	label := GrowthSpeedLabel(localizationManager, localizer, speed)
	return fmt.Sprintf("%s _(%s)_", emoji, label)
}

func FormatGrowthSpeed(speed float64) string {
	p := message.NewPrinter(language.Russian)
	return p.Sprintf("%.1f", speed)
}

func FormatTimeRemaining(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, endDate time.Time, now time.Time) string {
	duration := endDate.Sub(now)
	daysRemaining := int(duration.Hours() / 24)

	if daysRemaining < 0 {
		return localizationManager.Localize(localizer, UnitDay, map[string]any{"Count": 0})
	}

	if daysRemaining > 30 {
		months := daysRemaining / 30
		days := daysRemaining % 30

		if days == 0 {
			return localizationManager.Localize(localizer, UnitMonth, map[string]any{"Count": months})
		}
		monthsText := localizationManager.Localize(localizer, UnitMonth, map[string]any{"Count": months})
		daysText := localizationManager.Localize(localizer, UnitDay, map[string]any{"Count": days})
		return fmt.Sprintf("%s %s", monthsText, daysText)
	}

	return localizationManager.Localize(localizer, UnitDay, map[string]any{"Count": daysRemaining})
}

func FormatUserPullingPeriod(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, firstCockDate time.Time, now time.Time) string {
	years := now.Year() - firstCockDate.Year()
	months := int(now.Month()) - int(firstCockDate.Month())
	days := now.Day() - firstCockDate.Day()

	if days < 0 {
		months--
		prevMonth := now.AddDate(0, -1, 0)
		daysInPrevMonth := time.Date(prevMonth.Year(), prevMonth.Month()+1, 0, 0, 0, 0, 0, prevMonth.Location()).Day()
		days += daysInPrevMonth
	}

	if months < 0 {
		years--
		months += 12
	}

	dateStr := firstCockDate.Format("02.01.2006")

	var parts []string
	if years > 0 {
		parts = append(parts, localizationManager.Localize(localizer, UnitYear, map[string]any{"Count": years}))
	}
	if months > 0 {
		parts = append(parts, localizationManager.Localize(localizer, UnitMonth, map[string]any{"Count": months}))
	}
	if days > 0 || len(parts) == 0 {
		parts = append(parts, localizationManager.Localize(localizer, UnitDay, map[string]any{"Count": days}))
	}

	var result string
	if len(parts) == 1 {
		result = parts[0]
	} else if len(parts) == 2 {
		result = parts[0] + localizationManager.Localize(localizer, MsgListSeparatorLast, nil) + parts[1]
	} else if len(parts) == 3 {
		result = parts[0] + localizationManager.Localize(localizer, MsgListSeparator, nil) + parts[1] + localizationManager.Localize(localizer, MsgListSeparatorLast, nil) + parts[2]
	}

	return localizationManager.Localize(localizer, MsgUserPullingSince, map[string]any{
		"Period": result,
		"Date":   dateStr,
	})
}

// GenerateAchievementsText генерирует текст списка достижений с пагинацией
func GenerateAchievementsText(
	localizationManager *localization.LocalizationManager,
	localizer *i18n.Localizer,
	allAchievements []AchievementDef,
	apiAchievements []api.AchievementData,
	page int,
	itemsPerPage int,
) string {
	// Создаём map из API-данных для быстрого доступа
	apiMap := make(map[string]*api.AchievementData, len(apiAchievements))
	for i := range apiAchievements {
		apiMap[apiAchievements[i].ID] = &apiAchievements[i]
	}

	type AchievementWithStatus struct {
		Def         AchievementDef
		APIData     *api.AchievementData
		IsCompleted bool
	}

	achievementsWithStatus := make([]AchievementWithStatus, 0, len(allAchievements))

	for _, def := range allAchievements {
		apiData := apiMap[def.ID]
		isCompleted := apiData != nil && apiData.Completed

		achievementsWithStatus = append(achievementsWithStatus, AchievementWithStatus{
			Def:         def,
			APIData:     apiData,
			IsCompleted: isCompleted,
		})
	}

	sort.Slice(achievementsWithStatus, func(i, j int) bool {
		if achievementsWithStatus[i].IsCompleted != achievementsWithStatus[j].IsCompleted {
			return achievementsWithStatus[i].IsCompleted
		}
		return false
	})

	totalPages := (len(achievementsWithStatus) + itemsPerPage - 1) / itemsPerPage
	if page < 1 {
		page = 1
	}
	if page > totalPages {
		page = totalPages
	}

	startIdx := (page - 1) * itemsPerPage
	endIdx := startIdx + itemsPerPage
	if endIdx > len(achievementsWithStatus) {
		endIdx = len(achievementsWithStatus)
	}

	var lines []string
	for i := startIdx; i < endIdx; i++ {
		achStatus := achievementsWithStatus[i]
		line := FormatAchievementLine(localizationManager, localizer, achStatus.Def, achStatus.APIData, achStatus.IsCompleted)
		lines = append(lines, line)
	}

	return strings.Join(lines, "\n")
}

// FormatAchievementLine форматирует одну строку достижения
func FormatAchievementLine(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, def AchievementDef, apiData *api.AchievementData, isCompleted bool) string {
	escapedName := EscapeMarkdownV2(localizationManager.Localize(localizer, def.Name, nil))
	escapedDesc := EscapeMarkdownV2(localizationManager.Localize(localizer, def.Description, nil))

	emoji := "🔒"
	if apiData != nil {
		emoji = apiData.Emoji
	}

	if isCompleted {
		return localizationManager.Localize(localizer, MsgAchievementCompleted, map[string]any{
			"Emoji":       emoji,
			"Name":        escapedName,
			"Description": escapedDesc,
		})
	} else if apiData != nil && apiData.Progress > 0 && apiData.MaxProgress > 0 {
		return localizationManager.Localize(localizer, MsgAchievementInProgress, map[string]any{
			"Emoji":       emoji,
			"Name":        escapedName,
			"Progress":    apiData.Progress,
			"Max":         apiData.MaxProgress,
			"Description": escapedDesc,
		})
	}
	return localizationManager.Localize(localizer, MsgAchievementNotCompleted, map[string]any{
		"Emoji":       emoji,
		"Name":        escapedName,
		"Description": escapedDesc,
	})
}
