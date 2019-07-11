package model

type EventKind int

const (
	StopAreaEvent EventKind = iota
	LineEvent
)

type UpdateEvent interface {
	EventKind() EventKind
}

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
	return StopAreaEvent
}

type LineUpdateEvent struct {
	Origin string

	ObjectId ObjectID
	Name     string
}

func NewLineUpdateEvent() *LineUpdateEvent {
	return &LineUpdateEvent{}
}

func (ue *LineUpdateEvent) EventKind() EventKind {
	return LineEvent
}
