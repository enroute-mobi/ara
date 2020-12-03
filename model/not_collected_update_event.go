package model

type NotCollectedUpdateEvent struct {
	ObjectId ObjectID
}

func NewNotCollectedUpdateEvent(obj ObjectID) *NotCollectedUpdateEvent {
	return &NotCollectedUpdateEvent{
		ObjectId: obj,
	}
}

func (ue *NotCollectedUpdateEvent) EventKind() EventKind {
	return NOT_COLLECTED_EVENT
}
