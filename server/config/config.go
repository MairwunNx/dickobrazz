package config

import (
	"errors"
	"os"
	"strings"
)

const defaultAddress = ":8090"

type Config struct {
	Address        string
	InternalToken  string
	MongoURI       string
	MongoDatabase  string
	RedisAddress   string
	RedisPassword  string
	TelegramToken  string
	RandomOrgToken string
	CorsOrigins    []string
}

func Load() (Config, error) {
	cfg := Config{
		Address:        valueOrDefault("SERVER_ADDRESS", defaultAddress),
		InternalToken:  os.Getenv("SERVER_INTERNAL_TOKEN"),
		MongoURI:       os.Getenv("MONGODB_URI"),
		MongoDatabase:  os.Getenv("MONGO_INITDB_DATABASE"),
		RedisAddress:   os.Getenv("REDIS_ADDRESS"),
		RedisPassword:  os.Getenv("REDIS_PASSWORD"),
		TelegramToken:  os.Getenv("TELEGRAM_TOKEN"),
		RandomOrgToken: os.Getenv("RANDOMORG_TOKEN"),
		CorsOrigins:    parseCSVEnv("SERVER_CORS_ORIGINS"),
	}

	if cfg.MongoURI == "" {
		return Config{}, errors.New("MONGODB_URI is required")
	}
	if cfg.MongoDatabase == "" {
		return Config{}, errors.New("MONGO_INITDB_DATABASE is required")
	}
	if cfg.RedisAddress == "" {
		return Config{}, errors.New("REDIS_ADDRESS is required")
	}
	if cfg.RedisPassword == "" {
		return Config{}, errors.New("REDIS_PASSWORD is required")
	}
	if cfg.TelegramToken == "" {
		return Config{}, errors.New("TELEGRAM_TOKEN is required")
	}
	if cfg.RandomOrgToken == "" {
		return Config{}, errors.New("RANDOMORG_TOKEN is required")
	}

	return cfg, nil
}

func valueOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func parseCSVEnv(key string) []string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	for idx, value := range parts {
		parts[idx] = strings.TrimSpace(value)
	}
	return parts
}
