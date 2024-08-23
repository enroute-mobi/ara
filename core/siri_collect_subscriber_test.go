package core

import (
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/maps"
	"testing"
)

func Test_GetSubscriptionRequest(t *testing.T) {
	assert := assert.New(t)

	partners := createTestPartnerManager()
	partner := partners.New("slug")

	settings := map[string]string{
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	connector := NewSIRISubscriptionRequestDispatcher(partner)
	partner.connectors["test-startable-connector-connector"] = connector

	partner.subscriptionManager = NewMemorySubscriptions(partner)
	partners.Save(partner)

	partner.subscriptionManager.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	c, _ := partner.Connector("test-startable-connector-connector")
	subscriber := NewCollectSubcriber(c, "kind")
	assert.Len(subscriber.GetSubscriptionRequest(), 0, `No subscriptionRequest
without Subscription`)

	// Create a Subscription
	subscription := partner.Subscriptions().FindOrCreateByKind("kind")
	subscription.Save()

	subscriptionRequest := subscriber.GetSubscriptionRequest()
	assert.Empty(subscriptionRequest, "No subscriptionRequest with Subscription without Resource")

	// Create and add Resource to Subscription
	obj := model.NewCode("internal", "Value")
	reference := model.Reference{
		Code: &obj,
		Type: "type",
	}
	subscription.CreateAndAddNewResource(reference)
	subscriptionRequest = subscriber.GetSubscriptionRequest()

	assert.Len(subscriptionRequest, 1, "1 subscriptionRequest with Subscription with a Resource")
	assert.Equal(subscription.Id(), maps.Keys(subscriptionRequest)[0])

	modelsToRequest := subscriptionRequest[subscription.Id()].modelsToRequest
	assert.Len(modelsToRequest, 1)
	assert.Equal(obj, modelsToRequest[0].code)
	assert.Equal("type", modelsToRequest[0].kind)

	// Add another Resource to the Subscription
	obj2 := model.NewCode("internal", "AnotherValue")
	reference2 := model.Reference{
		Code: &obj2,
		Type: "type2",
	}
	subscription.CreateAndAddNewResource(reference2)
	subscriptionRequest = subscriber.GetSubscriptionRequest()

	assert.Len(subscriptionRequest, 1, "1 subscriptionRequest with Subscription with 2 Resources")
	assert.Equal(subscription.Id(), maps.Keys(subscriptionRequest)[0])

	modelsToRequest = subscriptionRequest[subscription.Id()].modelsToRequest
	assert.Len(modelsToRequest, 2, "2 Resources for the Subscription")
	var codes []model.Code
	var kinds []string
	for i := range modelsToRequest {
		codes = append(codes, modelsToRequest[i].code)
		kinds = append(kinds, modelsToRequest[i].kind)
	}
	assert.ElementsMatch(codes, []model.Code{obj, obj2})
	assert.ElementsMatch(kinds, []string{"type", "type2"})

	// Force Subscription 1st Resource RetryCount above 10
	resource := subscription.Resource(obj)
	resource.RetryCount += 11

	subscriptionRequest = subscriber.GetSubscriptionRequest()
	assert.Len(subscriptionRequest, 1, `1 subscriptionRequest with Subscription without 1 Resource
having 1 RetryCount > 10`)

	// Force Subscription 2nd Resource RetryCount above 10
	resource = subscription.Resource(obj2)
	resource.RetryCount += 11

	subscriptionRequest = subscriber.GetSubscriptionRequest()
	assert.Empty(subscriptionRequest)
}
