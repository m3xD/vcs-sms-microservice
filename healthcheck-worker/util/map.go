package util

import "sync"

type ConcurrentMap struct {
	sync.RWMutex
	items map[string]interface{}
}

// NewConcurrentMap creates a new concurrent map.
func NewConcurrentMap() *ConcurrentMap {
	return &ConcurrentMap{
		items: make(map[string]interface{}),
	}
}

// Set adds or updates an element in the map.
func (m *ConcurrentMap) Set(key string, value interface{}) {
	m.Lock()
	defer m.Unlock()
	m.items[key] = value
}

// Get retrieves an element from the map.
func (m *ConcurrentMap) Get(key string) (interface{}, bool) {
	m.RLock()
	defer m.RUnlock()
	value, exists := m.items[key]
	return value, exists
}

// Delete removes an element from the map.
func (m *ConcurrentMap) Delete(key string) {
	m.Lock()
	defer m.Unlock()
	delete(m.items, key)
}

// Items returns a copy of all items in the map for safe iteration.
func (m *ConcurrentMap) Items() map[string]interface{} {
	m.RLock()
	defer m.RUnlock()
	itemsCopy := make(map[string]interface{})
	for key, value := range m.items {
		itemsCopy[key] = value
	}
	return itemsCopy
}
