package fuzzychecker

import (
	"fmt"
	"time"
)

func FuzzyCompareDates(dateStr1, dateStr2 string) (int, error) {
	layout := "2006-01-02" // assuming same format for both

	date1, err1 := time.Parse(layout, dateStr1)
	date2, err2 := time.Parse(layout, dateStr2)
	if err1 != nil || err2 != nil {
		return 0, fmt.Errorf("invalid date format")
	}

	y1, m1, _ := date1.Date()
	y2, m2, _ := date2.Date()

	if y1 != y2 {
		return 0, nil
	}

	monthDiff := int(m1) - int(m2)
	if monthDiff < 0 {
		monthDiff = -monthDiff
	}

	return 1, nil
}
