package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheEntries map[string]cacheEntry
	mu           *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	data      []byte
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		cacheEntries: make(map[string]cacheEntry),
		mu:           &sync.Mutex{},
	}
	cache.reapLoop(interval)
	return cache
}

func (c *Cache) Add(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cacheEntries[key] = cacheEntry{
		createdAt: time.Now(),
		data:      data,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cacheEntry, found := c.cacheEntries[key]
	return cacheEntry.data, found
}

func (c *Cache) reapLoop(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	ticker := time.NewTicker(interval)

	go func() {
		for {
			t := <-ticker.C
			for key, cacheEntry := range c.cacheEntries {
				if t.Compare(cacheEntry.createdAt.Add(interval)) == 1 {
					delete(c.cacheEntries, key)
				}
			}
		}
	}()
}
