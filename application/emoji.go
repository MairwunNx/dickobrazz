package application

import (
	"time"

	"dickobrazz/application/datetime"
)

type EmojiRange struct {
	Min, Max int
	Emoji    string
}

var NewYearEmojiRanges = []EmojiRange{
	{0, 1, "🎁"},
	{2, 3, "❄️"},
	{4, 8, "🎄"},
	{9, 13, "⛄"},
	{14, 16, "🥒"},
	{17, 20, "🥕"},
	{21, 28, "🔥"},
	{29, 34, "🍆"},
	{35, 42, "🌶️"},
	{43, 49, "🎅"},
	{50, 54, "⭐"},
	{55, 59, "🎆"},
	{60, 60, "🚀"},
}

var ValentineEmojiRanges = []EmojiRange{
	{0, 1, "💔"},
	{2, 3, "💌"},
	{4, 8, "💕"},
	{9, 13, "💖"},
	{14, 16, "🌹"},
	{17, 20, "💐"},
	{21, 28, "🔥"},
	{29, 34, "🍆"},
	{35, 42, "🌶️"},
	{43, 49, "💝"},
	{50, 54, "💗"},
	{55, 59, "💘"},
	{60, 60, "🚀"},
}

var DefenderDayEmojiRanges = []EmojiRange{
	{0, 1, "🎖️"},
	{2, 3, "🪖"},
	{4, 8, "🛡️"},
	{9, 13, "⚔️"},
	{14, 16, "🥒"},
	{17, 20, "🥕"},
	{21, 28, "🔥"},
	{29, 34, "🍆"},
	{35, 42, "🌶️"},
	{43, 49, "🎖️"},
	{50, 54, "💪"},
	{55, 59, "🦅"},
	{60, 60, "🚀"},
}

var WomenDayEmojiRanges = []EmojiRange{
	{0, 1, "🌸"},
	{2, 3, "🌼"},
	{4, 8, "🌷"},
	{9, 13, "💐"},
	{14, 16, "🥦"},
	{17, 20, "🌻"},
	{21, 28, "🔥"},
	{29, 34, "🍆"},
	{35, 42, "🌶️"},
	{43, 49, "👩"},
	{50, 54, "🌺"},
	{55, 59, "🌹"},
	{60, 60, "🚀"},
}

var AprilFoolsEmojiRanges = []EmojiRange{
	{0, 1, "🤪"},
	{2, 3, "🤡"},
	{4, 8, "😜"},
	{9, 13, "🎭"},
	{14, 16, "🥦"},
	{17, 20, "🌻"},
	{21, 28, "🔥"},
	{29, 34, "🍆"},
	{35, 42, "🌶️"},
	{43, 49, "🃏"},
	{50, 54, "🎪"},
	{55, 59, "🎉"},
	{60, 60, "🚀"},
}

var CosmonauticsDayEmojiRanges = []EmojiRange{
	{0, 1, "⭐"},
	{2, 3, "🌟"},
	{4, 8, "✨"},
	{9, 13, "🌠"},
	{14, 16, "🥦"},
	{17, 20, "🌻"},
	{21, 28, "🔥"},
	{29, 34, "🍆"},
	{35, 42, "🌶️"},
	{43, 49, "👨‍🚀"},
	{50, 54, "🛸"},
	{55, 59, "🪐"},
	{60, 60, "🚀"},
}

var EasterEmojiRanges = []EmojiRange{
	{0, 1, "🐣"},
	{2, 3, "🥚"},
	{4, 8, "🐰"},
	{9, 13, "🌷"},
	{14, 16, "🥦"},
	{17, 20, "🌻"},
	{21, 28, "🔥"},
	{29, 34, "🍆"},
	{35, 42, "🌶️"},
	{43, 49, "🐇"},
	{50, 54, "🪺"},
	{55, 59, "🌸"},
	{60, 60, "🚀"},
}

