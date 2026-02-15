package formatting

import "testing"

func TestEscapeMarkdownV2(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"hello", "hello"},
		{"hello_world", "hello\\_world"},
		{"*bold*", "\\*bold\\*"},
		{"[link](url)", "\\[link\\]\\(url\\)"},
		{"~strikethrough~", "\\~strikethrough\\~"},
		{"`code`", "\\`code\\`"},
		{"#heading", "\\#heading"},
		{"a+b=c", "a\\+b\\=c"},
		{"|pipe|", "\\|pipe\\|"},
		{"{braces}", "\\{braces\\}"},
		{"end.", "end\\."},
		{"wow!", "wow\\!"},
		{"a>b", "a\\>b"},
		{"line1\nline2", "line1\nline2"},
		{"100% done", "100% done"},
		{"price: $10", "price: $10"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := EscapeMarkdownV2(tt.input)
			if got != tt.expected {
				t.Errorf("EscapeMarkdownV2(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
