package application

import (
	"fmt"
	"strings"
	"time"
)

const (
	CommonDots = "\\.\\.\\."

	MsgCockScoreboardNotFound = "\n🥀 *Тебе соболезнуем\\.\\.\\. потому что не смотрел на кок\\!*"

	MsgCockSize                    = "Мой кок: *%sсм* %s"
	MsgCockRulerScoreboardDefault  = "%s @%s *%sсм* %s"
	MsgCockRulerScoreboardSelected = "%s *@%s %sсм* %s"
	MsgCockRaceScoreboardDefault  = "%s @%s *%sсм*"
	MsgCockRaceScoreboardSelected = "%s *@%s %sсм*"

	MsgCockLadderScoreboardDefault  = "%s @%s *%sсм*"
	MsgCockLadderScoreboardSelected = "%s *@%s %sсм*"

	MsgCockRulerScoreboardTemplate = `*Линейка коков:*
👥 Участников: *%d*

🏆 Победители в номинации:

%s

🥀 Остальным соболезнуем:

%s

_Линейка коков – чистый рандом, сегодня ты бог, завтра ты лох\. Все коки сбрасываются каждые сутки по МСК\!_`

	MsgCockRulerScoreboardWinnersTemplate = `*Линейка коков:*
👥 Участников: *%d*

🏆 Победители в номинации:

%s

_Линейка коков – чистый рандом, сегодня ты бог, завтра ты лох\. Все коки сбрасываются каждые сутки по МСК\!_`

	MsgCockRaceScoreboardTemplate = `*Участники гонки коков:*
👥 Участников в сезоне: *%d*

🏆 Победители в номинации:

%s

🥀 Остальным соболезнуем:

%s

_Гонка коков – это соревнование, в котором коки участников суммируются за весь сезон\. Период обновления коков – сутки_
  
🚀 Текущий сезон гонки коков стартовал *%s*`

	MsgCockRaceScoreboardWinnersTemplate = `*Участники гонки коков:*
👥 Участников в сезоне: *%d*

🏆 Победители в номинации:

%s

_Гонка коков – это соревнование, в котором коки участников суммируются за весь сезон\. Период обновления коков – сутки_
  
🚀 Текущий сезон гонки коков стартовал *%s*`

	MsgCockLadderScoreboardTemplate = `*Ладдер коков:*
👥 Всего участников: *%d*

🏆 Лидеры кок–ладдера:

%s

🥀 Медленно, но верно поднимающиеся:

%s

_Ладдер коков – глобальный рейтинг участников по суммарному размеру кока за все время\. Покоряй вершины\!_`

	MsgCockLadderScoreboardWinnersTemplate = `*Ладдер коков:*
👥 Всего участников: *%d*

🏆 Лидеры кок–ладдера:

%s

_Ладдер коков – глобальный рейтинг участников по суммарному размеру кока за все время\. Покоряй вершины\!_`

	MsgCockDynamicsTemplate = `
📊 *Общая динамика коков*

Общий посчитанный кок: *%[1]s см* 🤭
Всего кокеров: *%[2]s* 🫡
Всего дёрнуто коков: *%[26]s* ✊🏻

День самого большого кока: *%[21]s*, нарастили аж *%[22]sсм* 🍾

Средний кок в системе _(5 дн.)_: *%[3]s см* %[4]s
Медиана кока в системе _(5 дн.)_: *%[5]s см* %[6]s

Соотношение коков _(5 дн.)_: 💪 *%[19]s%%* 🤏 *%[20]s%%*

📊 *Персональная динамика кока*

ИРК (Индекс Размера Кока): __*%[10]s*__
Коэффициент везения _(5 дн.)_: *%[28]s* %[29]s
Волатильность кока _(5 дн.)_: *%[30]s* %[31]s

Общий посчитанный кок: *%[7]s см* 🤯
В среднем размер кока _(5 дн.)_: *%[8]s см* %[9]s
Самый большой кок был: *%[11]s см* %[12]s (*%[13]s*)
Всего дёрнуто коков: *%[27]s* ✊🏻

Процент доминирования: *%[23]s%%* 👑

🏆 *Сезонные достижения*

Побед в сезонах: *%[24]s* 🎖️
Кок-респект: *%[25]s* 🚀

📈 *Кок-активы*

%[14]s Вчерашняя динамика: *%[15]s%%* (*%[16]s см*)
%[17]s Средний дневной прирост _(5 дн.)_: *%[18]s см/день*`

	MsgCockSeasonTemplate = `*Сезон коков* \(🟡 Текущий\)
⏱️ Период: *%[2]s \- %[3]s*

🔮 Претенденты сезона:

%[1]s`

	MsgCockSeasonWithWinnersTemplate = `*Сезон коков* \(🟢 Завершён\)
⏱️ Период: *%[2]s \- %[3]s*

🎖 Победители сезона:

%[1]s`

	MsgCockSeasonTemplateFooter = `*Каждый сезон* соревнуются игроки, отправляя свои коки в течение определенного периода времени\!

ℹ️ *Три участника* с наибольшим суммарным размером кока получают максимальные *кок\-респекты™*\!`

	MsgCockSeasonWinnerTemplate = "%[1]s *@%[2]s* с коком *%[3]sсм*"
	
	MsgCockSeasonNoSeasonsTemplate = `*Сезоны коков*\n\nВ данный момент нет активных сезонов\. Следите за обновлениями\!`
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

	/* Соотношение коков */

	totalBigCockRatio float64,
	totalSmallCockRatio float64,

	/* Самый большой кок */

	totalMaxCockDate time.Time,
	totalMaxCock int,

	/* % доминирование */

	userDominancePercent float64,

	/* Сезонные достижения */

	userSeasonWins int,
	userCockRespect int,

	/* Всего дёрнуто коков */

	totalCocksCount int,
	userCocksCount int,

	/* Коэффициент везения и волатильность */

	userLuckCoefficient float64,
	userVolatility float64,
) string {
	var userYesterdayChangePercentEmoji string
	var userYesterdayChangePercentSymbol string
	if userYesterdayChangePercent >= 0 {
		userYesterdayChangePercentEmoji = "🟩"
		userYesterdayChangePercentSymbol = "+"
	} else {
		userYesterdayChangePercentEmoji = "🟥"
		userYesterdayChangePercentSymbol = ""
	}

	var userDailyGrowthEmoji string
	var userDailyGrowthSymbol string
	if userDailyGrowth >= 0 {
		userDailyGrowthEmoji = "🟩"
		userDailyGrowthSymbol = "+"
	} else {
		userDailyGrowthEmoji = "🟥"
		userDailyGrowthSymbol = ""
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
		EscapeMarkdownV2(FormatDickSize(userMaxCock)), EmojiFromSize(userMaxCock), userMaxCockDate.Local().Format("02.01.06"),

		/* Кок-активы */
		userYesterdayChangePercentEmoji, fmt.Sprintf("%s%s", userYesterdayChangePercentSymbol, FormatDickPercent(userYesterdayChangePercent)), fmt.Sprintf("%s%s", userYesterdayChangePercentSymbol, FormatDickSize(userYesterdayChangeCock)),
		userDailyGrowthEmoji, fmt.Sprintf("%s%s", userDailyGrowthSymbol, FormatDickPercent(userDailyGrowth)),

		/* Соотношение коков */

		FormatDickPercent(totalBigCockRatio), FormatDickPercent(totalSmallCockRatio),

		/* Самый большой кок */
		totalMaxCockDate.Local().Format("02.01.06"), FormatDickSize(totalMaxCock),

		/* % Доминирования */
		FormatDickPercent(userDominancePercent),

		/* Сезонные достижения */
		FormatDickSize(userSeasonWins),
		FormatDickSize(userCockRespect),

		/* Всего дёрнуто коков */
		EscapeMarkdownV2(FormatDickSize(totalCocksCount)),
		EscapeMarkdownV2(FormatDickSize(userCocksCount)),

		/* Коэффициент везения и волатильность */
		EscapeMarkdownV2(FormatLuckCoefficient(userLuckCoefficient)), LuckEmoji(userLuckCoefficient),
		EscapeMarkdownV2(FormatVolatility(userVolatility)), VolatilityDisplay(userVolatility),
	)
}

