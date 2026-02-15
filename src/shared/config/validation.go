package config

import (
	"fmt"
	"strings"
)

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
