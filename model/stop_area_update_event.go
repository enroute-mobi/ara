package model

import "github.com/af83/edwig/siri"

type StopAreaUpdateEvent struct {
	UUIDConsumer

	id                    string
	stopVisitUpdateEvents []*StopVisitUpdateEvent
}

func NewStopAreaUpdateEvent(response *siri.XMLStopMonitoringResponse) *StopAreaUpdateEvent {
	event := &StopAreaUpdateEvent{}
	event.id = event.NewUUID()
	return event
}

func (event *StopAreaUpdateEvent) Id() string {
	return event.id
}

func (event *StopAreaUpdateEvent) StopVisitUpdateEvents() []*StopVisitUpdateEvent {
	return event.stopVisitUpdateEvents
}
