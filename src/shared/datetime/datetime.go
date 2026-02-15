package datetime

import (
	"strings"
	"time"
)

type LocalDateTime struct {
	time.Time
}

func (lt *LocalDateTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" || s == "null" {
		lt.Time = time.Time{}
		return nil
	}
	t, err := ParseUTC(s)
	if err != nil {
		return err
	}
	lt.Time = t
	return nil
}

func (lt LocalDateTime) MarshalJSON() ([]byte, error) {
	return lt.Time.MarshalJSON()
}

func (lt LocalDateTime) FormatDateMSK() string {
	if lt.Time.IsZero() {
		return ""
	}
	return lt.Time.Format("02.01.2006")
}

func NowLocation() *time.Location {
	location, _ := time.LoadLocation("Europe/Moscow")
	return location
}

func NowTime() time.Time {
	return time.Now().In(NowLocation())
}

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
