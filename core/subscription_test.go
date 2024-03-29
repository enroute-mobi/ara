package core

import (
	"testing"

	"bitbucket.org/enroute-mobi/ara/model"
	"github.com/stretchr/testify/assert"
)

func Test_Subscription_Id(t *testing.T) {
	subscription := Subscription{
		id: "6ba7b814-9dad-11d1-0-00c04fd430c8",
	}

	if subscription.Id() != "6ba7b814-9dad-11d1-0-00c04fd430c8" {
		t.Errorf("subscription.Id() returns wrong value, got: %s, required: %s", subscription.Id(), "6ba7b814-9dad-11d1-0-00c04fd430c8")
	}
}

func Test_subscription_MarshalJSON(t *testing.T) {
	subscriptions := NewMemorySubscriptions(NewPartner())

	subscription := subscriptions.New("salut")
	subscription.id = "6ba7b814-9dad-11d1-0-00c04fd430c8"
	subscription.CreateAndAddNewResource(*model.NewReference(model.NewCode("test", "value")))
	subscription.externalId = "externalId"

	expected := `{"SubscriptionRef":"6ba7b814-9dad-11d1-0-00c04fd430c8","ExternalId":"externalId","Kind":"salut","Resources":[{"Reference":{"Code":{"test":"value"}},"RetryCount":0,"SubscribedUntil":"1984-04-04T00:02:00Z","SubscribedAt":"0001-01-01T00:00:00Z"}]}`
	jsonBytes, err := subscription.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	jsonString := string(jsonBytes)
	if jsonString != expected {
		t.Errorf("subscription.MarshalJSON() returns wrong json:\n got: %s\n want: %s", jsonString, expected)
	}
}

func Test_MemorySubscription_New(t *testing.T) {
	subcriptions := NewMemorySubscriptions(NewPartner())

	subcription := subcriptions.New("kind")
	if subcription.Id() == "" {
		t.Errorf("New subcription identifier should be an empty string, got: %s", subcription.Id())
	}
}

func Test_MemorySubscriptions_Find_NotFound(t *testing.T) {
	subscriptions := NewMemorySubscriptions(NewPartner())
	_, ok := subscriptions.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if ok {
		t.Errorf("Find should return false when Subscription isn't found")
	}
}

func Test_MemorySubscriptions_Find(t *testing.T) {
	subscriptions := NewMemorySubscriptions(NewPartner())

	existingSubscription := subscriptions.New("kind")

	subscriptionId := existingSubscription.Id()

	subscription, ok := subscriptions.Find(subscriptionId)
	if !ok {
		t.Errorf("Find should return true when subscription is found")
	}
	if subscription.Id() != subscriptionId {
		t.Errorf("Find should return a subscription with the given Id")
	}
}

func Test_MemorySubscriptions_FindAll(t *testing.T) {
	subscriptions := NewMemorySubscriptions(NewPartner())

	for i := 0; i < 5; i++ {
		subscriptions.New("kind")
	}

	foundSubscriptions := subscriptions.FindAll()

	if len(foundSubscriptions) != 5 {
		t.Errorf("FindAll should return all subscriptions")
	}
}

func Test_MemorySubscriptions_Delete(t *testing.T) {
	subscriptions := NewMemorySubscriptions(NewPartner())
	existingSubscription := subscriptions.New("kind")

	subscriptions.Delete(existingSubscription)

	_, ok := subscriptions.Find(existingSubscription.Id())
	if ok {
		t.Errorf("Deleted subscription should not be findable")
	}
}

func Test_Subscription_byIdentifier(t *testing.T) {
	assert := assert.New(t)

	subscriptions := NewMemorySubscriptions(NewPartner())
	existingSubscription := subscriptions.New("kind")

	obj := model.NewCode("CodeSpace", "Value")
	reference := model.Reference{
		Code: &obj,
	}

	existingSubscription.CreateAndAddNewResource(reference)

	subs := subscriptions.FindByResourceId(obj.String(), "kind")

	assert.Len(subs, 1)
}

func Test_Subscriptions_byKindAndResourceId(t *testing.T) {
	assert := assert.New(t)

	subscriptions := NewMemorySubscriptions(NewPartner())
	existingSubscription := subscriptions.New("kind")

	assert.Len(subscriptions.byKindAndResourceId, 0)

	existingSubscription.Save()

	assert.Len(subscriptions.byKindAndResourceId, 0)

	obj := model.NewCode("CodeSpace", "Value")
	reference := model.Reference{
		Code: &obj,
	}

	existingSubscription.CreateAndAddNewResource(reference)

	assert.Len(subscriptions.byKindAndResourceId, 1)

	subs := subscriptions.FindByResourceId(obj.String(), "kind")
	assert.Len(subs, 1)

	existingSubscription.DeleteResource(obj.String())

	assert.Len(subscriptions.byKindAndResourceId, 1)

	subs = subscriptions.FindByResourceId(obj.String(), "kind")
	assert.Len(subs, 0)

	existingSubscription.Save()

	assert.Len(subscriptions.byKindAndResourceId, 0)

	subs = subscriptions.FindByResourceId(obj.String(), "kind")
	assert.Len(subs, 0)
}
