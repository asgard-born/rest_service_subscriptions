package utils

import (
	"fmt"
	"time"
)

func ParseToMonthYear(input string) (time.Time, error) {
	layouts := []string{
		"2006-01-02",
		"02.01.2006",
		"01-2006",
		"2006-01",
		"01/2006",
		"2006/01",
	}

	var t time.Time
	var err error

	for _, layout := range layouts {
		t, err = time.Parse(layout, input)
		if err == nil {
			return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC), nil
		}
	}

	return time.Time{}, fmt.Errorf("invalid date format, expected MM-YYYY")
}
