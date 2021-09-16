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

// Returns true if we send a DeleteSubscription request in the last 5 minutes
// Otherwise register it
func (cs *DeletedSubscriptions) AlreadySend(subID string) bool {
	cs.RLock()
	t, ok := cs.s[subID]
	cs.RUnlock()
	if ok && clock.DefaultClock().Now().Before(t.Add(DELETED_SUBSCRIPTION_TIMER)) {
		return true
	}

	cs.Lock()
	cs.s[subID] = clock.DefaultClock().Now()
	cs.Unlock()
	return false
}
