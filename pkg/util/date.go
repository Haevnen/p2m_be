package util

import (
	"time"
)

func Begin(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func ChangeToUTC(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), time.UTC)
}

func End(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 0, time.UTC)
}

func IsUTC(t time.Time) bool {
	return t.Location() == time.UTC
}

func ConvertToServerTimeZone(t time.Time) (*time.Time, error) {
	// Load the server's local timezone
	serverTimeZone, err := time.LoadLocation("Local")
	if err != nil {
		return nil, err
	}
	// Convert UTC to the server's timezone
	serverTime := t.In(serverTimeZone)
	return &serverTime, nil
}
