package model

type EventKind int

const (
	STOP_AREA_EVENT EventKind = iota
	STATUS_EVENT
	LINE_EVENT
	VEHICLE_JOURNEY_EVENT
	STOP_VISIT_EVENT
	NOT_COLLECTED_EVENT
	VEHICLE_EVENT
	SITUATION_EVENT
	FACILITY_EVENT
)

type UpdateEvent interface {
	EventKind() EventKind
}
