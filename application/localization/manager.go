package localization

import (
	"embed"
	"fmt"
	"os"
	"strings"

	"dickobrazz/application/logging"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pelletier/go-toml/v2"
	"golang.org/x/text/language"
)

//go:embed locales/*.toml
var localesFS embed.FS

const (
	defaultLanguage       = "en"
	supportedLanguagesEnv = "SUPPORTED_LANGUAGES"
)

type LocalizationManager struct {
	bundle *i18n.Bundle
	log    *logging.Logger
}

func NewLocalizationManager(log *logging.Logger) (*LocalizationManager, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	supportedLanguages := loadSupportedLanguages(log)
	for _, lang := range supportedLanguages {
		filename := fmt.Sprintf("locales/active.%s.toml", lang)

		data, err := localesFS.ReadFile(filename)
		if err != nil {
			log.E("Failed to read locale file", "filename", filename, logging.InnerError, err)
			return nil, fmt.Errorf("failed to read locale file %s: %w", filename, err)
		}

		if _, err := bundle.ParseMessageFileBytes(data, filename); err != nil {
			log.E("Failed to parse locale file", "filename", filename, logging.InnerError, err)
			return nil, fmt.Errorf("failed to parse locale file %s: %w", filename, err)
		}

		log.I("Loaded locale file", "filename", filename)
	}

	log.I("LocalizationManager initialized successfully")
	return &LocalizationManager{bundle: bundle, log: log}, nil
}

func loadSupportedLanguages(log *logging.Logger) []string {
	raw := strings.TrimSpace(os.Getenv(supportedLanguagesEnv))
	if raw == "" {
		return []string{"en", "ru", "es", "fr", "de", "zh"}
	}

	parts := strings.Split(raw, ",")
	languages := make([]string, 0, len(parts))
	for _, part := range parts {
		lang := strings.TrimSpace(part)
		if lang == "" {
			continue
		}
		languages = append(languages, lang)
	}

	if len(languages) == 0 {
		log.W("SUPPORTED_LANGUAGES is empty, using defaults")
		return []string{"en", "ru", "es", "fr", "de", "zh"}
	}

	return languages
}

func (x *LocalizationManager) LocalizerByUpdate(update *tgbotapi.Update) (*i18n.Localizer, string) {
	detectedLang := x.detectLanguage(update)
	return i18n.NewLocalizer(x.bundle, detectedLang, defaultLanguage), detectedLang
}

func (x *LocalizationManager) LocalizeByUpdate(update *tgbotapi.Update, messageID string) string {
	return x.LocalizeByUpdateTd(update, messageID, nil)
}

func (x *LocalizationManager) LocalizeByUpdateTd(update *tgbotapi.Update, messageID string, templateData map[string]any) string {
	localizer, _ := x.LocalizerByUpdate(update)
	return x.Localize(localizer, messageID, templateData)
}

func (x *LocalizationManager) Localize(localizer *i18n.Localizer, messageID string, templateData map[string]any) string {
	var pluralCount any
	if templateData != nil {
		if count, ok := templateData["Count"]; ok {
			pluralCount = count
		}
	}

	config := &i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
		PluralCount:  pluralCount,
	}

	msg, err := localizer.Localize(config)
	if err != nil {
		x.log.E("Failed to localize message", "message_id", messageID, logging.InnerError, err)
		return messageID
	}

	return msg
}

func (x *LocalizationManager) detectLanguage(update *tgbotapi.Update) string {
	if update == nil {
		return defaultLanguage
	}

	user := update.SentFrom()
	if user == nil || user.LanguageCode == "" {
		return defaultLanguage
	}

	return x.mapTelegramLanguageCode(user.LanguageCode)
}

func (x *LocalizationManager) mapTelegramLanguageCode(telegramCode string) string {
	lowerCode := strings.ToLower(telegramCode)

	switch {
	case strings.HasPrefix(lowerCode, "ru"),
		strings.HasPrefix(lowerCode, "uk"),
		strings.HasPrefix(lowerCode, "be"),
		strings.HasPrefix(lowerCode, "lv"),
		strings.HasPrefix(lowerCode, "lt"):
		return "ru"
	case strings.HasPrefix(lowerCode, "es"):
		return "es"
	case strings.HasPrefix(lowerCode, "fr"):
		return "fr"
	case strings.HasPrefix(lowerCode, "de"):
		return "de"
	case strings.HasPrefix(lowerCode, "zh"):
		return "zh"
	case strings.HasPrefix(lowerCode, "en"):
		return "en"
	default:
		return "en"
	}
}
