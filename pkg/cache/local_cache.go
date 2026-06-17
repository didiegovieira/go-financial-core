package cache

import (
	"sync"
	"time"
)

type LocalCache struct {
	mu    sync.RWMutex
	store map[string]CacheItem
}

type CacheItem struct {
	Value      interface{}
	Expiration int64
}

func NewLocalCache() *LocalCache {
	return &LocalCache{
		store: make(map[string]CacheItem),
	}
}

func (c *LocalCache) Set(key string, value interface{}, duration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = CacheItem{
		Value:      value,
		Expiration: time.Now().Add(duration).UnixNano(),
	}
}

func (c *LocalCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.store[key]
	if !found {
		return nil, false
	}

	if time.Now().UnixNano() > item.Expiration {
		delete(c.store, key)
		return nil, false
	}

	return item.Value, true
}

func (c *LocalCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
}
