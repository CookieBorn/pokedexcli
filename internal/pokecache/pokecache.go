package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	data map[string]cacheEntry
	mu   sync.Mutex
}

func NewCache(interval time.Duration) *Cache {

	cache := &Cache{
		data: make(map[string]cacheEntry),
		mu:   sync.Mutex{},
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			cache.reapLoop(interval)
		}
	}()

	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	var newEntry cacheEntry
	newEntry.val = val
	newEntry.createdAt = time.Now()
	c.data[key] = newEntry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for cacheKey := range c.data {
		if key == cacheKey {
			return c.data[key].val, true
		}
	}
	return nil, false
}

func (c *Cache) reapLoop(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, cache := range c.data {
		if time.Since(cache.createdAt) > interval {
			delete(c.data, key)
		}
	}
}
