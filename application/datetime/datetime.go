package datetime

import "time"

func NowLocation() *time.Location {
	location, _ := time.LoadLocation("Europe/Moscow")
	return location
}

func NowTime() time.Time {
	return time.Now().In(NowLocation())
}

// ParseUTC парсит дату из бэкэнда (формат 2026-02-11T21:00:00.000Z) и конвертирует в МСК
func ParseUTC(s string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		t, err = time.Parse("2006-01-02T15:04:05.000Z", s)
		if err != nil {
			return time.Time{}, err
		}
	}
	return t.In(NowLocation()), nil
}

// FormatDateMSK парсит UTC-дату с бэкэнда и форматирует в "02.01.2006" по МСК
func FormatDateMSK(s string) string {
	t, err := ParseUTC(s)
	if err != nil {
		return s
	}
	return t.Format("02.01.2006")
}

// FormatDateShortMSK парсит UTC-дату с бэкэнда и форматирует в "02.01.06" по МСК
func FormatDateShortMSK(s string) string {
	t, err := ParseUTC(s)
	if err != nil {
		return s
	}
	return t.Format("02.01.06")
}
