package config

import (
	"os"
	"regexp"
)

var envExpression = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)(?:(:-|-)([^}]*))?\}`)

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
