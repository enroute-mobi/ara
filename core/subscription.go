package core

import (
	"encoding/json"
	"time"

	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type SubscriptionId string

type Subscription struct {
	model.ClockConsumer
	manager Subscriptions

	id                  SubscriptionId
	kind                string
	externalId          string
	resourcesByObjectID map[string]*SubscribedResource
	subscriptionOptions map[string]string
}

type SubscribedResource struct {
	Reference        model.Reference
	RetryCount       int
	SubscribedAt     time.Time
	SubscribedUntil  time.Time
	LastStates       map[string]lastState `json:",omitempty"`
	resourcesOptions map[string]string
}

func (sr *SubscribedResource) ResourcesOptions() map[string]string {
	return sr.resourcesOptions
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

func (subscription *Subscription) SubscriptionOptions() map[string]string {
	return subscription.subscriptionOptions
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

func (subscription *Subscription) ResourcesByObjectID() map[string]*SubscribedResource {
	return subscription.resourcesByObjectID
}

func (subscription *Subscription) MarshalJSON() ([]byte, error) {
	resources := make([]*SubscribedResource, 0)
	for _, resource := range subscription.resourcesByObjectID {
		resources = append(resources, resource)
	}

	aux := struct {
		Id        SubscriptionId
		Kind      string                `json:",omitempty"`
		Resources []*SubscribedResource `json:",omitempty"`
	}{
		Id:        subscription.id,
		Kind:      subscription.kind,
		Resources: resources,
	}
	return json.Marshal(&aux)
}

func (subscription *Subscription) Resource(obj model.ObjectID) *SubscribedResource {
	sub, present := subscription.resourcesByObjectID[obj.String()]
	if !present {
		return nil
	}
	return sub
}

func (subscription *Subscription) Resources(now time.Time) []*SubscribedResource {
	ressources := []*SubscribedResource{}

	for _, ressource := range subscription.resourcesByObjectID {
		if ressource.SubscribedUntil.After(subscription.Clock().Now()) {
			ressources = append(ressources, ressource)
		}
	}
	return ressources
}

func (subscription *Subscription) CreateAddNewResource(reference model.Reference) *SubscribedResource {
	logger.Log.Debugf("Create subscribed resource for %v", reference.ObjectId.String())

	ressource := SubscribedResource{
		Reference:        reference,
		SubscribedUntil:  subscription.Clock().Now().Add(1 * time.Minute),
		LastStates:       make(map[string]lastState),
		resourcesOptions: make(map[string]string),
	}
	subscription.resourcesByObjectID[reference.ObjectId.String()] = &ressource
	return &ressource
}

type MemorySubscriptions struct {
	model.UUIDConsumer

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
	FindOrCreateByKind(string) (*Subscription, bool)
	FindByKind(string) (*Subscription, bool)
	Save(Subscription *Subscription) bool
	Delete(Subscription *Subscription) bool
	DeleteById(id SubscriptionId)
	CancelSubscriptions()
	FindByRessourceId(id string) (*Subscription, bool)
	FindByExternalId(externalId string) (*Subscription, bool)
}

func NewMemorySubscriptions(partner *Partner) *MemorySubscriptions {
	return &MemorySubscriptions{
		byIdentifier: make(map[SubscriptionId]*Subscription),
		partner:      partner,
	}
}

func (manager *MemorySubscriptions) New(kind string) *Subscription {
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
	for _, subscription := range manager.byIdentifier {
		if subscription.Kind() == kind {
			return subscription, true
		}
	}
	return nil, false
}

func (manager *MemorySubscriptions) FindByExternalId(externalId string) (*Subscription, bool) {
	for _, subscription := range manager.byIdentifier {
		if subscription.ExternalId() == externalId {
			return subscription, true
		}
	}
	return nil, false
}

func (manager *MemorySubscriptions) FindByRessourceId(id string) (*Subscription, bool) {
	for _, subscription := range manager.byIdentifier {
		_, ok := subscription.resourcesByObjectID[id]
		if !ok {
			continue
		}
		return subscription, true
	}
	return nil, false
}

func (manager *MemorySubscriptions) FindOrCreateByKind(kind string) (*Subscription, bool) {
	for _, subscription := range manager.byIdentifier {
		if subscription.Kind() == kind {
			return subscription, true
		}
	}

	subscription := manager.New(kind)
	return subscription, false
}

func (manager *MemorySubscriptions) Find(id SubscriptionId) (*Subscription, bool) {
	subscription, ok := manager.byIdentifier[id]
	if ok {
		return subscription, true
	} else {
		return nil, false
	}
}

func (manager *MemorySubscriptions) FindAll() (subscriptions []*Subscription) {
	if len(manager.byIdentifier) == 0 {
		return []*Subscription{}
	}
	for _, subscription := range manager.byIdentifier {
		subscriptions = append(subscriptions, subscription)
	}
	return
}

func (manager *MemorySubscriptions) Save(subscription *Subscription) bool {
	if subscription.Id() == "" {
		subscription.id = SubscriptionId(manager.NewUUID())
	}
	subscription.manager = manager
	manager.byIdentifier[subscription.Id()] = subscription
	return true
}

func (manager *MemorySubscriptions) Delete(subscription *Subscription) bool {
	delete(manager.byIdentifier, subscription.Id())
	return true
}

func (manager *MemorySubscriptions) DeleteById(id SubscriptionId) {
	delete(manager.byIdentifier, id)
}

func (manager *MemorySubscriptions) CancelSubscriptions() {
	for id := range manager.byIdentifier {
		delete(manager.byIdentifier, id)
	}
}
