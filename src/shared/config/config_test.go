package config

import (
	"os"
	"testing"
)

func TestExpandConfigurationEnvironment(t *testing.T) {
	os.Setenv("TEST_VAR", "hello")
	os.Setenv("EMPTY_VAR", "")
	defer os.Unsetenv("TEST_VAR")
	defer os.Unsetenv("EMPTY_VAR")

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple var", "${TEST_VAR}", "hello"},
		{"missing var", "${MISSING_VAR}", ""},
		{"default with colon-dash", "${MISSING_VAR:-fallback}", "fallback"},
		{"existing with colon-dash", "${TEST_VAR:-fallback}", "hello"},
		{"empty with colon-dash uses fallback", "${EMPTY_VAR:-fallback}", "fallback"},
		{"default with dash", "${MISSING_VAR-fallback}", "fallback"},
		{"empty with dash keeps empty", "${EMPTY_VAR-fallback}", ""},
		{"no substitution", "plain text", "plain text"},
		{"mixed", "url: ${TEST_VAR}:8080", "url: hello:8080"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandConfigurationEnvironment(tt.input)
			if got != tt.expected {
				t.Errorf("expandConfigurationEnvironment(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestValidateConfiguration(t *testing.T) {
	valid := &Configuration{
		Bot: BotConfiguration{
			CSOT:   "5s",
			Server: ServerConfiguration{BaseURL: "http://localhost"},
			Tg:     TelegramConfiguration{Token: "token123", Env: "dev"},
		},
	}

	if err := validateConfiguration(valid); err != nil {
		t.Errorf("validateConfiguration(valid) = %v, want nil", err)
	}

	tests := []struct {
		name string
		cfg  Configuration
	}{
		{"missing csot", Configuration{Bot: BotConfiguration{CSOT: "", Server: ServerConfiguration{BaseURL: "http://localhost"}, Tg: TelegramConfiguration{Token: "t", Env: "d"}}}},
		{"missing base_url", Configuration{Bot: BotConfiguration{CSOT: "5s", Server: ServerConfiguration{BaseURL: ""}, Tg: TelegramConfiguration{Token: "t", Env: "d"}}}},
		{"missing token", Configuration{Bot: BotConfiguration{CSOT: "5s", Server: ServerConfiguration{BaseURL: "http://localhost"}, Tg: TelegramConfiguration{Token: "", Env: "d"}}}},
		{"missing env", Configuration{Bot: BotConfiguration{CSOT: "5s", Server: ServerConfiguration{BaseURL: "http://localhost"}, Tg: TelegramConfiguration{Token: "t", Env: ""}}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateConfiguration(&tt.cfg); err == nil {
				t.Error("validateConfiguration() = nil, want error")
			}
		})
	}
}
