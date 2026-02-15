package dynamics

import (
	"dickobrazz/src/shared/datetime"
	"dickobrazz/src/shared/emoji"
	"dickobrazz/src/shared/formatting"
	"dickobrazz/src/shared/localization"
	"fmt"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func NewMsgCockDynamicsTemplate(
	localizationManager *localization.LocalizationManager,
	localizer *i18n.Localizer,
	totalCock int, totalUsers int, totalAvgCock int, totalMedianCock int,
	userTotalCock int, userAvgCock int, userIrk float64,
	userMaxCock int, userMaxCockDate time.Time,
	userYesterdayChangePercent float64, userYesterdayChangeCock int,
	userFiveCocksChangePercent float64, userFiveCocksChangeCock int,
	totalBigCockRatio float64, totalSmallCockRatio float64,
	totalMaxCockDate time.Time, totalMaxCock int,
	userDominancePercent float64,
	userSeasonWins int, userCockRespect int,
	totalCocksCount int, userCocksCount int,
	userLuckCoefficient float64, userVolatility float64,
	userGrowthSpeed float64, overallGrowthSpeed float64,
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
		"TotalCock":  formatting.EscapeMarkdownV2(formatting.FormatDickSize(totalCock)),
		"TotalUsers": formatting.EscapeMarkdownV2(formatting.FormatDickSize(totalUsers)),

		"TotalAvgCock":     formatting.EscapeMarkdownV2(formatting.FormatDickSize(totalAvgCock)),
		"TotalAvgEmoji":    emoji.EmojiFromSize(totalAvgCock),
		"TotalMedianCock":  formatting.EscapeMarkdownV2(formatting.FormatDickSize(totalMedianCock)),
		"TotalMedianEmoji": emoji.EmojiFromSize(totalMedianCock),

		"UserTotalCock": formatting.EscapeMarkdownV2(formatting.FormatDickSize(userTotalCock)),
		"UserAvgCock":   formatting.EscapeMarkdownV2(formatting.FormatDickSize(userAvgCock)),
		"UserAvgEmoji":  emoji.EmojiFromSize(userAvgCock),
		"UserIrk":       formatting.EscapeMarkdownV2(formatting.FormatDickIkr(userIrk)),
		"UserMaxCock":   formatting.EscapeMarkdownV2(formatting.FormatDickSize(userMaxCock)),
		"UserMaxEmoji":  emoji.EmojiFromSize(userMaxCock),
		"UserMaxDate":   userMaxCockDate.In(datetime.NowLocation()).Format("02.01.06"),

		"UserYesterdayChangeEmoji":   userYesterdayChangePercentEmoji,
		"UserYesterdayChangePercent": fmt.Sprintf("%s%s", userYesterdayChangePercentSymbol, formatting.FormatDickPercent(userYesterdayChangePercent)),
		"UserYesterdayChangeSize":    fmt.Sprintf("%s%s", userYesterdayChangePercentSymbol, formatting.FormatDickSize(userYesterdayChangeCock)),
		"UserFiveCocksChangePercent": fmt.Sprintf("%s%s", userFiveCocksChangeSymbol, formatting.FormatDickPercent(userFiveCocksChangePercent)),
		"UserFiveCocksChangeSize":    fmt.Sprintf("%s%s", userFiveCocksChangeSymbol, formatting.FormatDickSize(userFiveCocksChangeCock)),

		"TotalBigCockRatio":   formatting.FormatDickPercent(totalBigCockRatio),
		"TotalSmallCockRatio": formatting.FormatDickPercent(totalSmallCockRatio),

		"TotalMaxCockDate": totalMaxCockDate.In(datetime.NowLocation()).Format("02.01.06"),
		"TotalMaxCock":     formatting.FormatDickSize(totalMaxCock),

		"UserDominancePercent": formatting.FormatDickPercent(userDominancePercent),

		"UserSeasonWins":  formatting.FormatDickSize(userSeasonWins),
		"UserCockRespect": formatting.FormatDickSize(userCockRespect),

		"TotalCocksCount": formatting.EscapeMarkdownV2(formatting.FormatDickSize(totalCocksCount)),
		"UserCocksCount":  formatting.EscapeMarkdownV2(formatting.FormatDickSize(userCocksCount)),

		"UserLuckCoefficient":   formatting.EscapeMarkdownV2(formatting.FormatLuckCoefficient(userLuckCoefficient)),
		"UserLuckDisplay":       formatting.LuckDisplay(localizationManager, localizer, userLuckCoefficient),
		"UserVolatility":        formatting.EscapeMarkdownV2(formatting.FormatVolatility(userVolatility)),
		"UserVolatilityDisplay": formatting.VolatilityDisplay(localizationManager, localizer, userVolatility),

		"UserIrkLabel": formatting.IrkLabel(localizationManager, localizer, userIrk),

		"UserFiveCocksEmoji": userFiveCocksChangeEmoji,

		"UserGrowthSpeed":        formatting.EscapeMarkdownV2(formatting.FormatGrowthSpeed(userGrowthSpeed)),
		"UserGrowthSpeedDisplay": formatting.GrowthSpeedDisplay(localizationManager, localizer, userGrowthSpeed),

		"OverallGrowthSpeed": formatting.EscapeMarkdownV2(formatting.FormatGrowthSpeed(overallGrowthSpeed)),

		"UserPullingPeriod": userPullingPeriod,
	})
}
