package model

type StopAreaUpdateEvent struct {
	id                    string
	StopVisitUpdateEvents []*StopVisitUpdateEvent
}

func NewStopAreaUpdateEvent(id string) *StopAreaUpdateEvent {
	return &StopAreaUpdateEvent{id: id}
}

func (event *StopAreaUpdateEvent) Id() string {
	return event.id
}
