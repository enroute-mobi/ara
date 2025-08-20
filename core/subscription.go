package core

import (
	"encoding/json"
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core/ls"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
)

const (
	EstimatedTimetableBroadcast  = "EstimatedTimetableBroadcast"
	GeneralMessageBroadcast      = "GeneralMessageBroadcast"
	ProductionTimetableBroadcast = "ProductionTimetableBroadcast"
	SituationExchangeBroadcast   = "SituationExchangeBroadcast"
	StopMonitoringBroadcast      = "StopMonitoringBroadcast"
	VehicleMonitoringBroadcast   = "VehicleMonitoringBroadcast"
	FacilityMonitoringBroadcast  = "FacilityMonitoringBroadcast"

	EstimatedTimetableCollect = "EstimatedTimetableCollect"
	GeneralMessageCollect     = "GeneralMessageCollect"
	SituationExchangeCollect  = "SituationExchangeCollect"
	StopMonitoringCollect     = "StopMonitoringCollect"
	VehicleMonitoringCollect  = "VehicleMonitoringCollect"
	FacilityMonitoringCollect = "FacilityMonitoringCollect"
)

type SubscriptionId string

type Subscription struct {
	sync.RWMutex
	clock.ClockConsumer

	manager Subscriptions

	id         SubscriptionId
	kind       string
	externalId string
	subscribed bool

	SubscriberRef string

	resourcesByCode     map[string]*SubscribedResource
	subscriptionOptions map[string]string
}

type APISubscription struct {
	ExternalId            string
	Kind                  string
	SubscriberRef         string
	SubscribeResourcesNow bool

	References []model.Reference
}

func (subscription *Subscription) SetDefinition(apisub *APISubscription) {
	subscription.SetExternalId(apisub.ExternalId)
	subscription.kind = apisub.Kind
	subscription.SubscriberRef = apisub.SubscriberRef
	for _, ref := range apisub.References {
		if ref.Code != nil {
			subscription.CreateAndAddNewResource(ref)
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

func (subscription *Subscription) ResourcesByCodeCopy() map[string]*SubscribedResource {
	m := make(map[string]*SubscribedResource)
	subscription.RLock()
	for k, v := range subscription.resourcesByCode {
		m[k] = v
	}
	subscription.RUnlock()
	return m
}

func (subscription *Subscription) MarshalJSON() ([]byte, error) {
	resources := make([]*SubscribedResource, 0)

	subscription.RLock()
	for _, resource := range subscription.resourcesByCode {
		resources = append(resources, resource)
	}
	subscription.RUnlock()

	aux := struct {
		Id                  SubscriptionId        `json:"SubscriptionRef,omitempty"`
		ExternalId          string                `json:"ExternalId,omitempty"`
		Kind                string                `json:",omitempty"`
		Resources           []*SubscribedResource `json:",omitempty"`
		SubscriptionOptions map[string]string     `json:",omitempty"`
	}{
		Id:                  subscription.id,
		ExternalId:          subscription.externalId,
		Kind:                subscription.kind,
		Resources:           resources,
		SubscriptionOptions: subscription.subscriptionOptions,
	}
	return json.Marshal(&aux)
}

func (subscription *Subscription) Resource(obj model.Code) *SubscribedResource {
	subscription.RLock()
	sub, present := subscription.resourcesByCode[obj.String()]
	subscription.RUnlock()
	if !present {
		return nil
	}
	return sub
}

func (subscription *Subscription) UniqueResource() (r *SubscribedResource) {
	subscription.RLock()
	if len(subscription.resourcesByCode) != 1 {
		subscription.RUnlock()
		return
	}

	for _, ressource := range subscription.resourcesByCode {
		r = ressource
	}

	subscription.RUnlock()
	return
}

func (subscription *Subscription) Resources(now time.Time) (ressources []*SubscribedResource) {
	subscription.RLock()
	for _, ressource := range subscription.resourcesByCode {
		if ressource.SubscribedUntil.After(subscription.Clock().Now()) {
			ressources = append(ressources, ressource)
		}
	}
	subscription.RUnlock()
	return ressources
}

func (subscription *Subscription) AddNewResource(resource *SubscribedResource) {
	resource.subscription = subscription
	if !resource.subscribedAt.IsZero() {
		subscription.subscribed = true
	}
	subscription.Lock()
	subscription.resourcesByCode[resource.Reference.Code.String()] = resource
	subscription.Unlock()
}

func (subscription *Subscription) CreateAndAddNewResource(reference model.Reference) *SubscribedResource {
	logger.Log.Debugf("Create subscribed resource for %v", reference.Code.String())

	resource := SubscribedResource{
		subscription:     subscription,
		Reference:        reference,
		SubscribedUntil:  subscription.Clock().Now().Add(2 * time.Minute),
		lastStates:       make(map[string]ls.LastState),
		resourcesOptions: make(map[string]string),
	}
	subscription.Lock()
	subscription.resourcesByCode[reference.Code.String()] = &resource
	subscription.Unlock()

	subscription.manager.Index(subscription)

	return &resource
}

func (subscription *Subscription) DeleteResource(key string) {
	subscription.Lock()
	delete(subscription.resourcesByCode, key)
	subscription.Unlock()
}

func (subscription *Subscription) ResourcesLen() (i int) {
	subscription.RLock()
	i = len(subscription.resourcesByCode)
	subscription.RUnlock()
	return
}
