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

/* Add a CachedItem to the CacheTable without any specific loader */
func (table *CacheTable) Add(key string, lifeSpan time.Duration, data any) *CachedItem {
	item := NewCachedItem(key, lifeSpan, data, nil)

	table.Lock()
	table.items[item.key] = item
	table.Unlock()

	return item
}

/* Remove a specific CachedItem */
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

/* Remove all CachedItems */
func (table *CacheTable) Clear() {
	table.Lock()
	for k, v := range table.items {
		v.expire()
		delete(table.items, k)
	}
	table.Unlock()
}

/* Get a CachedItem saved data or fetch it with its dataloader */
func (table *CacheTable) Value(key string, args ...any) (any, error) {
	table.RLock()
	r, ok := table.items[key]
	table.RUnlock()

	if !ok {
		return nil, ErrNotFound(key)
	}
	return r.Value(args)
}

/* Get a CachedItem saved data or fetch it with the given func */
func (table *CacheTable) Fetch(key string, f func() (any, error)) (any, error) {
	table.RLock()
	r, ok := table.items[key]
	table.RUnlock()

	if !ok {
		return nil, ErrNotFound(key)
	}
	return r.Fetch(f)
}
