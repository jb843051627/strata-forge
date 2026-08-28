package engine

import "sync"

type StateCache struct {
	mu     sync.RWMutex
	values map[int64]string
}

func NewStateCache() *StateCache {
	return &StateCache{values: make(map[int64]string)}
}

func (c *StateCache) Set(key int64, value string) {
	c.mu.Lock()
	c.values[key] = value
	c.mu.Unlock()
}

func (c *StateCache) Get(key int64) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok := c.values[key]
	return value, ok
}

func (c *StateCache) Delete(key int64) {
	c.mu.Lock()
	delete(c.values, key)
	c.mu.Unlock()
}

func (c *StateCache) Snapshot() map[int64]string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	result := make(map[int64]string, len(c.values))
	for key, value := range c.values {
		result[key] = value
	}
	return c.values
}
