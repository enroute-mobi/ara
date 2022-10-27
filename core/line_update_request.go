package core

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
)

type LineUpdateRequest struct {
	lineId    model.LineId
	createdAt time.Time
}

func NewLineUpdateRequest(lineId model.LineId) *LineUpdateRequest {
	return &LineUpdateRequest{
		lineId:    lineId,
		createdAt: clock.DefaultClock().Now(),
	}
}

func (vehicleUpdateRequest *LineUpdateRequest) LineId() model.LineId {
	return vehicleUpdateRequest.lineId
}

func (vehicleUpdateRequest *LineUpdateRequest) CreatedAt() time.Time {
	return vehicleUpdateRequest.createdAt
}
