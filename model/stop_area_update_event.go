package model

type StopAreaUpdateEvent struct {
	id                          string
	StopAreaId                  StopAreaId
	StopAreaMonitoredEvent      *StopAreaMonitoredEvent
	StopVisitUpdateEvents       []*StopVisitUpdateEvent
	StopVisitNotCollectedEvents []*StopVisitNotCollectedEvent
}

type StopAreaMonitoredEvent struct {
	Monitored bool
}

type StopVisitNotCollectedEvent struct {
	StopVisitObjectId ObjectID
}

func NewStopAreaUpdateEvent(id string, stopAreaId StopAreaId) *StopAreaUpdateEvent {
	return &StopAreaUpdateEvent{id: id, StopAreaId: stopAreaId}
}

func NewStopAreaMonitoredEvent(id string, stopAreaId StopAreaId, monitored bool) *StopAreaUpdateEvent {
	event := &StopAreaUpdateEvent{id: id, StopAreaId: stopAreaId}
	event.StopAreaMonitoredEvent = &StopAreaMonitoredEvent{Monitored: monitored}
	return event
}

func (event *StopAreaUpdateEvent) Id() string {
	return event.id
}
