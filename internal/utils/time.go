package utils

import "time"

func ParseTime(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, value)
}

func ParseTimePointer(value string) (*time.Time, error) {
	parsed, err := ParseTime(value)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}
