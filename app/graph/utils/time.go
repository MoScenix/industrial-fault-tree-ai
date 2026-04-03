package utils

import "time"

const TimeLayout = "2006-01-02 15:04:05"

func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(TimeLayout)
}

func FormatTimeNow() string {
	return time.Now().Format(TimeLayout)
}

func ParseTimePtr(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}

	t, err := time.Parse(TimeLayout, value)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
