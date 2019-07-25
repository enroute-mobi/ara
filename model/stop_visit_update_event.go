package model

type StopVisitUpdateEvent struct {
	Origin string

	ObjectId               ObjectID
	StopAreaObjectId       ObjectID
	VehicleJourneyObjectId ObjectID

	Monitored bool
	Schedules StopVisitSchedules
}

func NewStopVisitUpdateEvent() *StopVisitUpdateEvent {
	return &StopVisitUpdateEvent{
		Schedules: NewStopVisitSchedules(),
	}
}

func (ue *StopVisitUpdateEvent) EventKind() EventKind {
	return STOP_VISIT_EVENT
}
