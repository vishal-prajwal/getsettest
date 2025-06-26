package fuzzychecker

import (
	"sort"
	"strings"

	fuzzy "github.com/paul-mannino/go-fuzzywuzzy"
)

func FuzzyMatch(name1, name2 string) (int, ExactMatch) {

	// Use fuzzywuzzy to get token set and partial ratios
	fuzzyScore := FuzzyMatching(name1, name2)
	exactScore := ExactMatching(name1, name2)

	return fuzzyScore, exactScore
}

func FuzzyMatching(original, target string) int {
	originals := strings.TrimSpace(strings.ToLower(original))
	targets := strings.TrimSpace(strings.ToLower(target))
	return fuzzy.Ratio(originals, targets)
}

type ExactMatch struct {
	Score   int
	Ordered bool
}

func ExactMatching(original, target string) ExactMatch {
	originals := strings.Split(strings.TrimSpace(strings.ToLower(original)), " ")
	targets := strings.Split(strings.TrimSpace(strings.ToLower(target)), " ")
	var matchIndexes []int
	if len(originals) > len(targets) {
		for _, v := range targets {
			if contains(originals, v) {
				matchIndexes = append(matchIndexes, indexOf(originals, v))
			}
		}
	} else {
		for _, v := range originals {
			if contains(targets, v) {
				matchIndexes = append(matchIndexes, indexOf(targets, v))
			}
		}
	}

	var size int
	if len(originals) > len(targets) {
		size = len(originals)
	} else {
		size = len(targets)
	}
	exactMatchPercentage := float64(len(matchIndexes)) / float64(size) * 100
	sorted := isSorted(matchIndexes, originals, targets)
	return ExactMatch{int(exactMatchPercentage), sorted}
}

func ExactMatchingSecondDoc(original, target string) ExactMatch {
	originals := strings.Split(strings.TrimSpace(strings.ToLower(original)), " ")
	targets := strings.Split(strings.TrimSpace(strings.ToLower(target)), " ")
	var matchIndexes []int
	if len(originals) > len(targets) {
		for _, v := range targets {
			if containsWithAtmostOneEditDistance(originals, v) {
				matchIndexes = append(matchIndexes, indexOfWithAtmostOneEditDistance(originals, v))
			}
		}
	} else {
		for _, v := range originals {
			if containsWithAtmostOneEditDistance(targets, v) {
				matchIndexes = append(matchIndexes, indexOfWithAtmostOneEditDistance(targets, v))
			}
		}
	}

	var size int
	if len(originals) > len(targets) {
		size = len(originals)
	} else {
		size = len(targets)
	}
	exactMatchPercentage := float64(len(matchIndexes)) / float64(size) * 100
	sorted := isSorted(matchIndexes, originals, targets)
	return ExactMatch{int(exactMatchPercentage), sorted}
}

func MatchWithoutSpace(s, t string) bool {
	s_new := strings.TrimSpace(strings.ToLower(strings.ReplaceAll(s, " ", "")))
	t_new := strings.TrimSpace(strings.ToLower(strings.ReplaceAll(t, " ", "")))
	return fuzzy.EditDistance(s_new, t_new) <= 1
}

func MatchFirstAndLastName(original string, target string) bool {
	originals := strings.Split(strings.TrimSpace(strings.ToLower(original)), " ")
	targets := strings.Split(strings.TrimSpace(strings.ToLower(target)), " ")
	len_original := len(originals)
	len_target := len(targets)
	if len_original >= 2 && len_target >= 2 {
		if fuzzy.EditDistance(originals[0], targets[0]) <= 1 && fuzzy.EditDistance(originals[len_original-1], targets[len_target-1]) <= 1 {
			return true
		}
	}
	return false
}

func MatchInitials(original, target string) bool {
	originals := strings.Split(strings.TrimSpace(strings.ToLower(original)), " ")
	targets := strings.Split(strings.TrimSpace(strings.ToLower(target)), " ")
	var mp = make(map[byte]int)
	for _, v := range originals {
		if len(v) == 0 {
			continue
		}
		mp[v[0]] += 1
	}
	for _, v := range targets {
		if len(v) == 0 {
			continue
		}
		if _, ok := mp[v[0]]; !ok {
			return false
		}
	}

	return true
}

func MatchRemovingFirstLastSingleLetterName(original string, target string) bool {
	originals := strings.Split(strings.TrimSpace(strings.ToLower(original)), " ")
	targets := strings.Split(strings.TrimSpace(strings.ToLower(target)), " ")
	len_original := len(originals)
	len_targets := len(targets)

	if len(original) == 0 || len(targets) == 0 {
		return false
	}

	if len(original) > 0 &&
		len(targets) > 0 &&
		len(originals[0]) != 1 &&
		len(originals[len_original-1]) != 1 &&
		len(targets[0]) != 1 &&
		len(targets[len_targets-1]) != 1 {
		return false
	}

	start_originals := 0
	end_originals := len_original - 1
	start_targets := 0
	end_targets := len_targets - 1

	if len_original <= 2 && len_targets <= 2 {
		return false
	}
	if originals[start_originals][0] == targets[start_targets][0] && len(originals[start_originals]) == 1 || len(targets[start_targets]) == 1 {
		start_targets++
		start_originals++
	}
	if MatchWithoutSpace(combine(originals[start_originals:]), combine(targets[start_targets:])) {
		return true
	}
	if originals[end_originals][0] == targets[end_targets][0] && len(originals[end_originals]) == 1 || len(targets[end_targets]) == 1 {
		end_originals--
		end_targets--
	}
	return MatchWithoutSpace(combine(originals[:end_originals+1]), combine(targets[:end_targets+1]))
}

func combine(s []string) string {
	return strings.Join(s, "")
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsWithAtmostOneEditDistance(s []string, e string) bool {
	for _, a := range s {
		if fuzzy.EditDistance(a, e) <= 1 {
			return true
		}
	}
	return false
}

func indexOf(s []string, e string) int {
	for i, a := range s {
		if a == e {
			return i
		}
	}
	return -1
}

func indexOfWithAtmostOneEditDistance(s []string, e string) int {
	for i, a := range s {
		if fuzzy.EditDistance(a, e) <= 1 {
			return i
		}
	}
	return -1
}
func isSorted(matchIndexes []int, originals, targets []string) bool {
	sortedIndexes := make([]int, len(matchIndexes))
	copy(sortedIndexes, matchIndexes)
	sort.Ints(sortedIndexes)
	if len(matchIndexes) == 0 {
		return false
	}

	if len(sortedIndexes) == 1 {
		pos := sortedIndexes[0]
		if len(originals) <= pos {
			return false
		}

		if len(targets) <= pos {
			return false
		}

		return strings.EqualFold(originals[pos], targets[pos])
	}

	return equal(sortedIndexes, matchIndexes)
}

func equal(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
