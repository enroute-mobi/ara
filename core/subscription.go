package core

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type SubscriptionId string

type Subscription struct {
	sync.RWMutex
	model.ClockConsumer

	manager Subscriptions

	id         SubscriptionId
	kind       string
	externalId string

	resourcesByObjectID map[string]*SubscribedResource
	subscriptionOptions map[string]string
}

type SubscribedResource struct {
	sync.RWMutex

	Reference        model.Reference
	RetryCount       int
	SubscribedAt     time.Time
	SubscribedUntil  time.Time
	lastStates       map[string]lastState
	resourcesOptions map[string]string
}

func NewResource(ref model.Reference) SubscribedResource {
	ressource := SubscribedResource{
		Reference:        ref,
		lastStates:       make(map[string]lastState),
		resourcesOptions: make(map[string]string),
	}

	return ressource
}

func (sr *SubscribedResource) ResourcesOptions() map[string]string {
	return sr.resourcesOptions
}

func (sr *SubscribedResource) LastState(state string) (l lastState, ok bool) {
	sr.RLock()
	l, ok = sr.lastStates[state]
	sr.RUnlock()
	return
}

func (sr *SubscribedResource) SetLastState(s string, l lastState) {
	sr.Lock()
	sr.lastStates[s] = l
	sr.Unlock()
}

type APISubscription struct {
	Kind       string
	References []model.Reference
}

func (subscription *Subscription) SetDefinition(apisub *APISubscription) {
	subscription.kind = apisub.Kind
	for _, ref := range apisub.References {
		if ref.ObjectId != nil {
			subscription.CreateAddNewResource(ref)
		}
	}
}

func (subscription *Subscription) Id() SubscriptionId {
	return subscription.id
}

func (subscription *Subscription) Kind() string {
	return subscription.kind
}

func (subscription *Subscription) ExternalId() string {
	return subscription.externalId
}

func (subscription *Subscription) SubscriptionOption(key string) (o string) {
	subscription.RLock()
	o = subscription.subscriptionOptions[key]
	subscription.RUnlock()
	return
}

func (subscription *Subscription) SetSubscriptionOption(key, value string) {
	subscription.Lock()
	subscription.subscriptionOptions[key] = value
	subscription.Unlock()
}

func (subscription *Subscription) SetKind(kind string) {
	subscription.kind = kind
}

func (subscription *Subscription) SetExternalId(externalId string) {
	subscription.externalId = externalId
}

func (subscription *Subscription) Save() (ok bool) {
	ok = subscription.manager.Save(subscription)
	return
}

func (subscription *Subscription) Delete() (ok bool) {
	ok = subscription.manager.Delete(subscription)
	return
}

func (subscription *Subscription) ResourcesByObjectIDCopy() map[string]*SubscribedResource {
	m := make(map[string]*SubscribedResource)
	subscription.RLock()
	for k, v := range subscription.resourcesByObjectID {
		m[k] = v
	}
	subscription.RUnlock()
	return m
}

func (subscription *Subscription) MarshalJSON() ([]byte, error) {
	resources := make([]*SubscribedResource, 0)

	subscription.RLock()
	for _, resource := range subscription.resourcesByObjectID {
		resources = append(resources, resource)
	}
	subscription.RUnlock()

	aux := struct {
		Id         SubscriptionId        `json:"SubscriptionRef,omitempty"`
		ExternalId string                `json:"ExternalId,omitempty"`
		Kind       string                `json:",omitempty"`
		Resources  []*SubscribedResource `json:",omitempty"`
	}{
		Id:         subscription.id,
		ExternalId: subscription.externalId,
		Kind:       subscription.kind,
		Resources:  resources,
	}
	return json.Marshal(&aux)
}

func (subscription *Subscription) Resource(obj model.ObjectID) *SubscribedResource {
	subscription.RLock()
	sub, present := subscription.resourcesByObjectID[obj.String()]
	subscription.RUnlock()
	if !present {
		return nil
	}
	return sub
}

func (subscription *Subscription) Resources(now time.Time) []*SubscribedResource {
	ressources := []*SubscribedResource{}

	subscription.RLock()
	for _, ressource := range subscription.resourcesByObjectID {
		if ressource.SubscribedUntil.After(subscription.Clock().Now()) {
			ressources = append(ressources, ressource)
		}
	}
	subscription.RUnlock()
	return ressources
}

func (subscription *Subscription) AddNewResource(resource SubscribedResource) {
	subscription.Lock()
	subscription.resourcesByObjectID[resource.Reference.ObjectId.String()] = &resource
	subscription.Unlock()
}

func (subscription *Subscription) CreateAddNewResource(reference model.Reference) *SubscribedResource {
	logger.Log.Debugf("Create subscribed resource for %v", reference.ObjectId.String())

	resource := SubscribedResource{
		Reference:        reference,
		SubscribedUntil:  subscription.Clock().Now().Add(1 * time.Minute),
		lastStates:       make(map[string]lastState),
		resourcesOptions: make(map[string]string),
	}
	subscription.Lock()
	subscription.resourcesByObjectID[reference.ObjectId.String()] = &resource
	subscription.Unlock()
	return &resource
}

func (subscription *Subscription) DeleteResource(key string) {
	subscription.Lock()
	delete(subscription.resourcesByObjectID, key)
	subscription.Unlock()
}

func (subscription *Subscription) ResourcesLen() (i int) {
	subscription.RLock()
	i = len(subscription.resourcesByObjectID)
	subscription.RUnlock()
	return
}

type MemorySubscriptions struct {
	model.UUIDConsumer

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

type Subscriptions interface {
	model.UUIDInterface

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
	CancelBroadcastSubscriptions()
	CancelCollectSubscriptions()
	FindByRessourceId(id, kind string) []*Subscription
	FindByExternalId(externalId string) (*Subscription, bool)
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
		if subscription.externalId != "" {
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

func (manager *MemorySubscriptions) FindByRessourceId(id, kind string) []*Subscription {
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
	maxResource, _ := strconv.Atoi(manager.partner.Setting("subscriptions.maximum_resources"))
	if maxResource == 1 {
		return manager.New(kind)
	}

	manager.mutex.RLock()
	for _, subscription := range manager.byIdentifier {
		if subscription.Kind() == kind && (maxResource < 1 || subscription.ResourcesLen() < maxResource) {
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
		generator := manager.partner.Generator("subscription_identifier")
		subscription.id = SubscriptionId(generator.NewIdentifier(IdentifierAttributes{Id: manager.NewUUID()}))
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
