package concurrentmap

import "bitbucket.org/junglee_games/getsetgo/pandora/collections/maps"

// ConcurrentMap interface that all concurrent maps implement
type ConcurrentMap interface {
	// PutIfAbsent inserts an entry into the map, if the key doesn't exists
	PutIfAbsent(key interface{}, value interface{})

	// extends Map Interface
	maps.Map
}
