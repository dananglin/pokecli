package pokecache

import (
	"maps"
	"sync"
	"time"
)

type Cache struct {
	stopChan chan struct{}
	mu       *sync.Mutex
	entries  map[string]cacheEntry
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		stopChan: make(chan struct{}),
		mu:       &sync.Mutex{},
		entries:  make(map[string]cacheEntry),
		interval: interval,
	}

	go cache.readLoop()

	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}

	c.entries[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	value, exist := c.entries[key]

	return value.val, exist
}

// func (c *Cache) Stop() {
// 	c.stopChan <- struct{}{}
// }

func (c *Cache) readLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

readloop:
	for {
		select {
		case <-c.stopChan:
			break readloop
		case <-ticker.C:
			c.cleanupEntries()
		}
	}
}

func (c *Cache) cleanupEntries() {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiredTime := time.Now().Add(-c.interval)

	for key := range maps.All(c.entries) {
		if c.entries[key].createdAt.Before(expiredTime) {
			delete(c.entries, key)
		}
	}
}
