package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	data      []byte
}

type Cache struct {
	Entries map[string]cacheEntry
	mu      sync.RWMutex
	stopCh  chan struct{}
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		Entries: make(map[string]cacheEntry),
		stopCh:  make(chan struct{}),
	}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) Add(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Entries[key] = cacheEntry{
		createdAt: time.Now(),
		data:      data,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	entry, exists := c.Entries[key]
	if !exists {
		return nil, false
	}
	return entry.data, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for key, entry := range c.Entries {
				if time.Since(entry.createdAt) > interval {
					delete(c.Entries, key)
				}
			}
			c.mu.Unlock()
		case <-c.stopCh:
			return
		}
	}
}

func (c *Cache) StopReaping() {
	close(c.stopCh)
}
