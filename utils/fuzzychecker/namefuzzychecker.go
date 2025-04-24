package fuzzychecker

import (
	"regexp"
	"strings"

	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
)

// List of salutation words to remove
var salutationWords = []string{
	"mr", "mrs", "ms", "bhai", "bhau", "bhoi", "bai", "kumar", "kumr", "kmr",
	"ben", "devi", "debi", "saheb", "kumaar", "sab", "kumari", "lal", "w", "s",
	"master", "singh",
}

// Replacement rules for similar sounding words
var phoneticReplacements = []struct {
	pattern     string
	replacement string
}{
	{"bh", "b"}, {"th", "t"}, {"dh", "d"}, {"sh", "s"}, {"ck", "q"},
	{"gh", "g"}, {"kh", "q"}, {"ch", "c"}, {"ph", "f"}, {"v", "b"},
	{"k", "q"}, {"w", "b"}, {"z", "j"},
}

func preprocessName(name string) string {
	name = strings.ToLower(name)

	// Remove salutation words
	for _, word := range salutationWords {
		regex := regexp.MustCompile(`\b` + word + `\b`)
		name = regex.ReplaceAllString(name, "")
	}

	// Remove junk characters
	regex := regexp.MustCompile(`[^A-Za-z ]+`)
	name = regex.ReplaceAllString(name, " ")

	// Replace similar sounding words
	for _, rule := range phoneticReplacements {
		name = strings.ReplaceAll(name, rule.pattern, rule.replacement)
	}

	// Remove extra spaces
	name = strings.Join(strings.Fields(name), " ")
	return name
}

func FuzzyMatch(name1, name2 string) int {
	pre1 := preprocessName(name1)
	pre2 := preprocessName(name2)

	// Use fuzzywuzzy to get token set and partial ratios
	tokenSet := fuzzy.TokenSetRatio(pre1, pre2)
	partial := fuzzy.PartialRatio(pre1, pre2)

	// Return max of the two
	if tokenSet > partial {
		return tokenSet
	}
	return partial
}
