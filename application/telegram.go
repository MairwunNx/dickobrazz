package application

import (
	"dickobrazz/application/logging"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	DefaultAPIEndpoint = "https://api.telegram.org/bot%s/%s"
)

func InitializeTelegramBot(log *logging.Logger) *tgbotapi.BotAPI {
	token, exist := os.LookupEnv("TELEGRAM_TOKEN")
	if !exist || token == "" {
		log.F("Telegram token must be set and non-empty")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.F("Failed to initialize Telegram Bot API", logging.InnerError, err)
	}

	telegramEnv := os.Getenv("TELEGRAM_ENV")
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

	bot.SetAPIEndpoint(apiEndpoint)

	log.I("Successfully connected to Telegram Bot API!", "environment", telegramEnv)
	return bot
}
