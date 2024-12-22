package application

import "time"

type EmojiRange struct {
	Min, Max int
	Emoji    string
}

var EmojiRanges = []EmojiRange{
	{0, 1, "🔍"},
	{2, 3, "🤏🏻"},
	{4, 8, "🍒🌡"},
	{9, 13, "📉"},
	{14, 16, "🍌"},
	{17, 20, "🌽"},
	{21, 28, "🔥"},
	{29, 34, "🍆"},
	{35, 42, "🌶️"},
	{43, 49, "👨🏿‍🦱"},
	{50, 54, "🌴"},
	{55, 59, "🗼"},
	{60, 60, "🚀"},
}

var WinterEmojiRanges = []EmojiRange{
	{0, 1, "🥶"},
	{2, 3, "🧊"},
	{4, 8, "❄️"},
	{9, 13, "☃️"},
	{14, 16, "🍌"},
	{17, 20, "🌽"},
	{21, 28, "🔥"},
	{29, 34, "🍆"},
	{35, 42, "🌶️"},
	{43, 49, "🎅"},
	{50, 54, "🌲"},
	{55, 59, "🏔️"},
	{60, 60, "🚀"},
}

func EmojiFromSize(size int) string {
	isWinter := IsWinter() // Определяем, зима ли сейчас
	emojiSet := EmojiRanges
	if isWinter {
		emojiSet = WinterEmojiRanges
	}

	for _, r := range emojiSet {
		if size >= r.Min && size <= r.Max {
			return r.Emoji
		}
	}
	return "❓❓❓"
}

func IsWinter() bool {
	now := time.Now()
	month := now.Month()
	return month == time.December || month == time.January || month == time.February
}
