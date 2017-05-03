package model

type TransactionalSubscriptions struct {
	UUIDConsumer

	model   Model
	saved   map[SubscriptionId]*Subscription
	deleted map[SubscriptionId]*Subscription
}

func NewTransactionalSubscriptions(model Model) *TransactionalSubscriptions {
	subscriptions := TransactionalSubscriptions{model: model}
	subscriptions.resetCaches()
	return &subscriptions
}

func (manager *TransactionalSubscriptions) resetCaches() {
	manager.saved = make(map[SubscriptionId]*Subscription)
	manager.deleted = make(map[SubscriptionId]*Subscription)
}

func (manager *TransactionalSubscriptions) New() Subscription {
	return *NewSubscription(manager.model)
}

func (manager *TransactionalSubscriptions) Find(id SubscriptionId) (Subscription, bool) {
	subscription, ok := manager.saved[id]
	if ok {
		return *subscription, ok
	}
	return manager.model.Subscriptions().Find(id)
}

func (manager *TransactionalSubscriptions) FindAll() []Subscription {
	subscriptions := []Subscription{}
	for _, subscription := range manager.saved {
		subscriptions = append(subscriptions, *subscription)
	}
	savedSubscriptions := manager.model.Subscriptions().FindAll()
	for _, subscription := range savedSubscriptions {
		_, ok := manager.saved[subscription.Id()]
		if !ok {
			subscriptions = append(subscriptions, subscription)
		}
	}
	return subscriptions
}

func (manager *TransactionalSubscriptions) Save(subscription *Subscription) bool {
	if subscription.Id() == "" {
		subscription.id = SubscriptionId(manager.NewUUID())
	}
	manager.saved[subscription.Id()] = subscription
	return true
}

func (manager *TransactionalSubscriptions) Delete(subscription *Subscription) bool {
	manager.deleted[subscription.Id()] = subscription
	return true
}

func (manager *TransactionalSubscriptions) Commit() error {
	for _, subscription := range manager.deleted {
		manager.model.Subscriptions().Delete(subscription)
	}
	for _, subscription := range manager.saved {
		manager.model.Subscriptions().Save(subscription)
	}
	return nil
}

func (manager *TransactionalSubscriptions) Rollback() error {
	manager.resetCaches()
	return nil
}
