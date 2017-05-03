package model

import "testing"

func Test_TransactionalSubscription_Find_NotFound(t *testing.T) {
	model := NewMemoryModel()
	subscriptions := NewTransactionalSubscriptions(model)

	_, ok := subscriptions.Find("6ba7b814-9dad-11d1-0-00c04fd430c8")
	if ok {
		t.Errorf("Find should return false when subscriptions isn't found")
	}
}

func Test_TransactionalSubcription_Find_Model(t *testing.T) {
	model := NewMemoryModel()
	subscriptions := NewTransactionalSubscriptions(model)

	existingSubscription := model.Subscriptions().New()
	model.Subscriptions().Save(&existingSubscription)

	subscriptionId := existingSubscription.Id()

	subscription, ok := subscriptions.Find(subscriptionId)
	if !ok {
		t.Errorf("Find should return true when subscription is found")
	}
	if subscription.Id() != subscriptionId {
		t.Errorf("Find should return a subscription with the given Id")
	}
}

func Test_TransactionalSubscriptions_Find_Saved(t *testing.T) {
	model := NewMemoryModel()
	subscriptions := NewTransactionalSubscriptions(model)

	existingSubscription := subscriptions.New()
	subscriptions.Save(&existingSubscription)

	subscriptionId := existingSubscription.Id()

	subscription, ok := subscriptions.Find(subscriptionId)
	if !ok {
		t.Errorf("Find should return true when Subscription is found")
	}
	if subscription.Id() != subscriptionId {
		t.Errorf("Find should return a Subscription with the given Id")
	}
}

func Test_TransactionSubscriptions_FindAll(t *testing.T) {
	model := NewMemoryModel()
	subscriptions := NewTransactionalSubscriptions(model)

	for i := 0; i < 5; i++ {
		existingSubscription := subscriptions.New()
		subscriptions.Save(&existingSubscription)
	}

	foundSubscriptions := subscriptions.FindAll()

	if len(foundSubscriptions) != 5 {
		t.Errorf("FindAll should return all Subscriptions")
	}

}

func Test_TransactionalSubscriptions_Save(t *testing.T) {
	model := NewMemoryModel()
	subscriptions := NewTransactionalSubscriptions(model)

	subscription := subscriptions.New()

	if success := subscriptions.Save(&subscription); !success {
		t.Errorf("Save should return true")
	}
	if subscription.Id() == "" {
		t.Errorf("New subscription identifier shouldn't be an empty string")
	}
	if _, ok := model.Subscriptions().Find(subscription.Id()); ok {
		t.Errorf("subscription shouldn't be saved before commit")
	}
}

func Test_TransactionalSubscriptions_Delete(t *testing.T) {
	model := NewMemoryModel()
	subscriptions := NewTransactionalSubscriptions(model)

	existingSubscription := model.Subscriptions().New()
	model.Subscriptions().Save(&existingSubscription)

	subscriptions.Delete(&existingSubscription)

	_, ok := subscriptions.Find(existingSubscription.Id())
	if !ok {
		t.Errorf("Subscription should not be deleted before commit")
	}
}

func Test_TransactionalSubscriptions_Commit(t *testing.T) {
	model := NewMemoryModel()
	subscriptions := NewTransactionalSubscriptions(model)

	// Test Save
	subscription := subscriptions.New()
	subscriptions.Save(&subscription)

	// Test Delete
	existingSubscription := model.Subscriptions().New()
	model.Subscriptions().Save(&existingSubscription)
	subscriptions.Delete(&existingSubscription)

	subscriptions.Commit()

	if _, ok := model.Subscriptions().Find(subscription.Id()); !ok {
		t.Errorf("Subscription should be saved after commit")
	}
	if _, ok := subscriptions.Find(existingSubscription.Id()); ok {
		t.Errorf("Subscription should be deleted after commit")
	}
}

func Test_TransactionalSubscriptions_Rollback(t *testing.T) {
	model := NewMemoryModel()
	subscriptions := NewTransactionalSubscriptions(model)

	subscription := subscriptions.New()
	subscriptions.Save(&subscription)

	subscriptions.Rollback()
	subscriptions.Commit()

	if _, ok := model.Subscriptions().Find(subscription.Id()); ok {
		t.Errorf("Subscription should not be saved with a rollback")
	}
}
