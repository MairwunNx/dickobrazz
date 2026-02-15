package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

	"dickobrazz/src/shared/logging"

	"gopkg.in/yaml.v3"
)

var configurationPaths = []string{
	"./config.yaml",
	"/etc/dickobrazz/config.yaml",
}

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

var envExpression = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)(?:(:-|-)([^}]*))?\}`)

func LoadConfiguration(log *logging.Logger) *Configuration {
	loadConfigurationOnce.Do(func() {
		cachedConfiguration, cachedConfigurationErr = loadConfiguration()
	})

	if cachedConfigurationErr != nil {
		log.F("Failed to load configuration", logging.InnerError, cachedConfigurationErr)
	}

	return cachedConfiguration
}

func loadConfiguration() (*Configuration, error) {
	var readErrors []string

	for _, path := range configurationPaths {
		raw, err := os.ReadFile(path)
		if err != nil {
			readErrors = append(readErrors, fmt.Sprintf("%s: %v", path, err))
			continue
		}

		expanded := expandConfigurationEnvironment(string(raw))

		var cfg Configuration
		if err := yaml.Unmarshal([]byte(expanded), &cfg); err != nil {
			return nil, fmt.Errorf("parse yaml from %q: %w", path, err)
		}

		if err := validateConfiguration(&cfg); err != nil {
			return nil, err
		}

		return &cfg, nil
	}

	return nil, fmt.Errorf("configuration file not found (%s)", strings.Join(readErrors, "; "))
}

func expandConfigurationEnvironment(content string) string {
	return envExpression.ReplaceAllStringFunc(content, func(match string) string {
		groups := envExpression.FindStringSubmatch(match)
		if len(groups) < 2 {
			return match
		}

		key := groups[1]
		operator := groups[2]
		fallback := groups[3]

		value, exists := os.LookupEnv(key)
		switch operator {
		case "":
			if !exists {
				return ""
			}
			return value
		case ":-":
			if !exists || value == "" {
				return fallback
			}
			return value
		case "-":
			if !exists {
				return fallback
			}
			return value
		default:
			return match
		}
	})
}

func validateConfiguration(cfg *Configuration) error {
	if strings.TrimSpace(cfg.Bot.CSOT) == "" {
		return fmt.Errorf("bot.csot must be set")
	}
	if strings.TrimSpace(cfg.Bot.Server.BaseURL) == "" {
		return fmt.Errorf("bot.server.base_url must be set")
	}
	if strings.TrimSpace(cfg.Bot.Tg.Token) == "" {
		return fmt.Errorf("bot.tg.token must be set")
	}
	if strings.TrimSpace(cfg.Bot.Tg.Env) == "" {
		return fmt.Errorf("bot.tg.env must be set")
	}

	return nil
}
