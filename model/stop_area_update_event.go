package model

type StopAreaUpdateEvent struct {
	Origin string

	ObjectId        ObjectID
	ParentObjectId  ObjectID
	Name            string
	CollectedAlways bool
	Longitude       float64
	Latitude        float64
}

func NewStopAreaUpdateEvent() *StopAreaUpdateEvent {
	return &StopAreaUpdateEvent{}
}

func (ue *StopAreaUpdateEvent) EventKind() EventKind {
	return STOP_AREA_EVENT
}
