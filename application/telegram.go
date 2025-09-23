package application

import (
	"dickobrazz/application/logging"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

	log.I("Successfully connected to Telegram Bot API!")
	return bot
}
