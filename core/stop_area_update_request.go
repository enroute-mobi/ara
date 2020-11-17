package core

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type StopAreaUpdateRequestId string

type StopAreaUpdateRequest struct {
	id         StopAreaUpdateRequestId
	stopAreaId model.StopAreaId
	createdAt  time.Time
}

func NewStopAreaUpdateRequest(stopAreaId model.StopAreaId) *StopAreaUpdateRequest {
	return &StopAreaUpdateRequest{
		id:         StopAreaUpdateRequestId(uuid.DefaultUUIDGenerator().NewUUID()),
		stopAreaId: stopAreaId,
		createdAt:  clock.DefaultClock().Now(),
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
