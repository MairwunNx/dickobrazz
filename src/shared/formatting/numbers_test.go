package formatting

import "testing"

func TestFormatDickSize(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{5, "5"},
		{42, "42"},
		{1000, "1\u00a0000"},
		{1234567, "1\u00a0234\u00a0567"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := FormatDickSize(tt.input)
			if got != tt.expected {
				t.Errorf("FormatDickSize(%d) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFormatDickPercent(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0.0, "0,0"},
		{50.5, "50,5"},
		{99.9, "99,9"},
		{100.12, "100,1"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := FormatDickPercent(tt.input)
			if got != tt.expected {
				t.Errorf("FormatDickPercent(%f) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFormatDickIkr(t *testing.T) {
	got := FormatDickIkr(1.234)
	if got != "1,234" {
		t.Errorf("FormatDickIkr(1.234) = %q, want %q", got, "1,234")
	}
}

func TestFormatLuckCoefficient(t *testing.T) {
	got := FormatLuckCoefficient(0.567)
	if got != "0,567" {
		t.Errorf("FormatLuckCoefficient(0.567) = %q, want %q", got, "0,567")
	}
}

func TestFormatVolatility(t *testing.T) {
	got := FormatVolatility(12.3)
	if got != "12,3" {
		t.Errorf("FormatVolatility(12.3) = %q, want %q", got, "12,3")
	}
}

func TestFormatGrowthSpeed(t *testing.T) {
	got := FormatGrowthSpeed(5.7)
	if got != "5,7" {
		t.Errorf("FormatGrowthSpeed(5.7) = %q, want %q", got, "5,7")
	}
}
