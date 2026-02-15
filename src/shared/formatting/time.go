package formatting

import (
	"dickobrazz/src/shared/localization"
	"fmt"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func FormatTimeRemaining(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, endDate time.Time, now time.Time) string {
	duration := endDate.Sub(now)
	daysRemaining := int(duration.Hours() / 24)

	if daysRemaining < 0 {
		return localizationManager.Localize(localizer, localization.UnitDay, map[string]any{"Count": 0})
	}

	if daysRemaining > 30 {
		months := daysRemaining / 30
		days := daysRemaining % 30

		if days == 0 {
			return localizationManager.Localize(localizer, localization.UnitMonth, map[string]any{"Count": months})
		}
		monthsText := localizationManager.Localize(localizer, localization.UnitMonth, map[string]any{"Count": months})
		daysText := localizationManager.Localize(localizer, localization.UnitDay, map[string]any{"Count": days})
		return fmt.Sprintf("%s %s", monthsText, daysText)
	}

	return localizationManager.Localize(localizer, localization.UnitDay, map[string]any{"Count": daysRemaining})
}

func FormatUserPullingPeriod(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, firstCockDate time.Time, now time.Time) string {
	years := now.Year() - firstCockDate.Year()
	months := int(now.Month()) - int(firstCockDate.Month())
	days := now.Day() - firstCockDate.Day()

	if days < 0 {
		months--
		prevMonth := now.AddDate(0, -1, 0)
		daysInPrevMonth := time.Date(prevMonth.Year(), prevMonth.Month()+1, 0, 0, 0, 0, 0, prevMonth.Location()).Day()
		days += daysInPrevMonth
	}

	if months < 0 {
		years--
		months += 12
	}

	dateStr := firstCockDate.Format("02.01.2006")

	var parts []string
	if years > 0 {
		parts = append(parts, localizationManager.Localize(localizer, localization.UnitYear, map[string]any{"Count": years}))
	}
	if months > 0 {
		parts = append(parts, localizationManager.Localize(localizer, localization.UnitMonth, map[string]any{"Count": months}))
	}
	if days > 0 || len(parts) == 0 {
		parts = append(parts, localizationManager.Localize(localizer, localization.UnitDay, map[string]any{"Count": days}))
	}

	var result string
	if len(parts) == 1 {
		result = parts[0]
	} else if len(parts) == 2 {
		result = parts[0] + localizationManager.Localize(localizer, localization.MsgListSeparatorLast, nil) + parts[1]
	} else if len(parts) == 3 {
		result = parts[0] + localizationManager.Localize(localizer, localization.MsgListSeparator, nil) + parts[1] + localizationManager.Localize(localizer, localization.MsgListSeparatorLast, nil) + parts[2]
	}

	return localizationManager.Localize(localizer, localization.MsgUserPullingSince, map[string]any{
		"Period": result,
		"Date":   dateStr,
	})
}
