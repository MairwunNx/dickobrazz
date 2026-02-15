package formatting

import (
	"dickobrazz/src/shared/localization"
	"fmt"
	"math"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

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
		return localizationManager.Localize(localizer, localization.LuckLabelGodRandom, nil)
	case luck >= 1.92:
		return localizationManager.Localize(localizer, localization.LuckLabelCosmicLuck, nil)
	case luck >= 1.833:
		return localizationManager.Localize(localizer, localization.LuckLabelFairyLuck, nil)
	case luck >= 1.7:
		return localizationManager.Localize(localizer, localization.LuckLabelSuperLuck, nil)
	case luck >= 1.5:
		return localizationManager.Localize(localizer, localization.LuckLabelIncredibleLuck, nil)
	case luck >= 1.2:
		return localizationManager.Localize(localizer, localization.LuckLabelVeryLucky, nil)
	case luck >= 1.1:
		return localizationManager.Localize(localizer, localization.LuckLabelLucky, nil)
	case luck >= 0.9:
		return localizationManager.Localize(localizer, localization.LuckLabelBalanced, nil)
	case luck >= 0.7:
		return localizationManager.Localize(localizer, localization.LuckLabelUnlucky, nil)
	case luck >= 0.5:
		return localizationManager.Localize(localizer, localization.LuckLabelBad, nil)
	case luck >= 0.3:
		return localizationManager.Localize(localizer, localization.LuckLabelGloom, nil)
	case luck >= 0.2:
		return localizationManager.Localize(localizer, localization.LuckLabelHellTilt, nil)
	default:
		return localizationManager.Localize(localizer, localization.LuckLabelBurningInHell, nil)
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
		return localizationManager.Localize(localizer, localization.VolatilityLabelStone, nil)
	case volatility < 3:
		return localizationManager.Localize(localizer, localization.VolatilityLabelStable, nil)
	case volatility < 6:
		return localizationManager.Localize(localizer, localization.VolatilityLabelModerate, nil)
	case volatility < 10:
		return localizationManager.Localize(localizer, localization.VolatilityLabelLivelySpread, nil)
	case volatility < 15:
		return localizationManager.Localize(localizer, localization.VolatilityLabelUneven, nil)
	case volatility < 25:
		return localizationManager.Localize(localizer, localization.VolatilityLabelChaotic, nil)
	default:
		return localizationManager.Localize(localizer, localization.VolatilityLabelRandom, nil)
	}
}

func VolatilityDisplay(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, volatility float64) string {
	return fmt.Sprintf("%s _(%s)_", VolatilityEmoji(volatility), VolatilityLabel(localizationManager, localizer, volatility))
}

func IrkLabel(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, irk float64) string {
	irk = clamp01(irk)

	bucket := int(math.Floor(irk * 10))
	if irk >= 1.0 {
		bucket = 10
	}

	labels := [...]string{
		localization.IrkLabelZero, localization.IrkLabelMinimal, localization.IrkLabelVerySmall,
		localization.IrkLabelSmall, localization.IrkLabelReduced, localization.IrkLabelAverage,
		localization.IrkLabelIncreased, localization.IrkLabelLarge, localization.IrkLabelVeryLarge,
		localization.IrkLabelMaximum, localization.IrkLabelUltimate,
	}

	return localizationManager.Localize(localizer, labels[bucket], nil)
}

func GrowthSpeedLabel(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, speed float64) string {
	switch {
	case speed >= 50:
		return localizationManager.Localize(localizer, localization.GrowthSpeedLabelCosmic, nil)
	case speed >= 40:
		return localizationManager.Localize(localizer, localization.GrowthSpeedLabelExtreme, nil)
	case speed >= 30:
		return localizationManager.Localize(localizer, localization.GrowthSpeedLabelVeryFast, nil)
	case speed >= 20:
		return localizationManager.Localize(localizer, localization.GrowthSpeedLabelFast, nil)
	case speed >= 15:
		return localizationManager.Localize(localizer, localization.GrowthSpeedLabelModerate, nil)
	case speed >= 10:
		return localizationManager.Localize(localizer, localization.GrowthSpeedLabelAverage, nil)
	case speed >= 5:
		return localizationManager.Localize(localizer, localization.GrowthSpeedLabelSlow, nil)
	case speed >= 2:
		return localizationManager.Localize(localizer, localization.GrowthSpeedLabelVerySlow, nil)
	case speed >= 0.5:
		return localizationManager.Localize(localizer, localization.GrowthSpeedLabelTurtle, nil)
	default:
		return localizationManager.Localize(localizer, localization.GrowthSpeedLabelStalled, nil)
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
