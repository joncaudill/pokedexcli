// internal/pokecache/pokecache.go
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
	mu    sync.Mutex
	entry map[string]cacheEntry
}

func NewCache(interval time.Duration) *Cache {
	theNewCache := &Cache{
		entry: make(map[string]cacheEntry),
	}
	theNewCache.reapLoop(interval)
	return theNewCache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entry[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entry[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for range ticker.C {
			c.reap()
		}
	}()
}

func (c *Cache) reap() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, entry := range c.entry {
		if time.Since(entry.createdAt) > time.Minute {
			delete(c.entry, key)
		}
	}
}
