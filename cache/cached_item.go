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

	data      any
	lifeSpan  time.Duration
	createdAt time.Time

	// To cancel AfterFunc
	cleanupTimer *time.Timer

	loadData func(...any) (any, error)
}

func NewCachedItem(key string, lifeSpan time.Duration, data any, loader func(...any) (any, error)) *CachedItem {
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

/* Unused for now: Define a data loader for the item */
func (item *CachedItem) SetDataLoader(f func(...any) (any, error)) {
	item.Lock()
	item.loadData = f
	item.Unlock()
}

/* Get the item saved data or fetch it with its dataloader */
func (item *CachedItem) Value(args ...any) (any, error) {
	var f func() (any, error)
	if item.loadData != nil {
		f = func() (any, error) { return item.loadData(args) }
	}

	return item.Fetch(f)
}

/* Get the item saved data or fetch it with the given func */
func (item *CachedItem) Fetch(f func() (any, error)) (any, error) {
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
		logger.Log.Debugf("Fetch data for cached item %v", item.key)
		item.data, err = f()
		if err != nil {
			return nil, err
		}
		// Set the expire method to execute after the item lifespan
		item.cleanupTimer = time.AfterFunc(item.lifeSpan, func() { item.expire() })
	}

	return item.data, nil
}

/* Thread safe stop the item cleanup Timer */
func (item *CachedItem) Stop() {
	item.Lock()
	if item.cleanupTimer != nil {
		item.cleanupTimer.Stop()
	}
	item.Unlock()
}

/* Internal method to delete the item data after its cleanup timer expires */
func (item *CachedItem) expire() {
	item.Lock()

	if item.cleanupTimer != nil {
		item.cleanupTimer.Stop()
	}

	logger.Log.Debugf("Cached item %v expired", item.key)
	item.data = nil

	item.Unlock()
}

/* Unused for now: Manually set a cached item data.  */
func (item *CachedItem) SetData(data any) {
	item.Lock()
	item.data = data
	item.Unlock()
}
