package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	entries  map[string]cacheEntry
	mu       *sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(inter time.Duration) Cache {
	res := Cache{
		entries:  make(map[string]cacheEntry),
		mu:       &sync.Mutex{},
		interval: inter,
	}
	go func() {
		tick := time.NewTicker(inter)
		for range tick.C {
			res.reapLoop()
		}
	}()
	return res
}

func (c *Cache) Add(key string, value []byte) {
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	if value, found := c.entries[key]; found {
		return value.val, true
	}
	return nil, false
}

func (c *Cache) reapLoop() {
	for key, val := range c.entries {
		if time.Now().Nanosecond()-val.createdAt.Nanosecond() > int(c.interval.Nanoseconds()) {
			c.mu.Lock()
			delete(c.entries, key)
			c.mu.Unlock()
		}
	}
}
