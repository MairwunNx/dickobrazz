package formatting

import (
	"testing"
	"time"
)

func TestIsMathDay(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected bool
	}{
		{"pi day", time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC), true},
		{"not pi day", time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC), false},
		{"different month", time.Date(2026, 4, 14, 0, 0, 0, 0, time.UTC), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isMathDay(tt.date)
			if got != tt.expected {
				t.Errorf("isMathDay(%v) = %v, want %v", tt.date, got, tt.expected)
			}
		})
	}
}

func TestIsProgrammersDay(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		expected bool
	}{
		{"programmers day 2026", time.Date(2026, 9, 13, 0, 0, 0, 0, time.UTC), true},
		{"not programmers day", time.Date(2026, 9, 14, 0, 0, 0, 0, time.UTC), false},
		{"programmers day 2024 leap", time.Date(2024, 9, 12, 0, 0, 0, 0, time.UTC), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isProgrammersDay(tt.date)
			if got != tt.expected {
				t.Errorf("isProgrammersDay(%v) = %v, want %v", tt.date, got, tt.expected)
			}
		})
	}
}

func TestFancyMathOrDefault(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "sin(0)"},
		{1, "0!"},
		{4, "2²"},
		{61, "√3721"},
		{100, "100"},
		{-5, "-5"},
	}

	for _, tt := range tests {
		got := fancyMathOrDefault(tt.input)
		if got != tt.expected {
			t.Errorf("fancyMathOrDefault(%d) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestToProgrammersNotation(t *testing.T) {
	got := toProgrammersNotation(10)
	if got != "0b1010" && got != "0xA" {
		t.Errorf("toProgrammersNotation(10) = %q, want binary or hex", got)
	}

	gotNeg := toProgrammersNotation(-5)
	if gotNeg != "-0b101" && gotNeg != "-0x5" {
		t.Errorf("toProgrammersNotation(-5) = %q, want negative binary or hex", gotNeg)
	}
}
