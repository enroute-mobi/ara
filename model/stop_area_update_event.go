package model

type StopAreaUpdateEvent struct {
	Origin string

	Code        Code
	ParentCode  Code
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
