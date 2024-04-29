package main

import (
	"sync"
)

type TypedMap[K any, V any] struct {
	sync.Map
}

// LoadBytes safely loads the value associated with the key as []byte
// If the value is not a []byte, returns nil and false
func (tm *TypedMap[K, V]) Load(key K) (V, bool) {
	value, ok := tm.Map.Load(key)
	return value.(V), ok
}


func (tm *TypedMap[K, V]) Store(key K, value V) {
	tm.Map.Store(key, value)
}

func (tm *TypedMap[K, V]) Delete(key K) {
	tm.Map.Delete(key)
}
