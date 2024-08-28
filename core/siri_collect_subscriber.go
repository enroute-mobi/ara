package core

import (
	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type CollectSubscriber struct {
	connector Connector
	name      string
}

func NewCollectSubcriber(c Connector, name string) *CollectSubscriber {
	return &CollectSubscriber{
		connector: c,
		name:      name,
	}
}

type subscriptionRequest struct {
	requestMessageRef string
	modelsToRequest   []*modelToRequest
}

type modelToRequest struct {
	code model.Code
	kind string
}

func (cs *CollectSubscriber) GetSubscriptionRequest() map[SubscriptionId]*subscriptionRequest {
	subscriptionRequests := make(map[SubscriptionId]*subscriptionRequest)
	subscriptions := cs.connector.Partner().Subscriptions().FindSubscriptionsByKind(cs.name)
	if len(subscriptions) == 0 {
		logger.Log.Debugf("%v"+"Subscriber visit without "+"%v "+"subscriptions", cs.name, cs.name)
		return nil
	}

	for _, subscription := range subscriptions {
		modelsToRequest := []*modelToRequest{}
		for _, resource := range subscription.ResourcesByCodeCopy() {
			if resource.SubscribedAt().IsZero() && resource.RetryCount <= 10 {
				toRequest := &modelToRequest{
					kind: resource.Reference.Type,
					code: *resource.Reference.Code,
				}
				modelsToRequest = append(modelsToRequest, toRequest)
			}
		}

		if len(modelsToRequest) == 0 {
			continue
		}

		_, ok := subscriptionRequests[subscription.Id()]
		if !ok {
			subscriptionRequests[subscription.Id()] = &subscriptionRequest{
				requestMessageRef: cs.connector.Partner().NewMessageIdentifier(),
			}
		}

		subscriptionRequests[subscription.Id()].modelsToRequest = append(subscriptionRequests[subscription.Id()].modelsToRequest, modelsToRequest...)
	}

	return subscriptionRequests
}

func (cs *CollectSubscriber) HandleResponse(subscriptionRequests map[SubscriptionId]*subscriptionRequest, message *audit.BigQueryMessage, response *sxml.XMLSubscriptionResponse) {
	for _, responseStatus := range response.ResponseStatus() {
		var subscriptionRequest *subscriptionRequest
		var ok bool

		// Find the subscriptionRef
		subscriptionRequest, ok = subscriptionRequests[SubscriptionId(responseStatus.SubscriptionRef())]
		if !ok {
			logger.Log.Debugf("%v"+"Subscriber ResponseStatus: SubscriptionRef %v not requested", cs.name, responseStatus.SubscriptionRef())
			continue
		}

		// Verify RequestMessageRef and skip if only 1 request
		if len(subscriptionRequests) != 1 {
			if subscriptionRequest.requestMessageRef != responseStatus.RequestMessageRef() {
				logger.Log.Debugf("%v"+"Subscriber ResponseStatus: RequestMessageRef unknown: %v", cs.name, responseStatus.RequestMessageRef())
				continue
			}
		}

		// Find the subscription
		subscription, ok := cs.connector.Partner().Subscriptions().Find(SubscriptionId(responseStatus.SubscriptionRef()))
		if !ok { // Should never happen
			logger.Log.Debugf("%v"+"Subscriber Response for unknown subscription %v", cs.name, responseStatus.SubscriptionRef())
			continue
		}

		// Find the models
		if len(subscriptionRequest.modelsToRequest) == 0 {
			if !ok { // Should never happen
				logger.Log.Debugf("%v"+"Subscriber: Error, no models to request for subscription %v", cs.name, subscription.Id())
				continue
			}
		}

		for _, modelToRequest := range subscriptionRequest.modelsToRequest {
			modelValue := modelToRequest.code.Value()
			var resource *SubscribedResource

			switch modelToRequest.code.String() {
			case "SituationExchangeCollect:all", "GeneralMessageCollect:all":
				resource = subscription.Resource(modelToRequest.code)
			default:
				resource = subscription.Resource(model.NewCode(cs.connector.RemoteCodeSpace(), modelValue))
			}

			if resource == nil { // Should never happen
				logger.Log.Debugf("%v"+"Subscriber Response for unknown subscription resource %v", cs.name, modelToRequest.code.String())
				continue
			}

			if !responseStatus.Status() {
				logger.Log.Debugf("%v"+"Subscriber Subscription status false for %v %v: %v %v ",
					cs.name,
					modelToRequest.kind,
					modelValue,
					responseStatus.ErrorType(),
					responseStatus.ErrorText())
				resource.RetryCount++
				message.Status = "Error"
				continue
			}
			resource.Subscribed(cs.connector.Clock().Now())
			resource.RetryCount = 0
		}
		delete(subscriptionRequests, subscription.Id())
	}
}

func (cs *CollectSubscriber) IncrementRetryCountFromMap(subscriptionRequests map[SubscriptionId]*subscriptionRequest) {
	for subId, subscriptionRequest := range subscriptionRequests {
		subscription, ok := cs.connector.Partner().Subscriptions().Find(subId)
		if !ok { // Should never happen
			continue
		}
		for _, l := range subscriptionRequest.modelsToRequest {
			resource := subscription.Resource(model.NewCode(cs.connector.RemoteCodeSpace(), l.code.Value()))
			if resource == nil { // Should never happen
				continue
			}
			resource.RetryCount++
		}
	}
}
