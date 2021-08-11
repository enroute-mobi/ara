package core

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
)

type VehicleUpdateRequest struct {
	lineId    model.LineId
	createdAt time.Time
}

func NewVehicleUpdateRequest(lineId model.LineId) *VehicleUpdateRequest {
	return &VehicleUpdateRequest{
		lineId:    lineId,
		createdAt: clock.DefaultClock().Now(),
	}
}

func (vehicleUpdateRequest *VehicleUpdateRequest) LineId() model.LineId {
	return vehicleUpdateRequest.lineId
}

func (vehicleUpdateRequest *VehicleUpdateRequest) CreatedAt() time.Time {
	return vehicleUpdateRequest.createdAt
}
