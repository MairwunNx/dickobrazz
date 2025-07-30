package application

import "time"

type EmojiRange struct {
	Min, Max int
	Emoji    string
}

var SummerEmojiRanges = []EmojiRange{
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
	{14, 16, "🥒"},
	{17, 20, "🥕"},
	{21, 28, "🔥"},
	{29, 34, "🍆"},
	{35, 42, "🌶️"},
	{43, 49, "🎅"},
	{50, 54, "🌲"},
	{55, 59, "🏔️"},
	{60, 60, "🚀"},
}

var AutumnEmojiRanges = []EmojiRange{
	{0, 1, "🍂"},
	{2, 3, "🍁"},
	{4, 8, "🌰"},
	{9, 13, "🍄"},
	{14, 16, "🥬"},
	{17, 20, "🎃"},
	{21, 28, "🔥"},
	{29, 34, "🍆"},
	{35, 42, "🌶️"},
	{43, 49, "👨🏻‍🌾"},
	{50, 54, "🍇"},
	{55, 59, "🌧️"},
	{60, 60, "🚀"},
}

var SpringEmojiRanges = []EmojiRange{
	{0, 1, "🌱"},
	{2, 3, "🌼"},
	{4, 8, "🌷"},
	{9, 13, "🌈"},
	{14, 16, "🥦"},
	{17, 20, "🌻"},
	{21, 28, "🔥"},
	{29, 34, "🍆"},
	{35, 42, "🌶️"},
	{43, 49, "👩‍🌾"},
	{50, 54, "🐝"},
	{55, 59, "🌺"},
	{60, 60, "🚀"},
}

type Season int

const (
	Winter Season = iota
	Spring
	Summer
	Autumn
)

func EmojiFromSize(size int) string {
	season := GetCurrentSeason()
	
	var emojiSet []EmojiRange
	
	switch season {
	case Winter:
		emojiSet = WinterEmojiRanges
	case Spring:
		emojiSet = SpringEmojiRanges
	case Autumn:
		emojiSet = AutumnEmojiRanges
	case Summer:
		emojiSet = SummerEmojiRanges
	}

	for _, r := range emojiSet {
		if size >= r.Min && size <= r.Max {
			return r.Emoji
		}
	}
	return "❓❓❓"
}

func GetCurrentSeason() Season {
	now := time.Now()
	month := now.Month()
	
	switch month {
	case time.December, time.January, time.February:
		return Winter
	case time.March, time.April, time.May:
		return Spring
	case time.June, time.July, time.August:
		return Summer
	case time.September, time.October, time.November:
		return Autumn
	default: // Компилятор "якобы" умнее меня, не догадывается, ведь сезонов то 4 *trollface.png*
		return Summer
	}
}
