package api

import "github.com/jonboulle/clockwork"

var defaultClock = clockwork.NewRealClock()

func DefaultClock() clockwork.Clock {
	return defaultClock
}

func SetDefaultClock(clock clockwork.Clock) {
	defaultClock = clock
}

type ClockConsumer struct {
	clock clockwork.Clock
}

func (consumer *ClockConsumer) SetClock(clock clockwork.Clock) {
	consumer.clock = clock
}

func (consumer *ClockConsumer) Clock() clockwork.Clock {
	if consumer.clock == nil {
		consumer.clock = DefaultClock()
	}
	return consumer.clock
}
