package application

import (
	"dickobrazz/application/logging"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	DefaultAPIEndpoint = "https://api.telegram.org/bot%s/%s"
)

func InitializeTelegramBot(log *logging.Logger, cfg *Configuration) *tgbotapi.BotAPI {
	token := cfg.Bot.Tg.Token
	if token == "" {
		log.F("Telegram token must be set and non-empty")
	}

	telegramEnv := cfg.Bot.Tg.Env
	if telegramEnv == "" {
		telegramEnv = "production"
	}

	var apiEndpoint string
	switch telegramEnv {
	case "test":
		apiEndpoint = "https://api.telegram.org/bot%s/test/%s"
		log.I("Using Telegram TEST API endpoint")
	case "production":
		apiEndpoint = DefaultAPIEndpoint
		log.I("Using Telegram PRODUCTION API endpoint")
	default:
		log.F("Invalid TELEGRAM_ENV value", "telegram_env", telegramEnv)
	}

	bot, err := tgbotapi.NewBotAPIWithAPIEndpoint(token, apiEndpoint)
	if err != nil {
		log.F("Failed to initialize Telegram Bot API", logging.InnerError, err)
	}

	log.I("Successfully connected to Telegram Bot API!", "environment", telegramEnv)
	return bot
}
