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

	loadData func(...interface{}) (interface{}, error)
}

func NewCachedItem(key string, lifeSpan time.Duration, data interface{}, loader func(...interface{}) (interface{}, error)) *CachedItem {
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

func (item *CachedItem) SetDataLoader(f func(...interface{}) (interface{}, error)) {
	item.Lock()
	item.loadData = f
	item.Unlock()
}

func (item *CachedItem) Value(args ...interface{}) (interface{}, error) {
	var f func() (interface{}, error)
	if item.loadData != nil {
		f = func() (interface{}, error) { return item.loadData(args) }
	}

	return item.Fetch(f)
}

func (item *CachedItem) Fetch(f func() (interface{}, error)) (interface{}, error) {
	item.RLock()
	d := item.data
	item.RUnlock()

	if d != nil {
		return d, nil
	}

	item.Lock()
	defer item.Unlock()
	var err error
	// Double check
	if item.data == nil && f != nil {
		// Ensure we never have 2 AfterFunc simustaniously
		if item.cleanupTimer != nil {
			item.cleanupTimer.Stop()
		}
		logger.Log.Debugf("Load data for item %v", item.key)
		item.data, err = f()
		if err != nil {
			return nil, err
		}
		item.cleanupTimer = time.AfterFunc(item.lifeSpan, func() { item.expire() })
	}

	return item.data, nil
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
