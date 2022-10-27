package core

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
)

type StopAreaUpdateRequest struct {
	stopAreaId model.StopAreaId
	createdAt  time.Time
}

func NewStopAreaUpdateRequest(stopAreaId model.StopAreaId) *StopAreaUpdateRequest {
	return &StopAreaUpdateRequest{
		stopAreaId: stopAreaId,
		createdAt:  clock.DefaultClock().Now(),
	}
}

func (stopAreaUpdateRequest *StopAreaUpdateRequest) StopAreaId() model.StopAreaId {
	return stopAreaUpdateRequest.stopAreaId
}

func (stopAreaUpdateRequest *StopAreaUpdateRequest) CreatedAt() time.Time {
	return stopAreaUpdateRequest.createdAt
}
