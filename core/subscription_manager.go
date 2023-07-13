package core

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type Subscriptions interface {
	uuid.UUIDInterface

	New(kind string) *Subscription
	Find(id SubscriptionId) (*Subscription, bool)
	FindAll() []*Subscription
	FindOrCreateByKind(string) *Subscription
	FindByKind(string) (*Subscription, bool)
	FindSubscriptionsByKind(string) []*Subscription
	FindBroadcastSubscriptions() []*Subscription
	Save(Subscription *Subscription) bool
	Delete(Subscription *Subscription) bool
	DeleteById(id SubscriptionId)
	CancelSubscriptions()
	CancelSubscriptionsResourcesBefore(time.Time)
	CancelBroadcastSubscriptions()
	CancelCollectSubscriptions()
	FindByResourceId(id, kind string) []*Subscription
	FindByExternalId(externalId string) (*Subscription, bool)
}

type MemorySubscriptions struct {
	uuid.UUIDConsumer

	mutex   *sync.RWMutex
	partner *Partner

	byIdentifier map[SubscriptionId]*Subscription
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
		mutex:        &sync.RWMutex{},
		byIdentifier: make(map[SubscriptionId]*Subscription),
		partner:      partner,
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
	subscriptions := []*Subscription{}

	for _, subscription := range manager.byIdentifier {
		subscription.RLock()
		_, ok := subscription.resourcesByObjectID[id]
		subscription.RUnlock()
		if ok && subscription.kind == kind {
			subscriptions = append(subscriptions, subscription)
		}
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
