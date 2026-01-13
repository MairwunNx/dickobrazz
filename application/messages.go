package application

import (
	"dickobrazz/application/datetime"
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

>📖 О линейке коков:
>
>Линейка коков – это daily рейтинг чистого рандома\. Размеры генерируются случайно каждый день \(от 0 до 61 см\) и сбрасываются в полночь по МСК\. Никаких накоплений – только удача здесь и сейчас\!
>
>🎲 Сегодня ты топ, завтра ты дно – рандом решает, кому повезло\!||`

	MsgCockRulerScoreboardWinnersTemplate = `*Линейка коков:*
👥 Участников: *%d*

🏆 Победители в номинации:

%s

>📖 О линейке коков:
>
>Линейка коков – это daily рейтинг чистого рандома\. Размеры генерируются случайно каждый день \(от 0 до 61 см\) и сбрасываются в полночь по МСК\. Никаких накоплений – только удача здесь и сейчас\!
>
>🎲 Сегодня ты топ, завтра ты дно – рандом решает, кому повезло\!||`

	// Версии без описаний для опытных пользователей
	MsgCockRulerScoreboardTemplateNoDesc = `*Линейка коков:*
👥 Участников: *%d*

🏆 Победители в номинации:

%s

🥀 Остальным соболезнуем:

%s`

	MsgCockRulerScoreboardWinnersTemplateNoDesc = `*Линейка коков:*
👥 Участников: *%d*

🏆 Победители в номинации:

%s`

	MsgCockRaceScoreboardTemplate = `*Участники гонки коков %[5]d %[6]s:*
👥 Участников в сезоне: *%[1]d*

🏆 Победители в номинации:

%[2]s

🥀 Остальным соболезнуем:

%[3]s

>📖 О гонке коков:
>
>Гонка коков – это сезонное соревнование длиной 3 месяца\. Измеряй свой кок ежедневно, все результаты суммируются автоматически\. Побеждают три участника с максимальным накопленным размером за весь сезон\.
>
>С началом нового сезона все коки сбрасываются, и гонка начинается заново для всех\.
>
>💡 Совет: Измеряй кок каждый день, чтобы максимизировать свои шансы на победу\!||

%[4]s`

	MsgCockRaceScoreboardWinnersTemplate = `*Участники гонки коков %[4]d %[5]s:*
👥 Участников в сезоне: *%[1]d*

🏆 Победители в номинации:

%[2]s

>📖 О гонке коков:
>
>Гонка коков – это сезонное соревнование длиной 3 месяца\. Измеряй свой кок ежедневно, все результаты суммируются автоматически\. Побеждают три участника с максимальным накопленным размером за весь сезон\.
>
>С началом нового сезона все коки сбрасываются, и гонка начинается заново для всех\.
>
>💡 Совет: Измеряй кок каждый день, чтобы максимизировать свои шансы на победу\!||

%[3]s`

	// Версии без описаний для опытных пользователей
	MsgCockRaceScoreboardTemplateNoDesc = `*Участники гонки коков %[5]d %[6]s:*
👥 Участников в сезоне: *%[1]d*

🏆 Победители в номинации:

%[2]s

🥀 Остальным соболезнуем:

%[3]s

%[4]s`

	MsgCockRaceScoreboardWinnersTemplateNoDesc = `*Участники гонки коков %[4]d %[5]s:*
👥 Участников в сезоне: *%[1]d*

🏆 Победители в номинации:

%[2]s

%[3]s`

	MsgCockLadderScoreboardTemplate = `*Ладдер коков:*
👥 Всего участников: *%d*

🏆 Лидеры кок–ладдера:

%s

🥀 Медленно, но верно поднимающиеся:

%s

>📖 О ладдере коков:
>
>Ладдер коков – это твой вечный путь к славе\. Здесь суммируется каждый кок за всю историю твоего участия\. В отличие от дневной линейки и сезонной гонки, ладдер никогда не обнуляется\.
>
>🔥 Топ ладдера – это легенды, измеряющие коки с первого дня\. Стань одним из них\!||`

	MsgCockLadderScoreboardWinnersTemplate = `*Ладдер коков:*
👥 Всего участников: *%d*

🏆 Лидеры кок–ладдера:

%s

>📖 О ладдере коков:
>
>Ладдер коков – это твой вечный путь к славе\. Здесь суммируется каждый кок за всю историю твоего участия\. В отличие от дневной линейки и сезонной гонки, ладдер никогда не обнуляется\.
>
>🔥 Топ ладдера – это легенды, измеряющие коки с первого дня\. Стань одним из них\!||`

	// Версии без описаний для опытных пользователей
	MsgCockLadderScoreboardTemplateNoDesc = `*Ладдер коков:*
👥 Всего участников: *%d*

🏆 Лидеры кок–ладдера:

%s

🥀 Медленно, но верно поднимающиеся:

%s`

	MsgCockLadderScoreboardWinnersTemplateNoDesc = `*Ладдер коков:*
👥 Всего участников: *%d*

🏆 Лидеры кок–ладдера:

%s`

	MsgCockDynamicsTemplate = `
📊 *Общая динамика коков*

Общий посчитанный кок: *%[1]s см* 🤭
Всего кокеров: *%[2]s* 🫡
Всего дёрнуто коков: *%[26]s* ✊🏻

День самого большого кока: *%[21]s*, нарастили аж *%[22]s см* 🍾

Средний кок в системе _(5 коков)_: *%[3]s см* %[4]s
Медиана кока в системе _(5 коков)_: *%[5]s см* %[6]s

Соотношение коков _(5 коков)_: 💪 *%[19]s%%* 🤏 *%[20]s%%*

📊 *Персональная динамика кока*

Общий посчитанный кок: *%[7]s см* 🤯
Всего дёрнуто коков: *%[27]s* ✊🏻

ИРК (Индекс Размера Кока): __*%[10]s*__ _(%[32]s)_
В среднем размер кока _(5 коков)_: *%[8]s см* %[9]s
Самый большой кок был: *%[11]s см* %[12]s (*%[13]s*)

Коэффициент везения _(5 коков)_: *%[28]s* %[29]s
Волатильность кока _(5 коков)_: *%[30]s* %[31]s

Процент доминирования: *%[23]s%%* 👑
Скорость роста кока _(5 коков)_: *%[34]s см/день* %[35]s

🏆 *Сезонные достижения*

Побед в сезонах: *%[24]s* 🎖️
Кок-респект: *%[25]s* 🚀

📈 *Кок-активы*

%[14]s Дневная динамика: *%[15]s%%* (*%[16]s см*)
%[33]s Динамика за 5 коков: *%[17]s%%* (*%[18]s см*)`

	MsgCockSeasonTemplate = `*Сезон коков %[4]d* \(🟡 Текущий\)
⏱️ Период: *%[2]s \- %[3]s*

🔮 Претенденты сезона:

%[1]s`

	MsgCockSeasonWithWinnersTemplate = `*Сезон коков %[4]d* \(🟢 Завершён\)
⏱️ Период: *%[2]s \- %[3]s*

🎖 Победители сезона:

%[1]s`

	MsgCockSeasonTemplateFooter = `>📖 О сезонах коков:
>
>Сезоны коков – это 3\-месячная битва за звание лучшего кокера\. Измеряй каждый день, суммируй результаты и борись за топ\-3\. Победители получают легендарные кок\-респекты™, которые можно обменять на мерч\!
>
>🔥 История помнит только победителей – стань одним из них\!||`

	MsgCockSeasonWinnerTemplate = "%[1]s *@%[2]s* с коком *%[3]s см*"
	
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
	userFiveCocksChangePercent float64,
	userFiveCocksChangeCock int,

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
	
	/* Средняя скорость прироста */
	
	userGrowthSpeed float64,
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

	var userFiveCocksChangeEmoji string
	var userFiveCocksChangeSymbol string
	if userFiveCocksChangePercent >= 0 {
		userFiveCocksChangeEmoji = "🟩"
		userFiveCocksChangeSymbol = "+"
	} else {
		userFiveCocksChangeEmoji = "🟥"
		userFiveCocksChangeSymbol = ""
	}

	return fmt.Sprintf(
		MsgCockDynamicsTemplate,

		/* 1-2: Общая динамика коков */
		EscapeMarkdownV2(FormatDickSize(totalCock)),           // %[1]s
		EscapeMarkdownV2(FormatDickSize(totalUsers)),          // %[2]s

		/* 3-6: Средний и медианный кок */
		EscapeMarkdownV2(FormatDickSize(totalAvgCock)),        // %[3]s
		EmojiFromSize(totalAvgCock),                           // %[4]s
		EscapeMarkdownV2(FormatDickSize(totalMedianCock)),     // %[5]s
		EmojiFromSize(totalMedianCock),                        // %[6]s

		/* 7-13: Персональная динамика кока */
		EscapeMarkdownV2(FormatDickSize(userTotalCock)),       // %[7]s
		EscapeMarkdownV2(FormatDickSize(userAvgCock)),         // %[8]s
		EmojiFromSize(userAvgCock),                            // %[9]s
		EscapeMarkdownV2(FormatDickIkr(userIrk)),              // %[10]s
		EscapeMarkdownV2(FormatDickSize(userMaxCock)),         // %[11]s
		EmojiFromSize(userMaxCock),                            // %[12]s
		userMaxCockDate.In(datetime.NowLocation()).Format("02.01.06"), // %[13]s

		/* 14-18: Кок-активы (дневная и 5 коков динамика) */
		userYesterdayChangePercentEmoji,                       // %[14]s
		fmt.Sprintf("%s%s", userYesterdayChangePercentSymbol, FormatDickPercent(userYesterdayChangePercent)), // %[15]s
		fmt.Sprintf("%s%s", userYesterdayChangePercentSymbol, FormatDickSize(userYesterdayChangeCock)),       // %[16]s
		fmt.Sprintf("%s%s", userFiveCocksChangeSymbol, FormatDickPercent(userFiveCocksChangePercent)),        // %[17]s
		fmt.Sprintf("%s%s", userFiveCocksChangeSymbol, FormatDickSize(userFiveCocksChangeCock)),              // %[18]s

		/* 19-20: Соотношение коков */
		FormatDickPercent(totalBigCockRatio),                  // %[19]s
		FormatDickPercent(totalSmallCockRatio),                // %[20]s

		/* 21-22: Самый большой кок */
		totalMaxCockDate.In(datetime.NowLocation()).Format("02.01.06"), // %[21]s
		FormatDickSize(totalMaxCock),                          // %[22]s

		/* 23: % Доминирования */
		FormatDickPercent(userDominancePercent),               // %[23]s

		/* 24-25: Сезонные достижения */
		FormatDickSize(userSeasonWins),                        // %[24]s
		FormatDickSize(userCockRespect),                       // %[25]s

		/* 26-27: Всего дёрнуто коков */
		EscapeMarkdownV2(FormatDickSize(totalCocksCount)),     // %[26]s
		EscapeMarkdownV2(FormatDickSize(userCocksCount)),      // %[27]s

		/* 28-31: Коэффициент везения и волатильность */
		EscapeMarkdownV2(FormatLuckCoefficient(userLuckCoefficient)), // %[28]s
		LuckDisplay(userLuckCoefficient),                      // %[29]s
		EscapeMarkdownV2(FormatVolatility(userVolatility)),    // %[30]s
		VolatilityDisplay(userVolatility),                     // %[31]s

		/* 32: Описание ИРК */
		IrkLabel(userIrk),                                     // %[32]s

		/* 33: Эмодзи динамики за 5 коков */
		userFiveCocksChangeEmoji,                              // %[33]s

		/* 34-35: Скорость прироста кока */
		EscapeMarkdownV2(FormatGrowthSpeed(userGrowthSpeed)),   // %[34]s
		GrowthSpeedDisplay(userGrowthSpeed),                    // %[35]s
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

func NewMsgCockSeasonTemplate(pretenders string, startDate, endDate string, seasonNum int) string {
	return fmt.Sprintf(
		MsgCockSeasonTemplate,
		pretenders,
		startDate,
		endDate,
		seasonNum,
	)
}

func NewMsgCockSeasonWithWinnersTemplate(winners string, startDate, endDate string, seasonNum int) string {
	return fmt.Sprintf(
		MsgCockSeasonWithWinnersTemplate,
		winners,
		startDate,
		endDate,
		seasonNum,
	)
}

func NewMsgCockSeasonWinnerTemplate(medal, nickname, totalSize string, respects int, showRespects bool) string {
	winnersLine := fmt.Sprintf(
		MsgCockSeasonWinnerTemplate,
		medal,
		EscapeMarkdownV2(nickname),
		EscapeMarkdownV2(totalSize),
	)
	
	// Показываем респекты только если showRespects = true (для завершенных сезонов)
	if showRespects {
		formattedRespects := EscapeMarkdownV2(FormatDickSize(respects))
		return fmt.Sprintf("%s *\\(\\+%s 🫡\\)*", winnersLine, formattedRespects)
	}
	
	return winnersLine
}

func NewMsgCockSeasonTemplateFooter() string {
	return MsgCockSeasonTemplateFooter
}

func NewMsgCockSeasonNoSeasonsTemplate() string {
	return MsgCockSeasonNoSeasonsTemplate
}

// NewMsgCockSeasonSinglePage генерирует текст для одной страницы сезона (постраничная навигация)
func NewMsgCockSeasonSinglePage(season CockSeason, getSeasonWinners func(CockSeason) []SeasonWinner, showDescription bool) string {
	startDate := EscapeMarkdownV2(season.StartDate.Format("02.01.2006"))
	endDate := EscapeMarkdownV2(season.EndDate.Format("02.01.2006"))
	
	winners := getSeasonWinners(season)
	var winnerLines []string
	
	for _, winner := range winners {
		medal := GetMedalByPosition(winner.Place - 1)
		normalizedNickname := NormalizeUsername(winner.Nickname, winner.UserID)
		respects := CalculateCockRespect(winner.Place)
		// Показываем респекты только для завершенных сезонов
		line := NewMsgCockSeasonWinnerTemplate(
			medal,
			normalizedNickname,
			FormatDickSize(int(winner.TotalSize)),
			respects,
			!season.IsActive, // showRespects = true только если сезон завершен
		)
		winnerLines = append(winnerLines, line)
	}
	
	winnersText := strings.Join(winnerLines, "\n")
	
	var seasonBlock string
	if season.IsActive {
		seasonBlock = NewMsgCockSeasonTemplate(winnersText, startDate, endDate, season.SeasonNum)
		// Футер показываем только для активного (текущего) сезона И если showDescription = true
		if showDescription {
			footer := NewMsgCockSeasonTemplateFooter()
			return seasonBlock + "\n\n" + footer
		}
		return seasonBlock
	} else {
		seasonBlock = NewMsgCockSeasonWithWinnersTemplate(winnersText, startDate, endDate, season.SeasonNum)
		return seasonBlock
	}
}

// MsgCockAchievementsTemplate - шаблон для списка достижений (первая страница с описанием)
const MsgCockAchievementsTemplate = `🏆 *Кок\-ачивки*
Выполнено: *%d/%d* _\(%d%%\)_ • 🌟 Респекты: *%d*

💡 _За каждую кок\-ачивку ты получаешь кок\-респекты™, которые скоро можно будет обменять на мерч в официальном магазине\!_

%s`

// MsgCockAchievementsTemplateOtherPages - шаблон для остальных страниц (без описания)
const MsgCockAchievementsTemplateOtherPages = `🏆 *Кок\-ачивки*
Выполнено: *%d/%d* _\(%d%%\)_ • 🌟 Респекты: *%d*

%s`