func GetMedalByPosition(position int) string {
	switch position {
	case 0:
		return "🥇"
	case 1:
		return "🥈"
	case 2:
		return "🥉"
	default:
		return ""
	}
}

func NewMsgCockSeasonTemplate(pretenders string, startDate, endDate string) string {
	return fmt.Sprintf(
		MsgCockSeasonTemplate,
		pretenders,
		startDate,
		endDate,
	)
}

func NewMsgCockSeasonWithWinnersTemplate(winners string, startDate, endDate string) string {
	return fmt.Sprintf(
		MsgCockSeasonWithWinnersTemplate,
		winners,
		startDate,
		endDate,
	)
}

func NewMsgCockSeasonWinnerTemplate(medal, nickname, totalSize string) string {
	return fmt.Sprintf(
		MsgCockSeasonWinnerTemplate,
		medal,
		EscapeMarkdownV2(nickname),
		EscapeMarkdownV2(totalSize),
	)
}

func NewMsgCockSeasonTemplateFooter() string {
	return MsgCockSeasonTemplateFooter
}

func NewMsgCockSeasonNoSeasonsTemplate() string {
	return MsgCockSeasonNoSeasonsTemplate
}

func NewMsgCockSeasonsFullText(seasons []CockSeason, totalSeasonsCount int, getSeasonWinners func(CockSeason) []SeasonWinner) string {
	if len(seasons) == 0 {
		return NewMsgCockSeasonNoSeasonsTemplate()
	}
	
	var seasonBlocks []string
	
	// Проходим сезоны в обратном порядке (от нового к старому)
	for i := len(seasons) - 1; i >= 0; i-- {
		season := seasons[i]
		startDate := EscapeMarkdownV2(season.StartDate.Format("02.01.2006"))
		endDate := EscapeMarkdownV2(season.EndDate.Format("02.01.2006"))
		
		winners := getSeasonWinners(season)
		var winnerLines []string
		
		for _, winner := range winners {
			medal := GetMedalByPosition(winner.Place - 1)
			// Нормализуем nickname (генерируем анонимное имя если пустой)
			normalizedNickname := NormalizeUsername(winner.Nickname, winner.UserID)
			line := NewMsgCockSeasonWinnerTemplate(
				medal,
				normalizedNickname,
				FormatDickSize(int(winner.TotalSize)),
			)
			winnerLines = append(winnerLines, line)
		}
		
		winnersText := strings.Join(winnerLines, "\n")
		
		var seasonBlock string
		if season.IsActive {
			seasonBlock = NewMsgCockSeasonTemplate(winnersText, startDate, endDate)
		} else {
			seasonBlock = NewMsgCockSeasonWithWinnersTemplate(winnersText, startDate, endDate)
		}
		
		seasonBlocks = append(seasonBlocks, seasonBlock)
	}
	
	allSeasonsText := strings.Join(seasonBlocks, "\n\n———————————————\n\n")
	footer := NewMsgCockSeasonTemplateFooter()
	
	var finalText string
	if totalSeasonsCount > len(seasons) {
		trimInfo := fmt.Sprintf("\n\n📋 _Показаны последние %d из %d сезонов_", len(seasons), totalSeasonsCount)
		finalText = allSeasonsText + "\n\n" + footer + EscapeMarkdownV2(trimInfo)
	} else {
		finalText = allSeasonsText + "\n\n" + footer
	}
	
	return finalText
}
