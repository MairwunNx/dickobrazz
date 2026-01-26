package application

import (
	"dickobrazz/application/database"
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

// GenerateAnonymousName генерирует анонимное имя для пользователя без username
// Использует PRNG с seed из userID для генерации стабильного номера (0-9999)
func GenerateAnonymousName(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, userID int64) string {
	// Создаем отдельный генератор с seed из userID для стабильности
	rng := rand.New(rand.NewSource(userID))
	number := rng.Intn(10000)
	numberStr := fmt.Sprintf("%04d", number)
	return localizationManager.Localize(localizer, AnonymousNameTemplate, map[string]any{"Number": numberStr})
}

// NormalizeUsername возвращает username пользователя или генерирует анонимное имя
func NormalizeUsername(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, username string, userID int64) string {
	if username == "" {
		return GenerateAnonymousName(localizationManager, localizer, userID)
	}
	return username
}

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

var (
	rnd   = rand.New(rand.NewSource(time.Now().UnixNano()))
	rndMu sync.Mutex
)

// isMathDay — 14 марта (International Day of Mathematics / Pi Day)
func isMathDay(t time.Time) bool {
	return t.Month() == time.March && t.Day() == 14
}

// isProgrammersDay — 256-й день года (12/13 сентября)
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
		// добавляем 1–3 случайных глитч символа
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

func GenerateCockSizeText(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, size int, emoji string) string {
	formattedSize := FormatCockSizeForDate(size)
	return localizationManager.Localize(localizer, MsgCockSize, map[string]any{
		"Size":  formattedSize,
		"Emoji": emoji,
	})
}

func (app *Application) GenerateCockRulerText(log *logging.Logger, localizer *i18n.Localizer, userID int64, cocks []UserCock, totalParticipants int, showDescription bool) string {
	var winners []string
	var others []string
	isUserInScoreboard := false

	for index, cock := range cocks {
		isCurrentUser := cock.UserId == userID
		emoji := GetPlaceEmoji(index+1, isCurrentUser)
		formattedSize := FormatCockSizeForDate(cock.Size)

		// Нормализуем username с учетом скрытия
		normalizedUsername := app.ResolveDisplayNickname(log, localizer, cock.UserId, cock.UserName)

		var line string
		if isCurrentUser {
			isUserInScoreboard = true
			line = app.localization.Localize(localizer, MsgCockRulerScoreboardSelected, map[string]any{
				"PlaceEmoji": emoji,
				"Username":   EscapeMarkdownV2(normalizedUsername),
				"Size":       formattedSize,
				"SizeEmoji":  EmojiFromSize(cock.Size),
			})
		} else {
			line = app.localization.Localize(localizer, MsgCockRulerScoreboardDefault, map[string]any{
				"PlaceEmoji": emoji,
				"Username":   EscapeMarkdownV2(normalizedUsername),
				"Size":       formattedSize,
				"SizeEmoji":  EmojiFromSize(cock.Size),
			})
		}

		if index < 3 {
			winners = append(winners, line)
		} else {
			others = append(others, line)
		}
	}

	if !isUserInScoreboard {
		if userCock := app.GetCockSizeFromCache(log, userID); userCock != nil {
			// Получаем все коки из кеша для определения позиции и соседей
			allCocks := app.GetCockSizesFromCache(log)
			sort.Slice(allCocks, func(i, j int) bool {
				return allCocks[i].Size > allCocks[j].Size
			})

			// Находим позицию пользователя
			userPosition := 0
			for idx, cock := range allCocks {
				if cock.UserId == userID {
					userPosition = idx + 1
					break
				}
			}

			// Определяем диапазон для показа (обрабатываем edge cases)
			var startIdx, endIdx int
			totalCount := len(allCocks)

			if userPosition == 14 {
				// Сразу после топ-13 - показываем только текущего и следующего
				startIdx = userPosition - 1 // индекс 13 (14-е место)
				endIdx = startIdx + 2
				if endIdx > totalCount {
					endIdx = totalCount
				}
			} else if userPosition >= totalCount-1 {
				// Последние 2 места - показываем предыдущего и текущего
				startIdx = userPosition - 2
				if startIdx < 13 {
					startIdx = 13 // не залезаем в топ-13
				}
				endIdx = totalCount
			} else {
				// Обычный случай - показываем предыдущего, текущего, следующего
				startIdx = userPosition - 2
				if startIdx < 13 {
					startIdx = 13 // не залезаем в топ-13
				}
				endIdx = startIdx + 3
				if endIdx > totalCount {
					endIdx = totalCount
				}
			}

			neighbors := allCocks[startIdx:endIdx]

			// Формируем строки для соседей
			var contextLines []string
			showTopDots := startIdx > 13          // Показываем точки сверху если есть пропуск после топ-13
			showBottomDots := endIdx < totalCount // Показываем точки снизу если есть что-то дальше

			for idx, neighbor := range neighbors {
				pos := startIdx + idx + 1
				isCurrentInContext := neighbor.UserId == userID
				normalizedNick := app.ResolveDisplayNickname(log, localizer, neighbor.UserId, neighbor.UserName)
				formattedSize := FormatCockSizeForDate(neighbor.Size)
				emoji := EmojiFromSize(neighbor.Size)
				posEmoji := GetPlaceEmojiForContext(pos, isCurrentInContext)

				if isCurrentInContext {
					contextLines = append(contextLines, app.localization.Localize(localizer, MsgCockRulerContextSelected, map[string]any{
						"PlaceEmoji": posEmoji,
						"Username":   EscapeMarkdownV2(normalizedNick),
						"Size":       EscapeMarkdownV2(formattedSize),
						"SizeEmoji":  emoji,
					}))
				} else {
					contextLines = append(contextLines, app.localization.Localize(localizer, MsgCockRulerContextDefault, map[string]any{
						"PlaceEmoji": posEmoji,
						"Username":   EscapeMarkdownV2(normalizedNick),
						"Size":       EscapeMarkdownV2(formattedSize),
						"SizeEmoji":  emoji,
					}))
				}
			}

			// Добавляем контекст с соседями
			var contextBlock string
			if showTopDots && showBottomDots {
				dots := app.localization.Localize(localizer, CommonDots, nil)
				contextBlock = "\n" + dots + "\n" + strings.Join(contextLines, "\n") + "\n" + dots
			} else if showTopDots {
				dots := app.localization.Localize(localizer, CommonDots, nil)
				contextBlock = "\n" + dots + "\n" + strings.Join(contextLines, "\n")
			} else if showBottomDots {
				dots := app.localization.Localize(localizer, CommonDots, nil)
				contextBlock = "\n" + strings.Join(contextLines, "\n") + "\n" + dots
			} else {
				contextBlock = "\n" + strings.Join(contextLines, "\n")
			}

			others = append(others, contextBlock)
		} else {
			others = append(others, app.localization.Localize(localizer, MsgCockScoreboardNotFound, nil))
		}
	}

	if len(others) != 0 {
		template := MsgCockRulerScoreboardTemplate
		if !showDescription {
			template = MsgCockRulerScoreboardTemplateNoDesc
		}
		return app.localization.Localize(localizer, template, map[string]any{
			"Participants": totalParticipants,
			"Winners":      strings.Join(winners, "\n"),
			"Others":       strings.Join(others, "\n"),
		})
	} else {
		template := MsgCockRulerScoreboardWinnersTemplate
		if !showDescription {
			template = MsgCockRulerScoreboardWinnersTemplateNoDesc
		}
		return app.localization.Localize(localizer, template, map[string]any{
			"Participants": totalParticipants,
			"Winners":      strings.Join(winners, "\n"),
		})
	}
}

func (app *Application) GenerateCockRaceScoreboard(log *logging.Logger, localizer *i18n.Localizer, userID int64, sizes []UserCockRace, seasonStart string, totalParticipants int, currentSeason *CockSeason, showDescription bool) string {
	var winners []string
	var others []string
	isUserInScoreboard := false

	for index, user := range sizes {
		isCurrentUser := user.UserID == userID
		emoji := GetPlaceEmoji(index+1, isCurrentUser)

		if isCurrentUser {
			isUserInScoreboard = true
		}

		// Нормализуем username с учетом скрытия
		normalizedNickname := app.ResolveDisplayNickname(log, localizer, user.UserID, user.Nickname)

		var scoreboardLine string
		if isCurrentUser {
			scoreboardLine = app.localization.Localize(localizer, MsgCockRaceScoreboardSelected, map[string]any{
				"PlaceEmoji": emoji,
				"Username":   EscapeMarkdownV2(normalizedNickname),
				"Size":       EscapeMarkdownV2(FormatDickSize(int(user.TotalSize))),
			})
		} else {
			scoreboardLine = app.localization.Localize(localizer, MsgCockRaceScoreboardDefault, map[string]any{
				"PlaceEmoji": emoji,
				"Username":   EscapeMarkdownV2(normalizedNickname),
				"Size":       EscapeMarkdownV2(FormatDickSize(int(user.TotalSize))),
			})
		}

		if index < 3 {
			winners = append(winners, scoreboardLine)
		} else {
			others = append(others, scoreboardLine)
		}
	}

	if !isUserInScoreboard {
		if cock := app.GetUserAggregatedCock(log, userID); cock != nil {
			// Получаем позицию пользователя
			var userPosition int
			var neighbors []UserCockRace

			if currentSeason != nil {
				userPosition = app.GetUserPositionInSeason(log, userID, *currentSeason)
				neighbors = app.GetUsersAroundPositionInSeason(log, userPosition, *currentSeason)
			} else {
				userPosition = app.GetUserPositionInLadder(log, userID)
				neighbors = app.GetUsersAroundPositionInLadder(log, userPosition)
			}

			// Формируем строки для соседей с учетом edge cases
			var contextLines []string
			var showTopDots, showBottomDots bool

			if userPosition == 14 {
				// Сразу после топ-13 - показываем только текущего и следующего
				showTopDots = false
				showBottomDots = len(neighbors) == 2 && userPosition < totalParticipants
			} else if userPosition >= totalParticipants-1 {
				// Последние 2 места
				showTopDots = userPosition > 14
				showBottomDots = false
			} else {
				// Обычный случай
				showTopDots = userPosition > 14
				showBottomDots = userPosition < totalParticipants
			}

			startPos := userPosition - len(neighbors) + 1
			if userPosition == 14 {
				startPos = 14
			}

			for idx, neighbor := range neighbors {
				pos := startPos + idx
				isCurrentInContext := neighbor.UserID == userID
				normalizedNick := app.ResolveDisplayNickname(log, localizer, neighbor.UserID, neighbor.Nickname)
				formattedSize := EscapeMarkdownV2(FormatDickSize(int(neighbor.TotalSize)))
				posEmoji := GetPlaceEmojiForContext(pos, isCurrentInContext)

				if isCurrentInContext {
					contextLines = append(contextLines, app.localization.Localize(localizer, MsgCockRaceContextSelected, map[string]any{
						"PlaceEmoji": posEmoji,
						"Username":   EscapeMarkdownV2(normalizedNick),
						"Size":       formattedSize,
					}))
				} else {
					contextLines = append(contextLines, app.localization.Localize(localizer, MsgCockRaceContextDefault, map[string]any{
						"PlaceEmoji": posEmoji,
						"Username":   EscapeMarkdownV2(normalizedNick),
						"Size":       formattedSize,
					}))
				}
			}

			// Добавляем контекст с соседями
			var contextBlock string
			if showTopDots && showBottomDots {
				dots := app.localization.Localize(localizer, CommonDots, nil)
				contextBlock = "\n" + dots + "\n" + strings.Join(contextLines, "\n") + "\n" + dots
			} else if showTopDots {
				dots := app.localization.Localize(localizer, CommonDots, nil)
				contextBlock = "\n" + dots + "\n" + strings.Join(contextLines, "\n")
			} else if showBottomDots {
				dots := app.localization.Localize(localizer, CommonDots, nil)
				contextBlock = "\n" + strings.Join(contextLines, "\n") + "\n" + dots
			} else {
				contextBlock = "\n" + strings.Join(contextLines, "\n")
			}

			others = append(others, contextBlock)
		} else {
			others = append(others, app.localization.Localize(localizer, MsgCockScoreboardNotFound, nil))
		}
	}

	// Формируем нижнюю строку с информацией о текущем сезоне
	var footerLine string
	var seasonNum int
	var seasonWord string

	if currentSeason != nil {
		now := datetime.NowTime()
		startDateFormatted := EscapeMarkdownV2(currentSeason.StartDate.Format("02.01.2006"))
		endDateFormatted := EscapeMarkdownV2(currentSeason.EndDate.Format("02.01.2006"))
		timeRemaining := FormatTimeRemaining(app.localization, localizer, currentSeason.EndDate, now)

		seasonNum = currentSeason.SeasonNum
		seasonWord = app.localization.Localize(localizer, UnitSeasonGenitive, map[string]any{"Count": seasonNum})

		footerLine = app.localization.Localize(localizer, MsgCockRaceFooterActiveSeason, map[string]any{
			"SeasonNum": seasonNum,
			"StartDate": startDateFormatted,
			"EndDate":   endDateFormatted,
			"Remaining": EscapeMarkdownV2(timeRemaining),
		})
	} else {
		seasonNum = 1
		seasonWord = app.localization.Localize(localizer, UnitSeasonGenitive, map[string]any{"Count": seasonNum})
		footerLine = app.localization.Localize(localizer, MsgCockRaceFooterNoSeason, map[string]any{
			"StartDate": seasonStart,
		})
	}

	if len(others) != 0 {
		template := MsgCockRaceScoreboardTemplate
		if !showDescription {
			template = MsgCockRaceScoreboardTemplateNoDesc
		}
		return app.localization.Localize(localizer, template, map[string]any{
			"Participants": totalParticipants,
			"Winners":      strings.Join(winners, "\n"),
			"Others":       strings.Join(others, "\n"),
			"Footer":       footerLine,
			"SeasonNum":    seasonNum,
			"SeasonWord":   seasonWord,
		})
	} else {
		template := MsgCockRaceScoreboardWinnersTemplate
		if !showDescription {
			template = MsgCockRaceScoreboardWinnersTemplateNoDesc
		}
		return app.localization.Localize(localizer, template, map[string]any{
			"Participants": totalParticipants,
			"Winners":      strings.Join(winners, "\n"),
			"Footer":       footerLine,
			"SeasonNum":    seasonNum,
			"SeasonWord":   seasonWord,
		})
	}
}

func (app *Application) GenerateCockLadderScoreboard(log *logging.Logger, localizer *i18n.Localizer, userID int64, sizes []UserCockRace, totalParticipants int, showDescription bool) string {
	var winners []string
	var others []string
	isUserInScoreboard := false

	for index, user := range sizes {
		isCurrentUser := user.UserID == userID
		emoji := GetPlaceEmoji(index+1, isCurrentUser)

		if isCurrentUser {
			isUserInScoreboard = true
		}

		// Нормализуем username с учетом скрытия
		normalizedNickname := app.ResolveDisplayNickname(log, localizer, user.UserID, user.Nickname)

		var scoreboardLine string
		if isCurrentUser {
			scoreboardLine = app.localization.Localize(localizer, MsgCockLadderScoreboardSelected, map[string]any{
				"PlaceEmoji": emoji,
				"Username":   EscapeMarkdownV2(normalizedNickname),
				"Size":       EscapeMarkdownV2(FormatDickSize(int(user.TotalSize))),
			})
		} else {
			scoreboardLine = app.localization.Localize(localizer, MsgCockLadderScoreboardDefault, map[string]any{
				"PlaceEmoji": emoji,
				"Username":   EscapeMarkdownV2(normalizedNickname),
				"Size":       EscapeMarkdownV2(FormatDickSize(int(user.TotalSize))),
			})
		}

		if index < 3 {
			winners = append(winners, scoreboardLine)
		} else {
			others = append(others, scoreboardLine)
		}
	}

	if !isUserInScoreboard {
		if cock := app.GetUserAggregatedCock(log, userID); cock != nil {
			// Получаем позицию пользователя в ладдере и соседей
			userPosition := app.GetUserPositionInLadder(log, userID)
			neighbors := app.GetUsersAroundPositionInLadder(log, userPosition)

			// Формируем строки для соседей с учетом edge cases
			var contextLines []string
			var showTopDots, showBottomDots bool

			if userPosition == 14 {
				// Сразу после топ-13 - показываем только текущего и следующего
				showTopDots = false
				showBottomDots = len(neighbors) == 2 && userPosition < totalParticipants
			} else if userPosition >= totalParticipants-1 {
				// Последние 2 места
				showTopDots = userPosition > 14
				showBottomDots = false
			} else {
				// Обычный случай
				showTopDots = userPosition > 14
				showBottomDots = userPosition < totalParticipants
			}

			startPos := userPosition - len(neighbors) + 1
			if userPosition == 14 {
				startPos = 14
			}

			for idx, neighbor := range neighbors {
				pos := startPos + idx
				isCurrentInContext := neighbor.UserID == userID
				normalizedNick := app.ResolveDisplayNickname(log, localizer, neighbor.UserID, neighbor.Nickname)
				formattedSize := EscapeMarkdownV2(FormatDickSize(int(neighbor.TotalSize)))
				posEmoji := GetPlaceEmojiForContext(pos, isCurrentInContext)

				if isCurrentInContext {
					contextLines = append(contextLines, app.localization.Localize(localizer, MsgCockLadderContextSelected, map[string]any{
						"PlaceEmoji": posEmoji,
						"Username":   EscapeMarkdownV2(normalizedNick),
						"Size":       formattedSize,
					}))
				} else {
					contextLines = append(contextLines, app.localization.Localize(localizer, MsgCockLadderContextDefault, map[string]any{
						"PlaceEmoji": posEmoji,
						"Username":   EscapeMarkdownV2(normalizedNick),
						"Size":       formattedSize,
					}))
				}
			}

			// Добавляем контекст с соседями
			var contextBlock string
			if showTopDots && showBottomDots {
				dots := app.localization.Localize(localizer, CommonDots, nil)
				contextBlock = "\n" + dots + "\n" + strings.Join(contextLines, "\n") + "\n" + dots
			} else if showTopDots {
				dots := app.localization.Localize(localizer, CommonDots, nil)
				contextBlock = "\n" + dots + "\n" + strings.Join(contextLines, "\n")
			} else if showBottomDots {
				dots := app.localization.Localize(localizer, CommonDots, nil)
				contextBlock = "\n" + strings.Join(contextLines, "\n") + "\n" + dots
			} else {
				contextBlock = "\n" + strings.Join(contextLines, "\n")
			}

			others = append(others, contextBlock)
		} else {
			others = append(others, app.localization.Localize(localizer, MsgCockScoreboardNotFound, nil))
		}
	}

	if len(others) != 0 {
		template := MsgCockLadderScoreboardTemplate
		if !showDescription {
			template = MsgCockLadderScoreboardTemplateNoDesc
		}
		return app.localization.Localize(localizer, template, map[string]any{
			"Participants": totalParticipants,
			"Winners":      strings.Join(winners, "\n"),
			"Others":       strings.Join(others, "\n"),
		})
	} else {
		template := MsgCockLadderScoreboardWinnersTemplate
		if !showDescription {
			template = MsgCockLadderScoreboardWinnersTemplateNoDesc
		}
		return app.localization.Localize(localizer, template, map[string]any{
			"Participants": totalParticipants,
			"Winners":      strings.Join(winners, "\n"),
		})
	}
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

		// Для мест 4+ добавляем номер места (точка экранирована для MarkdownV2)
		// Номер жирный для текущего пользователя
		if isCurrentUser {
			return fmt.Sprintf("%s *%d*\\.", emoji, place)
		}
		return fmt.Sprintf("%s %d\\.", emoji, place)
	}
}

