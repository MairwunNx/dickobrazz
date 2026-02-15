package formatting

import (
	"math"
	"testing"
)

func TestClamp01(t *testing.T) {
	tests := []struct {
		name     string
		input    float64
		expected float64
	}{
		{"negative", -0.5, 0},
		{"zero", 0, 0},
		{"middle", 0.5, 0.5},
		{"one", 1.0, 1.0},
		{"above_one", 1.5, 1},
		{"NaN", math.NaN(), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := clamp01(tt.input)
			if got != tt.expected {
				t.Errorf("clamp01(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestLuckEmoji(t *testing.T) {
	tests := []struct {
		luck     float64
		expected string
	}{
		{2.0, "👑🌌🌈🦄🍀🤩"},
		{1.95, "🌌🌈🦄🍀🤩"},
		{1.85, "🌈🦄🍀🤩"},
		{1.7, "🍀🤩"},
		{1.5, "🤩"},
		{1.2, "🍀✨"},
		{1.1, "🍀"},
		{0.9, "⚖️"},
		{0.7, "😕"},
		{0.5, "😔"},
		{0.3, "💀"},
		{0.2, "☠️"},
		{0.1, "🔥☠️🔥"},
	}

	for _, tt := range tests {
		got := LuckEmoji(tt.luck)
		if got != tt.expected {
			t.Errorf("LuckEmoji(%v) = %q, want %q", tt.luck, got, tt.expected)
		}
	}
}

func TestVolatilityEmoji(t *testing.T) {
	tests := []struct {
		vol      float64
		expected string
	}{
		{0.5, "🧱"},
		{2, "🧊"},
		{5, "📈"},
		{8, "📉📈"},
		{12, "🎢"},
		{20, "🎢🌪️"},
		{30, "🌪️💥"},
	}

	for _, tt := range tests {
		got := VolatilityEmoji(tt.vol)
		if got != tt.expected {
			t.Errorf("VolatilityEmoji(%v) = %q, want %q", tt.vol, got, tt.expected)
		}
	}
}

func TestGrowthSpeedEmoji(t *testing.T) {
	tests := []struct {
		speed    float64
		expected string
	}{
		{55, "👑🌌🚀💫"},
		{45, "🚀🔥⚡"},
		{35, "⚡💨🏎️"},
		{25, "🏃💨"},
		{15, "🚶‍♂️⏱️"},
		{10, "🚶"},
		{5, "🐢⏳"},
		{2, "🐌🕰️"},
		{0.5, "🐢🌿"},
		{0.1, "🗿⛔"},
	}

	for _, tt := range tests {
		got := GrowthSpeedEmoji(tt.speed)
		if got != tt.expected {
			t.Errorf("GrowthSpeedEmoji(%v) = %q, want %q", tt.speed, got, tt.expected)
		}
	}
}
