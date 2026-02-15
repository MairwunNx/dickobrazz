package formatting

import "testing"

func TestGetPlaceEmoji_TopThree(t *testing.T) {
	tests := []struct {
		place    int
		expected string
	}{
		{1, "🥇"},
		{2, "🥈"},
		{3, "🥉"},
	}

	for _, tt := range tests {
		got := GetPlaceEmoji(tt.place, false)
		if got != tt.expected {
			t.Errorf("GetPlaceEmoji(%d, false) = %q, want %q", tt.place, got, tt.expected)
		}
		gotCurrent := GetPlaceEmoji(tt.place, true)
		if gotCurrent != tt.expected {
			t.Errorf("GetPlaceEmoji(%d, true) = %q, want %q", tt.place, gotCurrent, tt.expected)
		}
	}
}

func TestGetMedalByPosition(t *testing.T) {
	tests := []struct {
		position int
		expected string
	}{
		{0, "🥇"},
		{1, "🥈"},
		{2, "🥉"},
		{3, ""},
		{10, ""},
		{-1, ""},
	}

	for _, tt := range tests {
		got := GetMedalByPosition(tt.position)
		if got != tt.expected {
			t.Errorf("GetMedalByPosition(%d) = %q, want %q", tt.position, got, tt.expected)
		}
	}
}

func TestGetPlaceEmojiForContext(t *testing.T) {
	got := GetPlaceEmojiForContext(5, false)
	if got != "🥀 5\\." {
		t.Errorf("GetPlaceEmojiForContext(5, false) = %q, want %q", got, "🥀 5\\.")
	}

	gotBold := GetPlaceEmojiForContext(5, true)
	if gotBold != "🥀 *5*\\." {
		t.Errorf("GetPlaceEmojiForContext(5, true) = %q, want %q", gotBold, "🥀 *5*\\.")
	}
}
