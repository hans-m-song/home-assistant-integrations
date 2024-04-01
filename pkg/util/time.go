package util

import (
	"fmt"
	"strings"
	"time"
)

// Midnight returns the time at midnight of the previous day
func Midnight() time.Time {
	t := time.Now()
	t = t.Truncate(24 * time.Hour)
	return t
}

// ParseTimeDate parses a time date string in the format "HH:MM DD/MM/YYYY" and returns a time.Time object.
func ParseTimeDate(timeDate string) (*time.Time, error) {
	timeRaw, dateRaw, ok := strings.Cut(timeDate, " ")
	if !ok {
		return nil, fmt.Errorf("invalid time date format: %s", timeDate)
	}

	// attempt to manipulate format into YYYY-MM-DD HH:MM:SS
	dateParts := strings.Split(dateRaw, "/")
	timeParts := strings.Split(timeRaw, ":")
	if len(dateParts) != 3 {
		return nil, fmt.Errorf("invalid date format: %s", dateRaw)
	}

	if len(timeParts) != 2 {
		return nil, fmt.Errorf("invalid time format: %s", timeRaw)
	}

	dateTime := fmt.Sprintf(
		"%s-%s-%s %s:%s:00",
		dateParts[2],
		dateParts[1],
		dateParts[0],
		timeParts[0],
		timeParts[1],
	)

	val, err := time.Parse(time.DateTime, dateTime)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp: %s", err)
	}

	return &val, nil
}
