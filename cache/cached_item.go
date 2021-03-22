package cache

import (
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

const (
	MIN_CACHE_LIFESPAN     = 10 * time.Second
	DEFAULT_CACHE_LIFESPAN = 60 * time.Second
)

type CachedItem struct {
	sync.RWMutex

	key string

	data      interface{}
	lifeSpan  time.Duration
	createdAt time.Time

	// To cancel AfterFunc
	cleanupTimer *time.Timer

	loadData func(...interface{}) interface{}
}

func NewCachedItem(key string, lifeSpan time.Duration, data interface{}, loader func(...interface{}) interface{}) *CachedItem {
	if lifeSpan < MIN_CACHE_LIFESPAN {
		lifeSpan = DEFAULT_CACHE_LIFESPAN
	}

	return &CachedItem{
		key:       key,
		lifeSpan:  lifeSpan,
		createdAt: time.Now(),
		data:      data,
		loadData:  loader,
	}
}

func (item *CachedItem) SetDataLoader(f func(...interface{}) interface{}) {
	item.Lock()
	item.loadData = f
	item.Unlock()
}

func (item *CachedItem) Value(args ...interface{}) interface{} {
	item.RLock()
	if item.data != nil {
		item.RUnlock()
		return item.data
	}
	item.RUnlock()

	item.Lock()
	// Double check
	if item.data == nil && item.loadData != nil {
		logger.Log.Debugf("Load data for item %v", item.key)
		item.data = item.loadData(args)
		item.cleanupTimer = time.AfterFunc(item.lifeSpan, func() { item.expire() })
	}

	item.Unlock()
	return item.data
}

func (item *CachedItem) Stop() {
	item.Lock()
	if item.cleanupTimer != nil {
		item.cleanupTimer.Stop()
	}
	item.Unlock()
}

func (item *CachedItem) expire() {
	item.Lock()

	if item.cleanupTimer != nil {
		item.cleanupTimer.Stop()
	}

	logger.Log.Debugf("Cached item %v expired", item.key)
	item.data = nil

	item.Unlock()
}

func (item *CachedItem) SetData(data interface{}) {
	item.Lock()
	item.data = data
	item.Unlock()
}
