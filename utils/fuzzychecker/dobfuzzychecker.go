package fuzzychecker

import (
	"strings"
)

func FuzzyCompareDates(dateStr1, dateStr2 string) bool {
	return strings.EqualFold(dateStr1, dateStr2)
}
