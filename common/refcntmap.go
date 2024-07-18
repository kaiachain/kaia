package common

import (
	"sync"
)

// RefCntMap is a map with reference counting.
type RefCountingMap struct {
	mu     sync.RWMutex
	values map[interface{}]interface{}
	counts map[interface{}]int
}

func NewRefCountingMap() *RefCountingMap {
	return &RefCountingMap{
		values: make(map[interface{}]interface{}),
		counts: make(map[interface{}]int),
	}
}

// Get returns the value associated with the given key.
func (r *RefCountingMap) Get(key interface{}) (interface{}, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	value, ok := r.values[key]
	return value, ok
}

// Add adds an element to the map and increments its reference count.
func (r *RefCountingMap) Add(key interface{}, value interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.values[key] = value
	r.counts[key]++
}

// Remove decrements the reference count of the element with the given key.
func (r *RefCountingMap) Remove(key interface{}) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.counts[key] > 0 {
		r.counts[key]--
	}
	if r.counts[key] == 0 {
		delete(r.values, key)
		delete(r.counts, key)
	}
}

// Len returns the number of elements in the map.
func (r *RefCountingMap) Len() int {
	return len(r.values)
}