// GetPlaceEmojiForContext возвращает эмодзи для контекста (пользователи вне топ-13)
// Параметр bold определяет, делать ли номер жирным (для текущего пользователя)
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

func LuckLabel(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, luck float64) string {
	switch {
	case luck >= 1.98: // типа бог рандома :)
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
	case luck >= 0.2: // адовый тильт
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

// IrkLabel возвращает краткое описание ИРК (0.0-1.0+)
func IrkLabel(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, irk float64) string {
	irk = clamp01(irk)

	bucket := int(math.Floor(irk * 10)) // 0..9
	if irk >= 1.0 {
		bucket = 10
	}

	labels := [...]string{
		IrkLabelZero,      // 0.0..0.099
		IrkLabelMinimal,   // 0.1..0.199
		IrkLabelVerySmall, // 0.2..0.299
		IrkLabelSmall,     // 0.3..0.399
		IrkLabelReduced,   // 0.4..0.499
		IrkLabelAverage,   // 0.5..0.599
		IrkLabelIncreased, // 0.6..0.699
		IrkLabelLarge,     // 0.7..0.799
		IrkLabelVeryLarge, // 0.8..0.899
		IrkLabelMaximum,   // 0.9..0.999
		IrkLabelUltimate,  // 1.0
	}

	return localizationManager.Localize(localizer, labels[bucket], nil)
}

// GrowthSpeedLabel возвращает описание скорости изменения
// Скорость всегда положительная (абсолютное значение), показывает интенсивность изменения
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

// FormatGrowthSpeed форматирует скорость роста кока (в см/день) с 1 знаком после запятой
// Скорость всегда положительная (абсолютное значение), как на спидометре
func FormatGrowthSpeed(speed float64) string {
	p := message.NewPrinter(language.Russian)
	return p.Sprintf("%.1f", speed)
}

// FormatTimeRemaining форматирует оставшееся время до конца периода
// Возвращает строку типа "1 месяц 3 дня" или "14 дней"
func FormatTimeRemaining(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, endDate time.Time, now time.Time) string {
	duration := endDate.Sub(now)
	daysRemaining := int(duration.Hours() / 24)

	if daysRemaining < 0 {
		return localizationManager.Localize(localizer, UnitDay, map[string]any{"Count": 0})
	}

	// Если больше месяца, показываем месяцы + дни
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

	// Если меньше месяца, показываем только дни
	return localizationManager.Localize(localizer, UnitDay, map[string]any{"Count": daysRemaining})
}

// FormatUserPullingPeriod форматирует период с первого кока пользователя
// Формат: "2 года, 10 месяцев и 3 дня (с 27.02.2020)"
func FormatUserPullingPeriod(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, firstCockDate time.Time, now time.Time) string {
	years := now.Year() - firstCockDate.Year()
	months := int(now.Month()) - int(firstCockDate.Month())
	days := now.Day() - firstCockDate.Day()

	// Корректируем если дни отрицательные
	if days < 0 {
		months--
		// Берем количество дней в предыдущем месяце
		prevMonth := now.AddDate(0, -1, 0)
		daysInPrevMonth := time.Date(prevMonth.Year(), prevMonth.Month()+1, 0, 0, 0, 0, 0, prevMonth.Location()).Day()
		days += daysInPrevMonth
	}

	// Корректируем если месяцы отрицательные
	if months < 0 {
		years--
		months += 12
	}

	// Форматируем дату первого кока
	dateStr := firstCockDate.Format("02.01.2006")

	var parts []string

	// Добавляем годы если есть
	if years > 0 {
		parts = append(parts, localizationManager.Localize(localizer, UnitYear, map[string]any{"Count": years}))
	}

	// Добавляем месяцы если есть
	if months > 0 {
		parts = append(parts, localizationManager.Localize(localizer, UnitMonth, map[string]any{"Count": months}))
	}

	// Добавляем дни если есть (или если нет ничего больше)
	if days > 0 || len(parts) == 0 {
		parts = append(parts, localizationManager.Localize(localizer, UnitDay, map[string]any{"Count": days}))
	}

	// Собираем строку
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
	allAchievements []database.Achievement,
	userAchievements map[string]*database.DocumentUserAchievement,
	page int,
	itemsPerPage int,
) (string, int, int, int) {
	// Сортируем достижения: сначала выполненные, затем в порядке определения
	type AchievementWithStatus struct {
		Achievement database.Achievement
		UserAch     *database.DocumentUserAchievement
		IsCompleted bool
	}

	achievementsWithStatus := make([]AchievementWithStatus, 0, len(allAchievements))
	completedCount := 0
	totalRespects := 0

	for _, ach := range allAchievements {
		userAch, exists := userAchievements[ach.ID]
		isCompleted := exists && userAch.Completed

		achievementsWithStatus = append(achievementsWithStatus, AchievementWithStatus{
			Achievement: ach,
			UserAch:     userAch,
			IsCompleted: isCompleted,
		})

		if isCompleted {
			completedCount++
			totalRespects += ach.Respects
		}
	}

	// Сортируем: выполненные в начало
	sort.Slice(achievementsWithStatus, func(i, j int) bool {
		if achievementsWithStatus[i].IsCompleted != achievementsWithStatus[j].IsCompleted {
			return achievementsWithStatus[i].IsCompleted
		}
		return false // Остальные в порядке определения
	})

	// Вычисляем пагинацию
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

	// Генерируем текст для текущей страницы
	var lines []string
	for i := startIdx; i < endIdx; i++ {
		achStatus := achievementsWithStatus[i]
		line := FormatAchievementLine(localizationManager, localizer, achStatus.Achievement, achStatus.UserAch, achStatus.IsCompleted)
		lines = append(lines, line)
	}

	achievementsList := strings.Join(lines, "\n")

	// Вычисляем процент
	percentComplete := 0
	if len(allAchievements) > 0 {
		percentComplete = (completedCount * 100) / len(allAchievements)
	}

	return achievementsList, completedCount, totalRespects, percentComplete
}

// FormatAchievementLine форматирует одну строку достижения
func FormatAchievementLine(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, ach database.Achievement, userAch *database.DocumentUserAchievement, isCompleted bool) string {
	escapedName := EscapeMarkdownV2(localizationManager.Localize(localizer, ach.Name, nil))
	escapedDesc := EscapeMarkdownV2(localizationManager.Localize(localizer, ach.Description, nil))

	if isCompleted {
		// Выполненное достижение
		return localizationManager.Localize(localizer, MsgAchievementCompleted, map[string]any{
			"Emoji":       ach.Emoji,
			"Name":        escapedName,
			"Description": escapedDesc,
		})
	} else if userAch != nil && userAch.Progress > 0 && ach.MaxProgress > 0 {
		// В процессе выполнения
		return localizationManager.Localize(localizer, MsgAchievementInProgress, map[string]any{
			"Emoji":       ach.Emoji,
			"Name":        escapedName,
			"Progress":    userAch.Progress,
			"Max":         ach.MaxProgress,
			"Description": escapedDesc,
		})
	} else {
		// Не выполнено
		return localizationManager.Localize(localizer, MsgAchievementNotCompleted, map[string]any{
			"Emoji":       ach.Emoji,
			"Name":        escapedName,
			"Description": escapedDesc,
		})
	}
}
