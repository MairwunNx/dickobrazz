package emoji

import (
	"testing"
	"time"
)

func TestOrthodoxEaster(t *testing.T) {
	tests := []struct {
		year     int
		expected time.Time
	}{
		{2024, time.Date(2024, 5, 5, 0, 0, 0, 0, time.UTC)},
		{2025, time.Date(2025, 4, 20, 0, 0, 0, 0, time.UTC)},
		{2026, time.Date(2026, 4, 12, 0, 0, 0, 0, time.UTC)},
		{2027, time.Date(2027, 5, 2, 0, 0, 0, 0, time.UTC)},
	}

	for _, tt := range tests {
		got := OrthodoxEaster(tt.year, time.UTC)
		if !got.Equal(tt.expected) {
			t.Errorf("OrthodoxEaster(%d) = %v, want %v", tt.year, got, tt.expected)
		}
	}
}
