package model

type StopAreaUpdateEvent struct {
	id                          string
	Origin                      string
	StopAreaId                  StopAreaId
	StopAreaAttributes          StopAreaAttributes
	StopAreaMonitoredEvent      *StopAreaMonitoredEvent
	StopVisitUpdateEvents       []*StopVisitUpdateEvent
	StopVisitNotCollectedEvents []*StopVisitNotCollectedEvent
}

type StopAreaMonitoredEvent struct {
	Partner string
	Status  bool
}

type StopVisitNotCollectedEvent struct {
	StopVisitObjectId ObjectID
}

func NewStopAreaUpdateEvent(id string, stopAreaId StopAreaId) *StopAreaUpdateEvent {
	return &StopAreaUpdateEvent{id: id, StopAreaId: stopAreaId}
}

func (update *StopAreaUpdateEvent) SetId(id string) {
	update.id = id
}

func NewStopAreaMonitoredEvent(id string, stopAreaId StopAreaId, partner string, status bool) *StopAreaUpdateEvent {
	event := &StopAreaUpdateEvent{id: id, StopAreaId: stopAreaId}
	event.StopAreaMonitoredEvent = &StopAreaMonitoredEvent{
		Partner: partner,
		Status:  status,
	}
	return event
}

func (event *StopAreaUpdateEvent) Id() string {
	return event.id
}
