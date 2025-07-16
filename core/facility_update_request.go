package core

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
)

type FacilityUpdateRequest struct {
	facilityId model.FacilityId
	createdAt  time.Time
}

func NewFacilityUpdateRequest(facilityId model.FacilityId) *FacilityUpdateRequest {
	return &FacilityUpdateRequest{
		facilityId: facilityId,
		createdAt:  clock.DefaultClock().Now(),
	}
}

func (facilityUpdateRequest *FacilityUpdateRequest) FacilityId() model.FacilityId {
	return facilityUpdateRequest.facilityId
}

func (facilityUpdateRequest *FacilityUpdateRequest) CreatedAt() time.Time {
	return facilityUpdateRequest.createdAt
}
