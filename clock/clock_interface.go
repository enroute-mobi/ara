package clock

import (
	"time"

	"github.com/jonboulle/clockwork"
)

type Clock interface {
	After(d time.Duration) <-chan time.Time
	Sleep(d time.Duration)
	Now() time.Time
	Since(t time.Time) time.Duration
}

type FakeClock interface {
	Clock

	Advance(d time.Duration)
	BlockUntil(n int)
}

var defaultClock = clockwork.NewRealClock()

func DefaultClock() clockwork.Clock {
	return defaultClock
}

func SetDefaultClock(clock clockwork.Clock) {
	defaultClock = clock
}

func NewRealClock() clockwork.Clock {
	return clockwork.NewRealClock()
}

func NewFakeClockAt(time time.Time) clockwork.FakeClock {
	return clockwork.NewFakeClockAt(time)
}

var FAKE_CLOCK_INITIAL_DATE = time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC)

func NewFakeClock() clockwork.FakeClock {
	// use a fixture that does not fulfill Time.IsZero()
	return NewFakeClockAt(FAKE_CLOCK_INITIAL_DATE)
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
