package ls

type LastState interface {
	SetSubscription(subscription)
}

type subscription interface {
	SubscriptionOption(string) string
}

type lastState struct {
	subscription subscription
}

func (ls *lastState) SetSubscription(sub subscription) {
	ls.subscription = sub
}
