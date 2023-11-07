package model

import "time"

type NotCollectedUpdateEvent struct {
	ObjectId       ObjectID
	NotCollectedAt time.Time
}

func NewNotCollectedUpdateEvent(obj ObjectID, t time.Time) *NotCollectedUpdateEvent {
	return &NotCollectedUpdateEvent{
		ObjectId:       obj,
		NotCollectedAt: t,
	}
}

func (ue *NotCollectedUpdateEvent) EventKind() EventKind {
	return NOT_COLLECTED_EVENT
}
