package datetime

import "time"

func NowLocation() *time.Location {
	location, _ := time.LoadLocation("Europe/Moscow")
	return location
}

func NowTime() time.Time {
	location, _ := time.LoadLocation("Europe/Moscow")
	return time.Now().In(location)
}
