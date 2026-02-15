package formatting

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

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

func FormatGrowthSpeed(speed float64) string {
	p := message.NewPrinter(language.Russian)
	return p.Sprintf("%.1f", speed)
}
