package model

type FacilityUpdateEvent struct {
	Origin string

	Code   Code
	Status string
}

func NewFacilityUpdateEvent() *FacilityUpdateEvent {
	return &FacilityUpdateEvent{}
}

func (ue *FacilityUpdateEvent) EventKind() EventKind {
	return FACILITY_EVENT
}
