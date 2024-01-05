package model

import "time"

type NotCollectedUpdateEvent struct {
	Code       Code
	NotCollectedAt time.Time
}

func NewNotCollectedUpdateEvent(obj Code, t time.Time) *NotCollectedUpdateEvent {
	return &NotCollectedUpdateEvent{
		Code:       obj,
		NotCollectedAt: t,
	}
}

func (ue *NotCollectedUpdateEvent) EventKind() EventKind {
	return NOT_COLLECTED_EVENT
}
