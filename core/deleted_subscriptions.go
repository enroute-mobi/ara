package core

import (
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
)

const (
	DELETED_SUBSCRIPTION_TIMER = 5 * time.Minute
)

type DeletedSubscriptions struct {
	sync.RWMutex

	s map[string]time.Time
}

func NewDeletedSubscriptions() *DeletedSubscriptions {
	return &DeletedSubscriptions{
		s: make(map[string]time.Time),
	}
}

// Reset empties the cache while keeping the same instance, so callers can reset
// it on (re)start without reassigning the pointer concurrently with readers.
func (ds *DeletedSubscriptions) Reset() {
	ds.Lock()
	ds.s = make(map[string]time.Time)
	ds.Unlock()
}

// Returns true if we send a DeleteSubscription request in the last 5 minutes
// Otherwise register it
func (ds *DeletedSubscriptions) AlreadySend(subID string) bool {
	ds.RLock()
	t, ok := ds.s[subID]
	ds.RUnlock()
	if ok && clock.DefaultClock().Now().Before(t.Add(DELETED_SUBSCRIPTION_TIMER)) {
		return true
	}

	ds.Lock()
	ds.s[subID] = clock.DefaultClock().Now()
	ds.Unlock()
	return false
}
