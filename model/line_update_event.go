package model

type LineUpdateEvent struct {
	Origin string

	Code Code
	Name     string
}

func NewLineUpdateEvent() *LineUpdateEvent {
	return &LineUpdateEvent{}
}

func (ue *LineUpdateEvent) EventKind() EventKind {
	return LINE_EVENT
}
