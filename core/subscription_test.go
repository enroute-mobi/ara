package core

import (
	"testing"

	"github.com/af83/edwig/model"
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
	subscription := &Subscription{
		resourcesByObjectID: make(map[string]*SubscribedResource),
		subscriptionOptions: make(map[string]string),
	}
	subscription.id = "6ba7b814-9dad-11d1-0-00c04fd430c8"
	subscription.kind = "salut"
	subscription.CreateAddNewResource(*model.NewReference(model.NewObjectID("test", "value")))

	expected := `{"SubscriptionRef":"6ba7b814-9dad-11d1-0-00c04fd430c8","Kind":"salut","Resources":[{"Reference":{"ObjectId":{"test":"value"}},"RetryCount":0,"SubscribedAt":"0001-01-01T00:00:00Z","SubscribedUntil":"1984-04-04T00:01:00Z"}]}`
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
	subscriptions := NewMemorySubscriptions(NewPartner())
	existingSubscription := subscriptions.New("kind")

	obj := model.NewObjectID("Kind", "Value")
	reference := model.Reference{
		ObjectId: &obj,
	}

	existingSubscription.CreateAddNewResource(reference)

	_, ok := subscriptions.FindByRessourceId(obj.String())

	if !ok {
		t.Errorf("Should have found the subscription")
	}
}
