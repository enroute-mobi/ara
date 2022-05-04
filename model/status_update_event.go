package model

type StatusUpdateEvent struct {
	StopAreaId StopAreaId
	Partner    string
	Status     bool
}

func NewStatusUpdateEvent(stopAreaId StopAreaId, partner string, status bool) *StatusUpdateEvent {
	return &StatusUpdateEvent{
		StopAreaId: stopAreaId,
		Partner:    partner,
		Status:     status,
	}
}

func (ue *StatusUpdateEvent) EventKind() EventKind {
	return STATUS_EVENT
}
