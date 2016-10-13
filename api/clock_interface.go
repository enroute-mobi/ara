package api

import (
	"github.com/jonboulle/clockwork"
	"time"
)

type Clock interface {
	After(d time.Duration) <-chan time.Time
	Sleep(d time.Duration)
	Now() time.Time
	Since(t time.Time) time.Duration
}

var defaultClock = clockwork.NewRealClock()

func DefaultClock() clockwork.Clock {
	return defaultClock
}

func SetDefaultClock(clock clockwork.Clock) {
	defaultClock = clock
}

func NewFakeClockAt(time time.Time) clockwork.FakeClock {
	return clockwork.NewFakeClockAt(time)
}

func NewFakeClock() clockwork.FakeClock {
	// use a fixture that does not fulfill Time.IsZero()
	return NewFakeClockAt(time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC))
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
