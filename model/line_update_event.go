package model

type LineUpdateEvent struct {
	Origin string

	ObjectId ObjectID
	Name     string
}

func NewLineUpdateEvent() *LineUpdateEvent {
	return &LineUpdateEvent{}
}

func (ue *LineUpdateEvent) EventKind() EventKind {
	return LINE_EVENT
}
