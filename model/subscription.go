package model

import (
	"encoding/json"
	"time"
)

type SubscriptionId string

type Subscription struct {
	ClockConsumer

	id                  SubscriptionId
	kind                string
	model               Model
	resourcesByObjectID map[string]*SubscribedResource
}

type SubscribedResource struct {
	Reference       Reference
	SubscribedUntil time.Time
}

func NewSubscription(model Model) *Subscription {
	return &Subscription{
		model:               model,
		resourcesByObjectID: make(map[string]*SubscribedResource),
	}
}

func (subscription *Subscription) Id() SubscriptionId {
	return subscription.id
}

func (subscription *Subscription) Kind() string {
	return subscription.kind
}

func (subscription *Subscription) Save() (ok bool) {
	ok = subscription.model.Subscriptions().Save(subscription)
	return
}

func (subscription *Subscription) ResourcesByObjectID() map[string]*SubscribedResource {
	return subscription.resourcesByObjectID
}

func (subscription *Subscription) UnmarshalJSON(data []byte) error {
	type Alias Subscription

	aux := &struct {
		ResourcesByObjectID map[string]*SubscribedResource
		*Alias
	}{
		Alias: (*Alias)(subscription),
	}
	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if aux.ResourcesByObjectID != nil {
		subscription.resourcesByObjectID = aux.ResourcesByObjectID
	}
	return nil
}

func (subscription *Subscription) MarshalJSON() ([]byte, error) {
	type Alias Subscription
	aux := struct {
		Id                  SubscriptionId
		Kind                string                         `json:",omitempty"`
		ResourcesByObjectID map[string]*SubscribedResource `json:",omitempty"`
		*Alias
	}{
		Id:    subscription.id,
		Kind:  subscription.kind,
		Alias: (*Alias)(subscription),
	}

	if len(subscription.resourcesByObjectID) != 0 {
		aux.ResourcesByObjectID = subscription.resourcesByObjectID
	}
	return json.Marshal(&aux)
}

func (subscription *Subscription) Resource(obj ObjectID) *SubscribedResource {
	ObjIdString := obj.Value() + obj.Kind()
	return subscription.resourcesByObjectID[ObjIdString]
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

func (subscription *Subscription) NewResource(reference Reference) *SubscribedResource {
	ressource := SubscribedResource{
		Reference:       reference,
		SubscribedUntil: subscription.Clock().Now().Add(1 * time.Minute),
	}

	ObjIdString := reference.ObjectId.Value() + reference.ObjectId.Kind()
	subscription.resourcesByObjectID[ObjIdString] = &ressource
	return &ressource
}

type MemorySubscriptions struct {
	UUIDConsumer

	model *MemoryModel

	byIdentifier map[SubscriptionId]*Subscription
}

type Subscriptions interface {
	UUIDInterface

	New() Subscription
	Find(id SubscriptionId) (Subscription, bool)
	FindAll() []Subscription
	Save(Subscription *Subscription) bool
	Delete(Subscription *Subscription) bool
}

func NewMemorySubscriptions() *MemorySubscriptions {
	return &MemorySubscriptions{
		byIdentifier: make(map[SubscriptionId]*Subscription),
	}
}

func (manager *MemorySubscriptions) New() Subscription {
	subscription := NewSubscription(manager.model)
	return *subscription
}

func (manager *MemorySubscriptions) Find(id SubscriptionId) (Subscription, bool) {
	subscription, ok := manager.byIdentifier[id]
	if ok {
		return *subscription, true
	} else {
		return Subscription{}, false
	}
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
	subscription.model = manager.model
	manager.byIdentifier[subscription.Id()] = subscription
	return true
}

func (manager *MemorySubscriptions) Delete(subscription *Subscription) bool {
	delete(manager.byIdentifier, subscription.Id())
	return true
}
