package application

import (
	"dickobrazz/application/datetime"
	"dickobrazz/application/logging"
	"fmt"
	"strconv"
	"strings"
	"time"

	"math/rand"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var glitchMarks = []rune{
	'\u0335', '\u0336', '\u0337', '\u0338', // зачёркивание
	'\u0300', '\u0301', '\u0302', '\u0303', // диакритика
	'\u0304', '\u0305', '\u0306', '\u0307', // черточки
	'\u0308', '\u0309', '\u030A', '\u030B',
	'\u0310', '\u0311', '\u0312', '\u0313',
	'\u0334', '\u034F', '\u0350', '\u0351',
	'\u0352', '\u0353', '\u0354', '\u0355', '\u0356',
}

var mathFancy = map[int]string{
	0:  "sin(0)",
	1:  "0!",
	2:  "C(2,1)",
	3:  "1! + 2!",
	4:  "2²",
	5:  "√25",
	6:  "3!",
	7:  "3! + 1",
	8:  "2³",
	9:  "3²",
	10: "C(5,2)",
	11: "(1011)₂",
	12: "4! / 2",
	13: "F₇",
	14: "Cat₄",
	15: "C(6,2)",
	16: "2⁴",
	17: "√289",
	18: "3! · 3",
	19: "3³ − 2³",
	20: "5! / 6",
	21: "F₈",
	22: "⌊π^e⌋",
	23: "⌈π^e⌉",
	24: "4!",
	25: "5²",
	26: "4! + 2!",
	27: "3³",
	28: "T₇ = 7·8/2",
	29: "2⁵ − 3",
	30: "2 · 5!!",
	31: "2⁵ − 1",
	32: "2⁵",
	33: "4! + 3! + 2! + 0!",
	34: "F₉",
	35: "C(7,3)",
	36: "6²",
	37: "⌊12π⌋",
	38: "(100110)₂",
	39: "3³ + 2·3!",
	40: "5! / 3",
	41: "n² + n + 41 |_{n=0}",
	42: "Cat₅",
	43: "⌊14π⌋",
	44: "⌊√2000⌋",
	45: "C(10,2)",
	46: "4! + 4! − 2!",
	47: "⌊15π⌋",
	48: "4! · 2",
	49: "7²",
	50: "⌊16π⌋",
	51: "4! + 3³",
	52: "6!! + 2²",
	53: "⌊17π⌋",
	54: "3³ + 3³",
	55: "F₁₀",
	56: "C(8,3)",
	57: "4! + 3! + 3³",
	58: "6!! + C(5,2)",
	59: "⌊19π⌋",
	60: "5! / 2",
	61: "√3721",
}

var rnd = rand.New(rand.NewSource(time.Now().UnixNano()))

// isMathDay — 14 марта (International Day of Mathematics / Pi Day)
func isMathDay(t time.Time) bool {
	return t.Month() == time.March && t.Day() == 14
}

// isProgrammersDay — 256-й день года (12/13 сентября)
func isProgrammersDay(t time.Time) bool {
	return t.YearDay() == 256
}

func toProgrammersNotation(n int) string {
	if rnd.Intn(2) == 0 {
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
		// добавляем 1–3 случайных глитч символа
		count := rnd.Intn(3) + 1
		for i := 0; i < count; i++ {
			sb.WriteRune(glitchMarks[rnd.Intn(len(glitchMarks))])
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
	
	// 1 апреля - День смеха: отрицательный размер
	if now.Month() == time.April && now.Day() == 1 {
		displaySize = -size
	}

	// 14 марта - День математика: математические выражения
	if isMathDay(now) {
		return fancyMathOrDefault(displaySize)
	}

	// 256-й день года - День программиста: двоичная/шестнадцатеричная нотация
	if isProgrammersDay(now) {
		return toProgrammersNotation(displaySize)
	}

	// 31 октября - Хэллоуин: глитченный текст
	if now.Month() == time.October && now.Day() == 31 {
		return glitchify(strconv.Itoa(displaySize))
	}

	return strconv.Itoa(displaySize)
}

func GenerateCockSizeText(size int, emoji string) string {
	formattedSize := FormatCockSizeForDate(size)
	return fmt.Sprintf(MsgCockSize, formattedSize, emoji)
}

func (app *Application) GenerateCockRulerText(log *logging.Logger, userID int64, cocks []UserCock) string {
	var winners []string
	var others []string
	isUserInScoreboard := false

	for index, cock := range cocks {
		isCurrentUser := cock.UserId == userID
		emoji := GetPlaceEmoji(index + 1)
		formattedSize := FormatCockSizeForDate(cock.Size)

		var line string
		if isCurrentUser {
			isUserInScoreboard = true
			line = fmt.Sprintf(MsgCockRulerScoreboardSelected, emoji, EscapeMarkdownV2(cock.UserName), formattedSize, EmojiFromSize(cock.Size))
		} else {
			line = fmt.Sprintf(MsgCockRulerScoreboardDefault, emoji, EscapeMarkdownV2(cock.UserName), formattedSize, EmojiFromSize(cock.Size))
		}

		if index < 3 {
			winners = append(winners, line)
		} else {
			others = append(others, line)
		}
	}

	if !isUserInScoreboard {
		if userCock := app.GetCockSizeFromCache(log, userID); userCock != nil {
			formattedUserCock := FormatCockSizeForDate(*userCock)
			others = append(others, fmt.Sprintf(MsgCockRulerScoreboardOut, EscapeMarkdownV2(CommonDots), formattedUserCock, EmojiFromSize(*userCock)))
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

func (app *Application) GenerateCockLadderScoreboard(log *logging.Logger, userID int64, sizes []UserCockRace) string {
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
			scoreboardLine = fmt.Sprintf(MsgCockLadderScoreboardSelected, emoji, EscapeMarkdownV2(user.Nickname), FormatDickSize(int(user.TotalSize)))
		} else {
			scoreboardLine = fmt.Sprintf(MsgCockLadderScoreboardDefault, emoji, EscapeMarkdownV2(user.Nickname), FormatDickSize(int(user.TotalSize)))
		}

		if index < 3 {
			winners = append(winners, scoreboardLine)
		} else {
			others = append(others, scoreboardLine)
		}
	}

	if !isUserInScoreboard {
		if cock := app.GetUserAggregatedCock(log, userID); cock != nil {
			others = append(others, fmt.Sprintf(MsgCockLadderScoreboardOut, EscapeMarkdownV2(cock.Nickname), FormatDickSize(int(cock.TotalSize))))
		} else {
			others = append(others, MsgCockScoreboardNotFound)
		}
	}

	if len(others) != 0 {
		return fmt.Sprintf(
			MsgCockLadderScoreboardTemplate,
			strings.Join(winners, "\n"),
			strings.Join(others, "\n"),
		)
	} else {
		return fmt.Sprintf(
			MsgCockLadderScoreboardWinnersTemplate,
			strings.Join(winners, "\n"),
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
		now := time.Now()
		month := now.Month()

		switch month {
		case time.March, time.April, time.May:
			return "🫠"
		case time.June, time.July, time.August:
			return "🥵"
		case time.September, time.October, time.November:
			return "🤧"
		default:
			return "🥶"
		}
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
  case luck >= 1.98: // типа бог рандома :)
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
	case luck >= 0.2: // адовый тильт
		return "☠️"
	default:
		return "🔥☠️🔥"
	}
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

func VolatilityLabel(volatility float64) string {
	switch {
	case volatility < 1:
		return "каменный"
	case volatility < 3:
		return "стабильный"
	case volatility < 6:
		return "умеренный"
	case volatility < 10:
		return "живой разброс"
	case volatility < 15:
		return "неровный"
	case volatility < 25:
		return "хаотичный"
	default:
		return "полный рандом"
	}
}

func VolatilityDisplay(volatility float64) string {
	return fmt.Sprintf("%s _(%s)_", VolatilityEmoji(volatility), VolatilityLabel(volatility))
}