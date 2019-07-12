package model

type EventKind int

const (
	STOP_AREA_EVENT EventKind = iota
	LINE_EVENT
	VEHICLE_JOURNEY_EVENT
	STOP_VISIT_EVENT
)

type UpdateEvent interface {
	EventKind() EventKind
}
