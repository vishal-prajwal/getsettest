package fuzzychecker

import (
	"fmt"
	"time"
)

func FuzzyCompareDates(dateStr1, dateStr2 string) (bool, error) {
	layout := "2006-01-02" // assuming same format for both

	date1, err1 := time.Parse(layout, dateStr1)
	date2, err2 := time.Parse(layout, dateStr2)
	if err1 != nil || err2 != nil {
		return false, fmt.Errorf("invalid date format")
	}

	y1, m1, _ := date1.Date()
	y2, m2, _ := date2.Date()

	if y1 != y2 {
		return false, nil
	}

	monthDiff := int(m1) - int(m2)
	if monthDiff < 0 {
		monthDiff = -monthDiff
	}

	return monthDiff <= 5, nil
}
