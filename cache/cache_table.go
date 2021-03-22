package cache

import (
	"sync"
	"time"
)

type CacheTable struct {
	sync.RWMutex

	items map[string]*CachedItem
}

func (table *CacheTable) Add(key string, lifeSpan time.Duration, data interface{}) *CachedItem {
	item := NewCachedItem(key, lifeSpan, data, nil)

	// Add item to cache.
	table.Lock()

	table.items[item.key] = item

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

func (table *CacheTable) Value(key string, args ...interface{}) (interface{}, error) {
	table.RLock()
	r, ok := table.items[key]
	table.RUnlock()

	if !ok {
		return nil, ErrNotFound(key)
	}
	return r.Value(args), nil
}
