package size

import (
	"dickobrazz/src/shared/formatting"
	"dickobrazz/src/shared/localization"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

func GenerateCockSizeText(localizationManager *localization.LocalizationManager, localizer *i18n.Localizer, size int, emoji string) string {
	formattedSize := formatting.FormatCockSizeForDate(size)
	return localizationManager.Localize(localizer, localization.MsgCockSize, map[string]any{
		"Size":  formattedSize,
		"Emoji": emoji,
	})
}
