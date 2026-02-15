package emoji

import (
	"dickobrazz/src/shared/datetime"
	"time"
)

func EmojiFromSize(size int) string {
	holidaySet := GetHolidayEmojiSet()
	for _, r := range holidaySet {
		if size >= r.Min && size <= r.Max {
			return r.Emoji
		}
	}

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
	return "💎"
}

func GetHolidayEmojiSet() []EmojiRange {
	now := datetime.NowTime()
	month := now.Month()
	day := now.Day()
	year := now.Year()

	if (month == time.December && day == 31) || (month == time.January && day <= 10) {
		return NewYearEmojiRanges
	}
	if month == time.February && day == 14 {
		return ValentineEmojiRanges
	}
	if month == time.February && day == 23 {
		return DefenderDayEmojiRanges
	}
	if month == time.March && day == 8 {
		return WomenDayEmojiRanges
	}
	if month == time.April && day == 1 {
		return AprilFoolsEmojiRanges
	}
	if month == time.April && day == 12 {
		return CosmonauticsDayEmojiRanges
	}

	easterDate := OrthodoxEaster(year, datetime.NowLocation())
	easterEnd := easterDate.AddDate(0, 0, 7)
	if (now.After(easterDate) || now.Equal(easterDate)) && now.Before(easterEnd) {
		return EasterEmojiRanges
	}

	if month == time.June && day == 12 {
		return RussiaDayEmojiRanges
	}
	if month == time.October && day == 31 {
		return HalloweenEmojiRanges
	}

	return nil
}

func GetCurrentSeason() Season {
	now := datetime.NowTime()
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
	default:
		return Summer
	}
}
