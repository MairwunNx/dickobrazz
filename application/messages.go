package application

import (
	"dickobrazz/application/datetime"
	"dickobrazz/application/localization"
	"fmt"
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

	MsgCockAchievementsTemplate           = "MsgCockAchievementsTemplate"
	MsgCockAchievementsTemplateOtherPages = "MsgCockAchievementsTemplateOtherPages"

	InlineTitleCockDynamic = "InlineTitleCockDynamic"
	InlineTitleCockSeason = "InlineTitleCockSeason"

	DescCockDynamic = "DescCockDynamic"
	DescCockSeason  = "DescCockSeason"
	MsgSeasonUnknownStartDate = "MsgSeasonUnknownStartDate"
	MsgSeasonButton           = "MsgSeasonButton"
	MsgSeasonNotFound         = "MsgSeasonNotFound"
	MsgCallbackInvalidFormat = "MsgCallbackInvalidFormat"
	MsgCallbackParseError    = "MsgCallbackParseError"

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

	MsgHelpText         = "MsgHelpText"
	MsgHidePrompt       = "MsgHidePrompt"
	MsgHideStatusHidden = "MsgHideStatusHidden"
	MsgHideButtonHide   = "MsgHideButtonHide"
	MsgHideButtonShow   = "MsgHideButtonShow"

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

	return localizationManager.Localize(localizer, "MsgCockDynamicsTemplate", map[string]any{
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
	return localizationManager.Localize(localizer, "MsgCockSeasonTemplate", map[string]any{
		"Pretenders": pretenders,
		"StartDate":  startDate,
		"EndDate":    endDate,
		"SeasonNum":  seasonNum,
	})
}

func NewMsgCockSeasonWithWinnersTemplate(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, winners string, startDate, endDate string, seasonNum int) string {
	return localizationManager.Localize(localizer, "MsgCockSeasonWithWinnersTemplate", map[string]any{
		"Winners":   winners,
		"StartDate": startDate,
		"EndDate":   endDate,
		"SeasonNum": seasonNum,
	})
}

func NewMsgCockSeasonWinnerTemplate(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, medal, nickname, totalSize string) string {
	return localizationManager.Localize(localizer, "MsgCockSeasonWinnerTemplate", map[string]any{
		"Medal":    medal,
		"Nickname": EscapeMarkdownV2(nickname),
		"Size":     EscapeMarkdownV2(totalSize),
	})
}

func NewMsgCockSeasonTemplateFooter(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer) string {
	return localizationManager.Localize(localizer, "MsgCockSeasonTemplateFooter", nil)
}

func NewMsgCockSeasonNoSeasonsTemplate(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer) string {
	return localizationManager.Localize(localizer, "MsgCockSeasonNoSeasonsTemplate", nil)
}