var RussiaDayEmojiRanges = []EmojiRange{
	{0, 1, "🇷🇺"},
	{2, 3, "🎉"},
	{4, 8, "🎊"},
	{9, 13, "🎈"},
	{14, 16, "🍌"},
	{17, 20, "🌽"},
	{21, 28, "🔥"},
	{29, 34, "🍆"},
	{35, 42, "🌶️"},
	{43, 49, "🇷🇺"},
	{50, 54, "🦅"},
	{55, 59, "⭐"},
	{60, 60, "🚀"},
}

var HalloweenEmojiRanges = []EmojiRange{
	{0, 1, "🦇"},
	{2, 3, "👻"},
	{4, 8, "🕷️"},
	{9, 13, "🕸️"},
	{14, 16, "🥬"},
	{17, 20, "🎃"},
	{21, 28, "🔥"},
	{29, 34, "🍆"},
	{35, 42, "🌶️"},
	{43, 49, "🧛"},
	{50, 54, "💀"},
	{55, 59, "🧟"},
	{60, 60, "🚀"},
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

	// Новый год: 31 декабря - 10 января
	if (month == time.December && day == 31) || (month == time.January && day <= 10) {
		return NewYearEmojiRanges
	}

	// День Святого Валентина: 14 февраля
	if month == time.February && day == 14 {
		return ValentineEmojiRanges
	}

	// День защитника Отечества: 23 февраля
	if month == time.February && day == 23 {
		return DefenderDayEmojiRanges
	}

	// 8 Марта
	if month == time.March && day == 8 {
		return WomenDayEmojiRanges
	}

	// День смеха: 1 апреля
	if month == time.April && day == 1 {
		return AprilFoolsEmojiRanges
	}

	// День космонавтики: 12 апреля
	if month == time.April && day == 12 {
		return CosmonauticsDayEmojiRanges
	}

	// Пасха: точное вычисление по православному календарю
	// Празднуем Светлую седмицу (7 дней после Пасхи)
	easterDate := OrthodoxEaster(year, datetime.NowLocation())
	easterEnd := easterDate.AddDate(0, 0, 7)
	if (now.After(easterDate) || now.Equal(easterDate)) && now.Before(easterEnd) {
		return EasterEmojiRanges
	}

	// День России: 12 июня
	if month == time.June && day == 12 {
		return RussiaDayEmojiRanges
	}

	// Хэллоуин: 31 октября
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
	default: // Компилятор "якобы" умнее меня, не догадывается, ведь сезонов то 4 *trollface.png*
		return Summer
	}
}

// OrthodoxEaster возвращает дату Православной Пасхи в григорианском календаре
// для заданного года в указанной временной зоне.
// Алгоритм: Пасха вычисляется в юлианском календаре (Meeus Julian algorithm),
// затем переводится в григорианскую дату через разницу календарей:
//   Δ = y/100 - y/400 - 2  (в сутках; верно для дат после 1600-03-01)
func OrthodoxEaster(year int, loc *time.Location) time.Time {
  // Шаг 1: Пасха в юлианском календаре (число марта/апреля по ЮК)
  a := year % 4
  b := year % 7
  c := year % 19
  d := (19*c + 15) % 30
  e := (2*a + 4*b - d + 34) % 7

  // Юлианская дата Пасхи: 22 марта + d + e (если >31 — это апрель)
  julianDay := 22 + d + e
  var month time.Month
  var day int
  if julianDay > 31 {
    month = time.April
    day = julianDay - 31
  } else {
    month = time.March
    day = julianDay
  }

  // Шаг 2: разница ЮК→ГК для данного года (суток)
  // Для 1900–2099 это 13, для 2100–2199 — 14 и т.д.
  delta := year/100 - year/400 - 2

  // Шаг 3: создаём "юлианскую" дату в Go (как григорианскую)
  // и прибавляем разницу календарей — получаем григорианскую Пасху.
  julianAsGregorian := time.Date(year, month, day, 0, 0, 0, 0, loc)
  return julianAsGregorian.AddDate(0, 0, delta)
}