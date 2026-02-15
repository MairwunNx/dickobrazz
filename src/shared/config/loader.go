package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var configurationPaths = []string{
	"./config.yaml",
	"/etc/dickobrazz/config.yaml",
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
