package core

import (
	"time"

	"github.com/af83/edwig/model"
)

type StopAreaUpdateRequestId string

type StopAreaUpdateRequest struct {
	id         StopAreaUpdateRequestId
	stopAreaId model.StopAreaId
	createdAt  time.Time
}

func NewStopAreaUpdateRequest(givenStopAreaId model.StopAreaId) *StopAreaUpdateRequest {
	return &StopAreaUpdateRequest{
		id:         StopAreaUpdateRequestId(model.DefaultUUIDGenerator().NewUUID()),
		stopAreaId: givenStopAreaId,
		createdAt:  model.DefaultClock().Now(),
	}
}

func (stopAreaUpdateRequest *StopAreaUpdateRequest) Id() StopAreaUpdateRequestId {
	return stopAreaUpdateRequest.id
}

func (stopAreaUpdateRequest *StopAreaUpdateRequest) StopAreaId() model.StopAreaId {
	return stopAreaUpdateRequest.stopAreaId
}

func (stopAreaUpdateRequest *StopAreaUpdateRequest) CreatedAt() time.Time {
	return stopAreaUpdateRequest.createdAt
}
