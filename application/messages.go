package application

import (
	"fmt"
	"time"
)

const (
	CommonDots = "..."

	MsgCockScoreboardNotFound = "\n🥀 *Тебе соболезнуем... потому что не смотрел на кок!*"

	MsgCockSize                    = "Мой кок: *%dсм* %s"
	MsgCockRulerScoreboardDefault  = "%s @%s *%dсм* %s"
	MsgCockRulerScoreboardSelected = "%s *@%s %dсм* %s"
	MsgCockRulerScoreboardOut      = "\n🥀 *И %s твой кок %dсм* %s"

	MsgCockRaceScoreboardDefault  = "%s @%s *%sсм*"
	MsgCockRaceScoreboardSelected = "%s *@%s %sсм*"
	MsgCockRaceScoreboardOut      = "\n🥀 *И %s твой кок %sсм*"

	MsgCockRulerScoreboardTemplate = `*Линейка коков:*

🏆 Победители в номинации:

%s

🥀 Остальным соболезнуем:

%s

_Линейка коков – чистый рандом, сегодня ты бог, завтра ты лох. Все коки сбрасываются каждые сутки по МСК!_`

	MsgCockRulerScoreboardWinnersTemplate = `*Линейка коков:*

🏆 Победители в номинации:

%s

_Линейка коков – чистый рандом, сегодня ты бог, завтра ты лох. Все коки сбрасываются каждые сутки по МСК!_`

	MsgCockRaceScoreboardTemplate = `*Участники гонки коков:*

🏆 Победители в номинации:

%s

🥀 Остальным соболезнуем:

%s

_Гонка коков – это соревнование, в котором коки участников суммируются за весь сезон. Период обновления коков – сутки_
  
🚀 Текущий сезон гонки коков стартовал *16.12.2024*`

	MsgCockRaceScoreboardWinnersTemplate = `*Участники гонки коков:*

🏆 Победители в номинации:

%s

_Гонка коков – это соревнование, в котором коки участников суммируются за весь сезон. Период обновления коков – сутки_
  
🚀 Текущий сезон гонки коков стартовал *16.12.2024*`

	MsgCockDynamicsTemplate = `
📊 *Общая динамика коков*

Общий посчитанный кок: *%[1]s см* 🤭
Всего уникальных коков в системе: *%[2]s* 🫡

Средний кок в системе: *%[3]s см* %[4]s
Медиана кока в системе: *%[5]s см* %[6]s

📊 *Персональная динамика кока*

Общий посчитанный кок: *%[7]s см* 🤯
В среднем размер кока: *%[8]s см* %[9]s
ИРК (Индекс Размера Кока): %[10]s
Самый большой кок был: *%[11]s см* %[12]s (%[13]s)

📈 *Кок-активы*

%[14]s Вчерашняя динамика: *%[15]s%%* (*%[16]s см*) %[17]s
%[18]s Средний дневной прирост: *%[19]s см/день* %[20]s

⚠️ _Могут быть недоработки, динамика коков тестируется._
`
)

func NewMsgCockDynamicsTemplate(
	/* Общая динамика коков */

	totalCock int,
	totalUsers int,
	totalAvgCock int,
	totalMedianCock int,

	/* Персональная динамика кока */

	userTotalCock int,
	userAvgCock int,
	userIrk float64,
	userMaxCock int,
	userMaxCockDate time.Time,

	/* Кок-активы */

	userYesterdayChangePercent float64,
	userYesterdayChangeCock int,

	userDailyGrowth float64,
) string {
	var userYesterdayChangePercentEmoji string
	var userYesterdayChangePercentSymbol string
	var userYesterdayChangePercentEmojiEnd string
	if userYesterdayChangePercent >= 0 {
		userYesterdayChangePercentEmoji = "🟩"
		userYesterdayChangePercentSymbol = "+"
		userYesterdayChangePercentEmojiEnd = "🔺"
	} else {
		userYesterdayChangePercentEmoji = "🟥"
		userYesterdayChangePercentSymbol = "-"
		userYesterdayChangePercentEmojiEnd = "🔻"
	}

	var userDailyGrowthEmoji string
	var userDailyGrowthSymbol string
	var userDailyGrowthEmojiEnd string
	if userDailyGrowth >= 0 {
		userDailyGrowthEmoji = "🟩"
		userDailyGrowthSymbol = "+"
		userDailyGrowthEmojiEnd = "🔺"
	} else {
		userDailyGrowthEmoji = "🟥"
		userDailyGrowthSymbol = "-"
		userDailyGrowthEmojiEnd = "🔻"
	}

	return fmt.Sprintf(
		MsgCockDynamicsTemplate,

		/* Общая динамика коков */

		EscapeMarkdownV2(FormatDickSize(totalCock)),
		EscapeMarkdownV2(FormatDickSize(totalUsers)),
		EscapeMarkdownV2(FormatDickSize(totalAvgCock)), EmojiFromSize(totalAvgCock),
		EscapeMarkdownV2(FormatDickSize(totalMedianCock)), EmojiFromSize(totalMedianCock),

		/* Персональная динамика кока */

		EscapeMarkdownV2(FormatDickSize(userTotalCock)),
		EscapeMarkdownV2(FormatDickSize(userAvgCock)), EmojiFromSize(userAvgCock),
		EscapeMarkdownV2(FormatDickIkr(userIrk)),
		EscapeMarkdownV2(FormatDickSize(userMaxCock)), EscapeMarkdownV2(userMaxCockDate.Format("02.01.06")), EmojiFromSize(userMaxCock),

		/* Кок-активы */
		userYesterdayChangePercentEmoji, fmt.Sprintf("%s%s", EscapeMarkdownV2(userYesterdayChangePercentSymbol), EscapeMarkdownV2(FormatDickPercent(userYesterdayChangePercent))), EscapeMarkdownV2(FormatDickSize(userYesterdayChangeCock)), userYesterdayChangePercentEmojiEnd,
		userDailyGrowthEmoji, fmt.Sprintf("%s%s", EscapeMarkdownV2(userDailyGrowthSymbol), EscapeMarkdownV2(FormatDickPercent(userDailyGrowth))), userDailyGrowthEmojiEnd,
	)
}
