package application

import (
	"dickobrazz/application/datetime"
	"dickobrazz/application/localization"
	"fmt"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

const (
	CommonDots = "CommonDots"

	MsgCockScoreboardNotFound = "MsgCockScoreboardNotFound"

	MsgCockSize                    = "MsgCockSize"
	MsgCockRulerScoreboardDefault  = "MsgCockRulerScoreboardDefault"
	MsgCockRulerScoreboardSelected = "MsgCockRulerScoreboardSelected"
	MsgCockRaceScoreboardDefault   = "MsgCockRaceScoreboardDefault"
	MsgCockRaceScoreboardSelected  = "MsgCockRaceScoreboardSelected"

	MsgCockRulerContextDefault   = "MsgCockRulerContextDefault"
	MsgCockRulerContextSelected  = "MsgCockRulerContextSelected"
	MsgCockRaceContextDefault    = "MsgCockRaceContextDefault"
	MsgCockRaceContextSelected   = "MsgCockRaceContextSelected"
	MsgCockLadderContextDefault  = "MsgCockLadderContextDefault"
	MsgCockLadderContextSelected = "MsgCockLadderContextSelected"

	MsgCockLadderScoreboardDefault  = "MsgCockLadderScoreboardDefault"
	MsgCockLadderScoreboardSelected = "MsgCockLadderScoreboardSelected"

	MsgCockRulerScoreboardTemplate              = "MsgCockRulerScoreboardTemplate"
	MsgCockRulerScoreboardWinnersTemplate       = "MsgCockRulerScoreboardWinnersTemplate"
	MsgCockRulerScoreboardTemplateNoDesc        = "MsgCockRulerScoreboardTemplateNoDesc"
	MsgCockRulerScoreboardWinnersTemplateNoDesc = "MsgCockRulerScoreboardWinnersTemplateNoDesc"

	MsgCockRaceScoreboardTemplate              = "MsgCockRaceScoreboardTemplate"
	MsgCockRaceScoreboardWinnersTemplate       = "MsgCockRaceScoreboardWinnersTemplate"
	MsgCockRaceScoreboardTemplateNoDesc        = "MsgCockRaceScoreboardTemplateNoDesc"
	MsgCockRaceScoreboardWinnersTemplateNoDesc = "MsgCockRaceScoreboardWinnersTemplateNoDesc"

	MsgCockLadderScoreboardTemplate              = "MsgCockLadderScoreboardTemplate"
	MsgCockLadderScoreboardWinnersTemplate       = "MsgCockLadderScoreboardWinnersTemplate"
	MsgCockLadderScoreboardTemplateNoDesc        = "MsgCockLadderScoreboardTemplateNoDesc"
	MsgCockLadderScoreboardWinnersTemplateNoDesc = "MsgCockLadderScoreboardWinnersTemplateNoDesc"

	MsgCockRaceFooterActiveSeason = "MsgCockRaceFooterActiveSeason"
	MsgCockRaceFooterNoSeason     = "MsgCockRaceFooterNoSeason"

	MsgCockDynamicsTemplate = "MsgCockDynamicsTemplate"

	MsgCockSeasonTemplate            = "MsgCockSeasonTemplate"
	MsgCockSeasonWithWinnersTemplate = "MsgCockSeasonWithWinnersTemplate"
	MsgCockSeasonTemplateFooter      = "MsgCockSeasonTemplateFooter"
	MsgCockSeasonWinnerTemplate      = "MsgCockSeasonWinnerTemplate"
	MsgCockSeasonWinnerWithRespects  = "MsgCockSeasonWinnerWithRespects"
	MsgCockSeasonNoSeasonsTemplate   = "MsgCockSeasonNoSeasonsTemplate"

	MsgCockAchievementsTemplate           = "MsgCockAchievementsTemplate"
	MsgCockAchievementsTemplateOtherPages = "MsgCockAchievementsTemplateOtherPages"

	MsgSystemInfoTemplate = "MsgSystemInfoTemplate"

	InlineTitleCockSize         = "InlineTitleCockSize"
	InlineTitleCockRuler        = "InlineTitleCockRuler"
	InlineTitleCockLadder       = "InlineTitleCockLadder"
	InlineTitleCockRace         = "InlineTitleCockRace"
	InlineTitleCockDynamic      = "InlineTitleCockDynamic"
	InlineTitleCockSeason       = "InlineTitleCockSeason"
	InlineTitleCockAchievements = "InlineTitleCockAchievements"
	InlineTitleSystemInfo       = "InlineTitleSystemInfo"

	DescCockSize         = "DescCockSize"
	DescCockRuler        = "DescCockRuler"
	DescCockLadder       = "DescCockLadder"
	DescCockRace         = "DescCockRace"
	DescCockDynamic      = "DescCockDynamic"
	DescCockSeason       = "DescCockSeason"
	DescCockAchievements = "DescCockAchievements"
	DescSystemInfo       = "DescSystemInfo"

	MsgCockDynamicNoData      = "MsgCockDynamicNoData"
	MsgSeasonUnknownStartDate = "MsgSeasonUnknownStartDate"
	MsgSeasonButton           = "MsgSeasonButton"
	MsgSeasonNotFound         = "MsgSeasonNotFound"
	MsgCallbackInvalidFormat  = "MsgCallbackInvalidFormat"
	MsgCallbackParseError     = "MsgCallbackParseError"

	MsgUserPullingRecently = "MsgUserPullingRecently"
	MsgUserPullingSince    = "MsgUserPullingSince"
	MsgListSeparator       = "MsgListSeparator"
	MsgListSeparatorLast   = "MsgListSeparatorLast"

	UnitDay            = "UnitDay"
	UnitMonth          = "UnitMonth"
	UnitYear           = "UnitYear"
	UnitSeason         = "UnitSeason"
	UnitSeasonGenitive = "UnitSeasonGenitive"
	UnitHour           = "UnitHour"
	UnitMinute         = "UnitMinute"

	UptimeDayShort    = "UptimeDayShort"
	UptimeHourShort   = "UptimeHourShort"
	UptimeMinuteShort = "UptimeMinuteShort"

	MsgUnknownValue = "MsgUnknownValue"

	AnonymousNameTemplate = "AnonymousNameTemplate"

	MsgAchievementCompleted    = "MsgAchievementCompleted"
	MsgAchievementInProgress   = "MsgAchievementInProgress"
	MsgAchievementNotCompleted = "MsgAchievementNotCompleted"

	LuckLabelGodRandom      = "LuckLabelGodRandom"
	LuckLabelCosmicLuck     = "LuckLabelCosmicLuck"
	LuckLabelFairyLuck      = "LuckLabelFairyLuck"
	LuckLabelSuperLuck      = "LuckLabelSuperLuck"
	LuckLabelIncredibleLuck = "LuckLabelIncredibleLuck"
	LuckLabelVeryLucky      = "LuckLabelVeryLucky"
	LuckLabelLucky          = "LuckLabelLucky"
	LuckLabelBalanced       = "LuckLabelBalanced"
	LuckLabelUnlucky        = "LuckLabelUnlucky"
	LuckLabelBad            = "LuckLabelBad"
	LuckLabelGloom          = "LuckLabelGloom"
	LuckLabelHellTilt       = "LuckLabelHellTilt"
	LuckLabelBurningInHell  = "LuckLabelBurningInHell"

	VolatilityLabelStone        = "VolatilityLabelStone"
	VolatilityLabelStable       = "VolatilityLabelStable"
	VolatilityLabelModerate     = "VolatilityLabelModerate"
	VolatilityLabelLivelySpread = "VolatilityLabelLivelySpread"
	VolatilityLabelUneven       = "VolatilityLabelUneven"
	VolatilityLabelChaotic      = "VolatilityLabelChaotic"
	VolatilityLabelRandom       = "VolatilityLabelRandom"

	IrkLabelZero      = "IrkLabelZero"
	IrkLabelMinimal   = "IrkLabelMinimal"
	IrkLabelVerySmall = "IrkLabelVerySmall"
	IrkLabelSmall     = "IrkLabelSmall"
	IrkLabelReduced   = "IrkLabelReduced"
	IrkLabelAverage   = "IrkLabelAverage"
	IrkLabelIncreased = "IrkLabelIncreased"
	IrkLabelLarge     = "IrkLabelLarge"
	IrkLabelVeryLarge = "IrkLabelVeryLarge"
	IrkLabelMaximum   = "IrkLabelMaximum"
	IrkLabelUltimate  = "IrkLabelUltimate"

	GrowthSpeedLabelCosmic   = "GrowthSpeedLabelCosmic"
	GrowthSpeedLabelExtreme  = "GrowthSpeedLabelExtreme"
	GrowthSpeedLabelVeryFast = "GrowthSpeedLabelVeryFast"
	GrowthSpeedLabelFast     = "GrowthSpeedLabelFast"
	GrowthSpeedLabelModerate = "GrowthSpeedLabelModerate"
	GrowthSpeedLabelAverage  = "GrowthSpeedLabelAverage"
	GrowthSpeedLabelSlow     = "GrowthSpeedLabelSlow"
	GrowthSpeedLabelVerySlow = "GrowthSpeedLabelVerySlow"
	GrowthSpeedLabelTurtle   = "GrowthSpeedLabelTurtle"
	GrowthSpeedLabelStalled  = "GrowthSpeedLabelStalled"
)

func NewMsgCockDynamicsTemplate(
	localizationManager *localization.LocalizationManager,
	localizer *i18n.Localizer,
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

	/* Скорость роста общей статистики */

	overallGrowthSpeed float64,

	/* Период дергания кока пользователем */

	userPullingPeriod string,
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

	return localizationManager.Localize(localizer, MsgCockDynamicsTemplate, map[string]any{
		/* 1-2: Общая динамика коков */
		"TotalCock":  EscapeMarkdownV2(FormatDickSize(totalCock)),
		"TotalUsers": EscapeMarkdownV2(FormatDickSize(totalUsers)),

		/* 3-6: Средний и медианный кок */
		"TotalAvgCock":     EscapeMarkdownV2(FormatDickSize(totalAvgCock)),
		"TotalAvgEmoji":    EmojiFromSize(totalAvgCock),
		"TotalMedianCock":  EscapeMarkdownV2(FormatDickSize(totalMedianCock)),
		"TotalMedianEmoji": EmojiFromSize(totalMedianCock),

		/* 7-13: Персональная динамика кока */
		"UserTotalCock": EscapeMarkdownV2(FormatDickSize(userTotalCock)),
		"UserAvgCock":   EscapeMarkdownV2(FormatDickSize(userAvgCock)),
		"UserAvgEmoji":  EmojiFromSize(userAvgCock),
		"UserIrk":       EscapeMarkdownV2(FormatDickIkr(userIrk)),
		"UserMaxCock":   EscapeMarkdownV2(FormatDickSize(userMaxCock)),
		"UserMaxEmoji":  EmojiFromSize(userMaxCock),
		"UserMaxDate":   userMaxCockDate.In(datetime.NowLocation()).Format("02.01.06"),

		/* 14-18: Кок-активы (дневная и 5 коков динамика) */
		"UserYesterdayChangeEmoji":   userYesterdayChangePercentEmoji,
		"UserYesterdayChangePercent": fmt.Sprintf("%s%s", userYesterdayChangePercentSymbol, FormatDickPercent(userYesterdayChangePercent)),
		"UserYesterdayChangeSize":    fmt.Sprintf("%s%s", userYesterdayChangePercentSymbol, FormatDickSize(userYesterdayChangeCock)),
		"UserFiveCocksChangePercent": fmt.Sprintf("%s%s", userFiveCocksChangeSymbol, FormatDickPercent(userFiveCocksChangePercent)),
		"UserFiveCocksChangeSize":    fmt.Sprintf("%s%s", userFiveCocksChangeSymbol, FormatDickSize(userFiveCocksChangeCock)),

		/* 19-20: Соотношение коков */
		"TotalBigCockRatio":   FormatDickPercent(totalBigCockRatio),
		"TotalSmallCockRatio": FormatDickPercent(totalSmallCockRatio),

		/* 21-22: Самый большой кок */
		"TotalMaxCockDate": totalMaxCockDate.In(datetime.NowLocation()).Format("02.01.06"),
		"TotalMaxCock":     FormatDickSize(totalMaxCock),

		/* 23: % Доминирования */
		"UserDominancePercent": FormatDickPercent(userDominancePercent),

		/* 24-25: Сезонные достижения */
		"UserSeasonWins":  FormatDickSize(userSeasonWins),
		"UserCockRespect": FormatDickSize(userCockRespect),

		/* 26-27: Всего дёрнуто коков */
		"TotalCocksCount": EscapeMarkdownV2(FormatDickSize(totalCocksCount)),
		"UserCocksCount":  EscapeMarkdownV2(FormatDickSize(userCocksCount)),

		/* 28-31: Коэффициент везения и волатильность */
		"UserLuckCoefficient":   EscapeMarkdownV2(FormatLuckCoefficient(userLuckCoefficient)),
		"UserLuckDisplay":       LuckDisplay(localizationManager, localizer, userLuckCoefficient),
		"UserVolatility":        EscapeMarkdownV2(FormatVolatility(userVolatility)),
		"UserVolatilityDisplay": VolatilityDisplay(localizationManager, localizer, userVolatility),

		/* 32: Описание ИРК */
		"UserIrkLabel": IrkLabel(localizationManager, localizer, userIrk),

		/* 33: Эмодзи динамики за 5 коков */
		"UserFiveCocksEmoji": userFiveCocksChangeEmoji,

		/* 34-35: Скорость прироста кока */
		"UserGrowthSpeed":        EscapeMarkdownV2(FormatGrowthSpeed(userGrowthSpeed)),
		"UserGrowthSpeedDisplay": GrowthSpeedDisplay(localizationManager, localizer, userGrowthSpeed),

		/* 36: Скорость роста общей статистики */
		"OverallGrowthSpeed": EscapeMarkdownV2(FormatGrowthSpeed(overallGrowthSpeed)),

		/* 37: Период дергания кока */
		"UserPullingPeriod": userPullingPeriod,
	})
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

func NewMsgCockSeasonTemplate(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, pretenders string, startDate, endDate string, seasonNum int) string {
	return localizationManager.Localize(localizer, MsgCockSeasonTemplate, map[string]any{
		"Pretenders": pretenders,
		"StartDate":  startDate,
		"EndDate":    endDate,
		"SeasonNum":  seasonNum,
	})
}

func NewMsgCockSeasonWithWinnersTemplate(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, winners string, startDate, endDate string, seasonNum int) string {
	return localizationManager.Localize(localizer, MsgCockSeasonWithWinnersTemplate, map[string]any{
		"Winners":   winners,
		"StartDate": startDate,
		"EndDate":   endDate,
		"SeasonNum": seasonNum,
	})
}

func NewMsgCockSeasonWinnerTemplate(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, medal, nickname, totalSize string, respects int, showRespects bool) string {
	winnersLine := localizationManager.Localize(localizer, MsgCockSeasonWinnerTemplate, map[string]any{
		"Medal":    medal,
		"Nickname": EscapeMarkdownV2(nickname),
		"Size":     EscapeMarkdownV2(totalSize),
	})

	// Показываем респекты только если showRespects = true (для завершенных сезонов)
	if showRespects {
		formattedRespects := EscapeMarkdownV2(FormatDickSize(respects))
		return localizationManager.Localize(localizer, MsgCockSeasonWinnerWithRespects, map[string]any{
			"WinnerLine": winnersLine,
			"Respects":   formattedRespects,
		})
	}

	return winnersLine
}

func NewMsgCockSeasonTemplateFooter(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer) string {
	return localizationManager.Localize(localizer, MsgCockSeasonTemplateFooter, nil)
}

func NewMsgCockSeasonNoSeasonsTemplate(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer) string {
	return localizationManager.Localize(localizer, MsgCockSeasonNoSeasonsTemplate, nil)
}

func NewMsgSystemInfoTemplate(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, info *SystemInfo) string {
	return localizationManager.Localize(localizer, MsgSystemInfoTemplate, map[string]any{
		// Сервис
		"Uptime":   EscapeMarkdownV2(info.Uptime),
		"Version":  EscapeMarkdownV2(info.Version),
		"BuildRev": EscapeMarkdownV2(info.BuildRev),
		"BuildAt":  EscapeMarkdownV2(info.BuildAt),
		// Окружение
		"OS":            EscapeMarkdownV2(info.OS),
		"Arch":          EscapeMarkdownV2(info.Arch),
		"GoVersion":     EscapeMarkdownV2(info.GoVersion),
		"MemoryUsed":    info.MemoryUsed,
		"MemoryTotal":   info.MemoryTotal,
		"MemoryPercent": EscapeMarkdownV2(fmt.Sprintf("%.1f", info.MemoryPercent)),
		// Базы данных
		"MongoVersion": EscapeMarkdownV2(info.MongoVersion),
		"RedisVersion": EscapeMarkdownV2(info.RedisVersion),
		// Запрос
		"Username": EscapeMarkdownV2(info.Username),
		"UserID":   info.UserID,
		"BotID":    info.BotID,
	})
}

// NewMsgCockSeasonSinglePage генерирует текст для одной страницы сезона (постраничная навигация)
func NewMsgCockSeasonSinglePage(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, season CockSeason, getSeasonWinners func(CockSeason) []SeasonWinner, showDescription bool) string {
	startDate := EscapeMarkdownV2(season.StartDate.Format("02.01.2006"))
	endDate := EscapeMarkdownV2(season.EndDate.Format("02.01.2006"))

	winners := getSeasonWinners(season)
	var winnerLines []string

	for _, winner := range winners {
		medal := GetMedalByPosition(winner.Place - 1)
		normalizedNickname := NormalizeUsername(localizationManager, localizer, winner.Nickname, winner.UserID)
		respects := CalculateCockRespect(winner.Place)
		// Показываем респекты только для завершенных сезонов
		line := NewMsgCockSeasonWinnerTemplate(
			localizationManager,
			localizer,
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
		seasonBlock = NewMsgCockSeasonTemplate(localizationManager, localizer, winnersText, startDate, endDate, season.SeasonNum)
		// Футер показываем только для активного (текущего) сезона И если showDescription = true
		if showDescription {
			footer := NewMsgCockSeasonTemplateFooter(localizationManager, localizer)
			return seasonBlock + "\n\n" + footer
		}
		return seasonBlock
	} else {
		seasonBlock = NewMsgCockSeasonWithWinnersTemplate(localizationManager, localizer, winnersText, startDate, endDate, season.SeasonNum)
		return seasonBlock
	}
}

const (
	MsgCockAchievementsTitle = "MsgCockAchievementsTitle"
	MsgSystemInfoTitle       = "MsgSystemInfoTitle"
)
