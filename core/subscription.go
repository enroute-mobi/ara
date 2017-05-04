package core

import (
	"encoding/json"
	"time"

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

type SubscribedResource struct {
	Reference       model.Reference
	SubscribedUntil time.Time
}

func NewSubscription(partner *Partner) *Subscription {
	return &Subscription{
		partner:             partner,
		resourcesByObjectID: make(map[string]*SubscribedResource),
	}
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

func (subscription *Subscription) Resource(obj model.ObjectID) *SubscribedResource {
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

func (subscription *Subscription) NewResource(reference model.Reference) *SubscribedResource {
	ressource := SubscribedResource{
		Reference:       reference,
		SubscribedUntil: subscription.Clock().Now().Add(1 * time.Minute),
	}

	ObjIdString := reference.ObjectId.Value() + reference.ObjectId.Kind()
	subscription.resourcesByObjectID[ObjIdString] = &ressource
	return &ressource
}

type MemorySubscriptions struct {
	model.UUIDConsumer

	partner *Partner

	byIdentifier map[SubscriptionId]*Subscription
}

type Subscriptions interface {
	model.UUIDInterface

	New() Subscription
	Find(id SubscriptionId) (Subscription, bool)
	FindAll() []Subscription
	Save(Subscription *Subscription) bool
	Delete(Subscription *Subscription) bool
}

func NewMemorySubscriptions(partner *Partner) *MemorySubscriptions {
	return &MemorySubscriptions{
		byIdentifier: make(map[SubscriptionId]*Subscription),
		partner:      partner,
	}
}

func (manager *MemorySubscriptions) New() Subscription {
	subscription := NewSubscription(manager.partner)
	return *subscription
}

func (manager *MemorySubscriptions) FindOrCreateByKind(kind string) *Subscription {
	for _, subscription := range manager.byIdentifier {
		if subscription.Kind() == kind {
			return subscription
		}
	}

	subscription := NewSubscription(manager.partner)
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
