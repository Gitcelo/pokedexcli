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
	cache map[string]cacheEntry
	mu    sync.Mutex
}

func NewCache(interval time.Duration) *Cache {
	c := Cache{
		cache: make(map[string]cacheEntry),
	}
	go c.reapLoop(interval)
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	ce := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = ce
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	ce, ok := c.cache[key]
	if !ok {
		return []byte{}, false
	}
	return ce.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		t := <-ticker.C
		c.mu.Lock()
		for key, ce := range c.cache {
			age := t.Sub(ce.createdAt)
			if age > interval {
				delete(c.cache, key)
			}
		}
		c.mu.Unlock()
	}
}
