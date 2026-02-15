package config

import (
	"dickobrazz/src/shared/logging"
	"sync"
)

type Configuration struct {
	Bot BotConfiguration `yaml:"bot"`
}

type BotConfiguration struct {
	CSOT   string                `yaml:"csot"`
	Server ServerConfiguration   `yaml:"server"`
	Tg     TelegramConfiguration `yaml:"tg"`
}

type ServerConfiguration struct {
	BaseURL string `yaml:"base_url"`
}

type TelegramConfiguration struct {
	Token string `yaml:"token"`
	Env   string `yaml:"env"`
}

var (
	loadConfigurationOnce  sync.Once
	cachedConfiguration    *Configuration
	cachedConfigurationErr error
)

func NewConfiguration(log *logging.Logger) *Configuration {
	loadConfigurationOnce.Do(func() {
		cachedConfiguration, cachedConfigurationErr = loadConfiguration()
	})

	if cachedConfigurationErr != nil {
		log.F("Failed to load configuration", logging.InnerError, cachedConfigurationErr)
	}

	return cachedConfiguration
}
