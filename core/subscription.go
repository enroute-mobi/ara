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

	id                  SubscriptionId
	kind                string
	partner             *Partner
	resourcesByObjectID map[string]*SubscribedResource
}

func NewSubscription() *Subscription {
	return &Subscription{
		resourcesByObjectID: make(map[string]*SubscribedResource),
	}
}

type SubscribedResource struct {
	Reference       model.Reference
	SubscribedAt    time.Time
	SubscribedUntil time.Time
}

func (subscription *Subscription) Id() SubscriptionId {
	return subscription.id
}

func (subscription *Subscription) Kind() string {
	return subscription.kind
}

func (subscription *Subscription) SetKind(kind string) {
	subscription.kind = kind
}

func (subscription *Subscription) Save() (ok bool) {
	ok = subscription.partner.Subscriptions().Save(subscription)
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
		Reference:       reference,
		SubscribedUntil: subscription.Clock().Now().Add(1 * time.Minute),
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

	New() Subscription
	Find(id SubscriptionId) (Subscription, bool)
	FindAll() []Subscription
	FindOrCreateByKind(string) *Subscription
	Save(Subscription *Subscription) bool
	Delete(Subscription *Subscription) bool
	NewSubscription() *Subscription
}

func NewMemorySubscriptions(partner *Partner) *MemorySubscriptions {
	return &MemorySubscriptions{
		byIdentifier: make(map[SubscriptionId]*Subscription),
		partner:      partner,
	}
}

func (manager *MemorySubscriptions) New() Subscription {
	subscription := manager.NewSubscription()
	return *subscription
}

func (manager *MemorySubscriptions) FindOrCreateByKind(kind string) *Subscription {
	for _, subscription := range manager.byIdentifier {
		if subscription.Kind() == kind {
			return subscription
		}
	}

	subscription := manager.NewSubscription()
	subscription.SetKind(kind)
	return subscription
}

func (manager *MemorySubscriptions) Find(id SubscriptionId) (Subscription, bool) {
	subscription, ok := manager.byIdentifier[id]
	if ok {
		return *subscription, true
	} else {
		return Subscription{}, false
	}
}

func (manager *MemorySubscriptions) NewSubscription() *Subscription {
	sub := NewSubscription()
	manager.Save(sub)

	return sub
}

func (manager *MemorySubscriptions) FindAll() (subscriptions []Subscription) {
	if len(manager.byIdentifier) == 0 {
		return []Subscription{}
	}
	for _, subscription := range manager.byIdentifier {
		subscriptions = append(subscriptions, *subscription)
	}
	return
}

func (manager *MemorySubscriptions) Save(subscription *Subscription) bool {
	if subscription.Id() == "" {
		subscription.id = SubscriptionId(manager.NewUUID())
	}
	subscription.partner = manager.partner
	manager.byIdentifier[subscription.Id()] = subscription
	return true
}

func (manager *MemorySubscriptions) Delete(subscription *Subscription) bool {
	delete(manager.byIdentifier, subscription.Id())
	return true
}
