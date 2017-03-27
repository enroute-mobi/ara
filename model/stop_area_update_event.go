package model

type StopAreaUpdateEvent struct {
	id                          string
	StopAreaId                  StopAreaId
	StopVisitUpdateEvents       []*StopVisitUpdateEvent
	StopVisitNotCollectedEvents []*StopVisitNotCollectedEvent
}

type StopVisitNotCollectedEvent struct {
	StopVisitObjectId ObjectID
}

func NewStopAreaUpdateEvent(id string, stopAreaId StopAreaId) *StopAreaUpdateEvent {
	return &StopAreaUpdateEvent{id: id, StopAreaId: stopAreaId}
}

func (event *StopAreaUpdateEvent) Id() string {
	return event.id
}
