package core

import (
	"encoding/json"
	"testing"
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
	subscription := Subscription{
		id:   "6ba7b814-9dad-11d1-0-00c04fd430c8",
		kind: "salut",
	}
	expected := `{"Id":"6ba7b814-9dad-11d1-0-00c04fd430c8","Kind":"salut"}`
	jsonBytes, err := subscription.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	jsonString := string(jsonBytes)
	if jsonString != expected {
		t.Errorf("subscription.MarshalJSON() returns wrong json:\n got: %s\n want: %s", jsonString, expected)
	}
}

func Test_Subscription_UnmarshalJSON(t *testing.T) {
	text := `{
    "ResourcesByObjectID": {
		"une ressource" : {
			"Reference": {
					"ObjectId": {
						"kind":"value"
					},
					"Id": "une id",
					"Type": "un type"
				}
			}
		}
	}`

	subscription := Subscription{}
	err := json.Unmarshal([]byte(text), &subscription)
	if err != nil {
		t.Fatal(err)
	}

	ressources := subscription.ResourcesByObjectID()

	if len(ressources) != 1 {
		t.Errorf("ResourcesByObjectID  should have len == 1 after UnmarshalJSON()")
	}

	if ressources["une ressource"] == nil {
		t.Errorf("ResourcesByObjectID  should have a ressource 'une ressource' after UnmarshalJSON()")
	}
}

func Test_MemorySubscription_New(t *testing.T) {
	subcriptions := NewMemorySubscriptions(NewPartner())

	subcription := subcriptions.New()
	if subcription.Id() != "" {
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

	existingSubscription := subscriptions.New()
	subscriptions.Save(&existingSubscription)

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
		existingSubscription := subscriptions.New()
		subscriptions.Save(&existingSubscription)
	}

	foundSubscriptions := subscriptions.FindAll()

	if len(foundSubscriptions) != 5 {
		t.Errorf("FindAll should return all subscriptions")
	}
}

func Test_MemorySubscriptions_Delete(t *testing.T) {
	subscriptions := NewMemorySubscriptions(NewPartner())
	existingSubscription := subscriptions.New()
	subscriptions.Save(&existingSubscription)

	subscriptions.Delete(&existingSubscription)

	_, ok := subscriptions.Find(existingSubscription.Id())
	if ok {
		t.Errorf("Deleted subscription should not be findable")
	}
}
