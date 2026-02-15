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
