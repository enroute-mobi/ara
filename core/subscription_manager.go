package core

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type Subscriptions interface {
	uuid.UUIDInterface

	New(string) *Subscription
	Find(SubscriptionId) (*Subscription, bool)
	FindAll() []*Subscription
	FindOrCreateByKind(string) *Subscription
	FindByKind(string) (*Subscription, bool)
	FindSubscriptionsByKind(string) []*Subscription
	FindBroadcastSubscriptions() []*Subscription
	Index(*Subscription)
	Save(*Subscription) bool
	Delete(*Subscription) bool
	DeleteById(SubscriptionId)
	CancelSubscriptions()
	CancelSubscriptionsResourcesBefore(time.Time)
	CancelBroadcastSubscriptions()
	CancelCollectSubscriptions()
	FindByResourceId(id, kind string) []*Subscription
	FindByExternalId(string) (*Subscription, bool)
}

type MemorySubscriptions struct {
	uuid.UUIDConsumer

	mutex   *sync.RWMutex
	partner *Partner

	byIdentifier        map[SubscriptionId]*Subscription
	byKindAndResourceId map[string][]SubscriptionId
}

func (manager *MemorySubscriptions) MarshalJSON() ([]byte, error) {
	subscriptions := make([]*Subscription, 0)
	for _, subscription := range manager.byIdentifier {
		subscriptions = append(subscriptions, subscription)
	}

	return json.Marshal(subscriptions)
}

func NewMemorySubscriptions(partner *Partner) *MemorySubscriptions {
	return &MemorySubscriptions{
		mutex:               &sync.RWMutex{},
		byIdentifier:        make(map[SubscriptionId]*Subscription),
		byKindAndResourceId: make(map[string][]SubscriptionId),
		partner:             partner,
	}
}

func (manager *MemorySubscriptions) New(kind string) *Subscription {
	logger.Log.Debugf("Creating subscription with kind %v", kind)
	subscription := &Subscription{
		kind:                kind,
		manager:             manager,
		resourcesByObjectID: make(map[string]*SubscribedResource),
		subscriptionOptions: make(map[string]string),
	}
	subscription.Save()
	return subscription
}

func (manager *MemorySubscriptions) FindByKind(kind string) (*Subscription, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	for _, subscription := range manager.byIdentifier {
		if subscription.Kind() == kind {
			return subscription, true
		}
	}
	return nil, false
}

func (manager *MemorySubscriptions) FindSubscriptionsByKind(kind string) (subscriptions []*Subscription) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	for _, subscription := range manager.byIdentifier {
		if subscription.Kind() == kind {
			subscriptions = append(subscriptions, subscription)
		}
	}
	return subscriptions
}

func (manager *MemorySubscriptions) FindBroadcastSubscriptions() (subscriptions []*Subscription) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	for _, subscription := range manager.byIdentifier {
		if strings.HasSuffix(subscription.Kind(), "Broadcast") {
			subscriptions = append(subscriptions, subscription)
		}
	}
	return subscriptions
}

func (manager *MemorySubscriptions) FindByExternalId(externalId string) (*Subscription, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	for _, subscription := range manager.byIdentifier {
		if subscription.ExternalId() == externalId {
			return subscription, true
		}
	}
	return nil, false
}

func (manager *MemorySubscriptions) FindByResourceId(id, kind string) []*Subscription {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	kindAndResourceId := fmt.Sprintf("%s-%s", kind, id)
	subscriptionIdentifiers, found := manager.byKindAndResourceId[kindAndResourceId]

	subscriptions := []*Subscription{}

	if found {
		for _, subscriptionId := range subscriptionIdentifiers {
			subscription, ok := manager.byIdentifier[subscriptionId]
			if !ok {
				// Subscription no longer exists
				continue
			}

			if subscription.kind != kind {
				// Subscription wrong kind ? changed kind ?
				continue
			}

			subscription.RLock()
			_, ok = subscription.resourcesByObjectID[id]
			subscription.RUnlock()

			if !ok {
				// Resource no longer presents in subscription
				continue
			}

			subscriptions = append(subscriptions, subscription)
		}
	}

	if len(subscriptions) > 0 {
		if len(subscriptions) != len(subscriptionIdentifiers) {
			// Some of subscriptions are no longer associated to this kindAndResourceId
			subscriptionIds := []SubscriptionId{}
			for _, subscription := range subscriptions {
				subscriptionIds = append(subscriptionIds, subscription.Id())
			}

			// "Update" lock to RW one
			manager.mutex.RUnlock()
			manager.mutex.Lock()

			manager.byKindAndResourceId[kindAndResourceId] = subscriptionIds

			manager.mutex.Unlock()
			manager.mutex.RLock()
		}
	} else {
		// No subscription found for this kindAndResourceId

		// "Update" lock to RW one
		manager.mutex.RUnlock()
		manager.mutex.Lock()

		delete(manager.byKindAndResourceId, kindAndResourceId)

		manager.mutex.Unlock()
		manager.mutex.RLock()
	}

	return subscriptions
}

