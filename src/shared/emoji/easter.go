package emoji

import "time"

// OrthodoxEaster возвращает дату Православной Пасхи в григорианском календаре
// для заданного года в указанной временной зоне.
func OrthodoxEaster(year int, loc *time.Location) time.Time {
	a := year % 4
	b := year % 7
	c := year % 19
	d := (19*c + 15) % 30
	e := (2*a + 4*b - d + 34) % 7

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

	delta := year/100 - year/400 - 2

	julianAsGregorian := time.Date(year, month, day, 0, 0, 0, 0, loc)
	return julianAsGregorian.AddDate(0, 0, delta)
}
