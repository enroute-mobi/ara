package model

type LegacyStopAreaUpdateEvent struct {
	id                          string
	Origin                      string
	StopAreaId                  StopAreaId
	StopAreaAttributes          StopAreaAttributes
	StopAreaMonitoredEvent      *StopAreaMonitoredEvent
	LegacyStopVisitUpdateEvents []*LegacyStopVisitUpdateEvent
	StopVisitNotCollectedEvents []*StopVisitNotCollectedEvent
}

type StopAreaMonitoredEvent struct {
	Partner string
	Status  bool
}

type StopVisitNotCollectedEvent struct {
	StopVisitObjectId ObjectID
}

func NewLegacyStopAreaUpdateEvent(id string, stopAreaId StopAreaId) *LegacyStopAreaUpdateEvent {
	return &LegacyStopAreaUpdateEvent{id: id, StopAreaId: stopAreaId}
}

func (update *LegacyStopAreaUpdateEvent) SetId(id string) {
	update.id = id
}

func NewStopAreaMonitoredEvent(id string, stopAreaId StopAreaId, partner string, status bool) *LegacyStopAreaUpdateEvent {
	event := &LegacyStopAreaUpdateEvent{id: id, StopAreaId: stopAreaId}
	event.StopAreaMonitoredEvent = &StopAreaMonitoredEvent{
		Partner: partner,
		Status:  status,
	}
	return event
}

func (event *LegacyStopAreaUpdateEvent) Id() string {
	return event.id
}
