package core

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type SituationUpdateRequestId string

const (
	SITUATION_UPDATE_REQUEST_ALL       = "requestAll"
	SITUATION_UPDATE_REQUEST_LINE      = "requestLine"
	SITUATION_UPDATE_REQUEST_STOP_AREA = "requestStopArea"
)

type SituationUpdateRequest struct {
	id          SituationUpdateRequestId
	kind        string
	requestedId string
	createdAt   time.Time
}

func NewSituationUpdateRequest(kind, requestedId string) *SituationUpdateRequest {
	return &SituationUpdateRequest{
		id:          SituationUpdateRequestId(uuid.DefaultUUIDGenerator().NewUUID()),
		kind:        kind,
		requestedId: requestedId,
		createdAt:   clock.DefaultClock().Now(),
	}
}

func (situationUpdateRequest *SituationUpdateRequest) Id() SituationUpdateRequestId {
	return situationUpdateRequest.id
}

func (situationUpdateRequest *SituationUpdateRequest) Kind() string {
	return situationUpdateRequest.kind
}

func (situationUpdateRequest *SituationUpdateRequest) RequestedId() string {
	return situationUpdateRequest.requestedId
}

func (situationUpdateRequest *SituationUpdateRequest) CreatedAt() time.Time {
	return situationUpdateRequest.createdAt
}
