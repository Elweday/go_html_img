package main

import (
	"sync"
)

type TypedMap[K any, V any] struct {
	sync.Map
}

// Load safely loads the value associated with the key
// If the value does not exist returns nil and false
func (tm *TypedMap[K, V]) Load(key K) (V, bool) {
	value, ok1 := tm.Map.Load(key)
	casted, ok2 := value.(V)
	return casted, ok1 && ok2
}


func (tm *TypedMap[K, V]) Store(key K, value V) {
	tm.Map.Store(key, value)
}

func (tm *TypedMap[K, V]) Delete(key K) {
	tm.Map.Delete(key)
}
