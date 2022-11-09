package core

import (
	"encoding/json"
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/core/ls"
	"bitbucket.org/enroute-mobi/ara/model"
)

type SubscribedResource struct {
	sync.RWMutex

	subscription *Subscription

	subscribedAt     time.Time
	lastStates       map[string]ls.LastState
	resourcesOptions map[string]string

	Reference       model.Reference
	RetryCount      int
	SubscribedUntil time.Time
}

func NewResource(ref model.Reference) *SubscribedResource {
	return &SubscribedResource{
		Reference:        ref,
		lastStates:       make(map[string]ls.LastState),
		resourcesOptions: make(map[string]string),
	}
}

func (sr *SubscribedResource) MarshalJSON() ([]byte, error) {
	type Alias SubscribedResource

	aux := struct {
		*Alias
		SubscribedAt time.Time
	}{
		SubscribedAt: sr.subscribedAt,
		Alias:        (*Alias)(sr),
	}
	return json.Marshal(&aux)
}

func (sr *SubscribedResource) ResourcesOptions() map[string]string {
	return sr.resourcesOptions
}

func (sr *SubscribedResource) LastState(state string) (l ls.LastState, ok bool) {
	sr.RLock()
	l, ok = sr.lastStates[state]
	sr.RUnlock()
	return
}

func (sr *SubscribedResource) SetLastState(s string, l ls.LastState) {
	sr.Lock()
	sr.lastStates[s] = l
	sr.Unlock()
}

func (sr *SubscribedResource) SubscribedAt() time.Time {
	return sr.subscribedAt
}

func (sr *SubscribedResource) Subscribed(t time.Time) {
	sr.subscribedAt = t
	if sr.subscription != nil {
		sr.subscription.subscribed = true
	}
}
