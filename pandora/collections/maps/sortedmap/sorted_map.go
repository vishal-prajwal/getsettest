package sortedmap

import (
	"bitbucket.org/junglee_games/getsetgo/pandora/collections"
	"bitbucket.org/junglee_games/getsetgo/pandora/collections/maps"
)

// SortedMap - A Map that further provides a total ordering on its keys.
//
// The map is ordered according to the natural ordering of its keys, or by a
//
//	Comparator typically provided at sorted map creation time.
type SortedMap interface {
	// First returns the first(min) entry in the map
	First() *collections.Entry
	// Last returns the last(max) entry in the map
	Last() *collections.Entry
	// extends Map, Iterable and Comparable interfaces
	maps.Map
	collections.Iterable
	collections.Comparable
}