func (manager *MemorySubscriptions) FindOrCreateByKind(kind string) *Subscription {
	maxResource := manager.partner.SubscriptionMaximumResources()
	if maxResource == 1 {
		return manager.New(kind)
	}

	manager.mutex.RLock()
	for _, subscription := range manager.byIdentifier {
		if subscription.Kind() == kind && (maxResource < 1 || subscription.ResourcesLen() < maxResource) && !subscription.subscribed {
			manager.mutex.RUnlock()
			return subscription
		}
	}
	manager.mutex.RUnlock()

	return manager.New(kind)
}

func (manager *MemorySubscriptions) Find(id SubscriptionId) (*Subscription, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	subscription, ok := manager.byIdentifier[id]
	if ok {
		return subscription, true
	} else {
		return nil, false
	}
}

func (manager *MemorySubscriptions) FindAll() (subscriptions []*Subscription) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	if len(manager.byIdentifier) == 0 {
		return []*Subscription{}
	}
	for _, subscription := range manager.byIdentifier {
		subscriptions = append(subscriptions, subscription)
	}
	return
}

func (manager *MemorySubscriptions) Save(subscription *Subscription) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if subscription.Id() == "" {
		generator := manager.partner.SubscriptionIdentifierGenerator()
		subscription.id = SubscriptionId(generator.NewIdentifier(idgen.IdentifierAttributes{Id: manager.NewUUID()}))
	}

	subscription.manager = manager
	manager.byIdentifier[subscription.Id()] = subscription

	manager.unsafeIndex(subscription)

	return true
}

// Index the subscription with all associated resources
func (manager *MemorySubscriptions) Index(subscription *Subscription) {
	manager.mutex.Lock()
	manager.unsafeIndex(subscription)
	manager.mutex.Unlock()
}

// Unsafe method, need to handle manager mutex before calling
func (manager *MemorySubscriptions) unsafeIndex(subscription *Subscription) {
	subscription.RLock()

	for resourceId := range subscription.resourcesByObjectID {
		kindAndResourceId := fmt.Sprintf("%s-%s", subscription.Kind(), resourceId)
		subscriptionIdentifiers, found := manager.byKindAndResourceId[kindAndResourceId]

		if !found {
			// No subscription associated to this kindAndResourceId, create a new SubscriptionId slice
			manager.byKindAndResourceId[kindAndResourceId] = []SubscriptionId{subscription.Id()}
			continue
		}

		// Is the Subscription already associated to this kindAndResourceId ?
		if newSubscriptionForKindAndResourceId(subscriptionIdentifiers, subscription.Id()) {
			// Associate this Subscription to the kindAndResourceId
			manager.byKindAndResourceId[kindAndResourceId] = append(subscriptionIdentifiers, subscription.Id())
		}
	}
	subscription.RUnlock()
}

// Check if a SubscriptionId is in a Slice of SubscriptionIds
func newSubscriptionForKindAndResourceId(s []SubscriptionId, id SubscriptionId) bool {
	for _, subscriptionId := range s {
		if subscriptionId == id {
			// The Subscription is already associated to this kindAndResourceId
			return false
		}
	}
	return true
}

func (manager *MemorySubscriptions) Delete(subscription *Subscription) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	delete(manager.byIdentifier, subscription.Id())
	return true
}

func (manager *MemorySubscriptions) DeleteById(id SubscriptionId) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	delete(manager.byIdentifier, id)
}

func (manager *MemorySubscriptions) CancelSubscriptionsResourcesBefore(time time.Time) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	for _, sub := range manager.byIdentifier {
		for key, resource := range sub.ResourcesByObjectIDCopy() {
			if resource.SubscribedAt().After(time) || resource.SubscribedAt().IsZero() {
				continue
			}
			sub.DeleteResource(key)
			logger.Log.Debugf("Deleting ressource %v from subscription with id %v after partner reload", key, sub.Id())

		}
		if sub.ResourcesLen() == 0 {
			delete(manager.byIdentifier, sub.Id())
		}
	}
}

func (manager *MemorySubscriptions) CancelSubscriptions() {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	for id := range manager.byIdentifier {
		delete(manager.byIdentifier, id)
	}
}

func (manager *MemorySubscriptions) CancelBroadcastSubscriptions() {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	for id, subscription := range manager.byIdentifier {
		if subscription.externalId != "" {
			delete(manager.byIdentifier, id)
		}
	}
}

func (manager *MemorySubscriptions) CancelCollectSubscriptions() {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	for id, subscription := range manager.byIdentifier {
		if subscription.externalId == "" {
			delete(manager.byIdentifier, id)
		}
	}
}
