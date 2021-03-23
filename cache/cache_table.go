package cache

import (
	"sync"
	"time"
)

type CacheTable struct {
	sync.RWMutex

	items map[string]*CachedItem
}

func NewCacheTable() *CacheTable {
	return &CacheTable{
		items: make(map[string]*CachedItem),
	}
}

func (table *CacheTable) Add(key string, lifeSpan time.Duration, data interface{}) *CachedItem {
	item := NewCachedItem(key, lifeSpan, data, nil)

	// Add item to cache.
	table.Lock()

	table.items[item.key] = item

	table.Unlock()

	return item
}

func (table *CacheTable) Delete(key string) (*CachedItem, error) {
	table.Lock()

	r, ok := table.items[key]
	if !ok {
		table.Unlock()
		return nil, ErrNotFound(key)
	}

	r.Stop()
	delete(table.items, key)

	table.Unlock()

	return r, nil
}

func (table *CacheTable) Clear() {
	table.Lock()
	for k, v := range table.items {
		v.expire()
		delete(table.items, k)
	}
	table.Unlock()
}

func (table *CacheTable) Value(key string, args ...interface{}) (interface{}, error) {
	table.RLock()
	r, ok := table.items[key]
	table.RUnlock()

	if !ok {
		return nil, ErrNotFound(key)
	}
	return r.Value(args)
}

func (table *CacheTable) Fetch(key string, f func() (interface{}, error)) (interface{}, error) {
	table.RLock()
	r, ok := table.items[key]
	table.RUnlock()

	if !ok {
		return nil, ErrNotFound(key)
	}
	return r.Fetch(f)
}
