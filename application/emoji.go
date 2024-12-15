package application

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
	{50, 54, "🥥🌴🥥"},
	{55, 59, "🏰🗼🏰"},
	{60, 60, "🚀🌌"},
}

func EmojiFromSize(size int) string {
	for _, r := range EmojiRanges {
		if size >= r.Min && size <= r.Max {
			return r.Emoji
		}
	}
	return "❓❓❓"
}
