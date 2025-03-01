package util

import "time"

func ParseTimeFromString(t string) (time.Time, error) {
	return time.Parse("15:04:05", t)
}
