package core

import (
	"time"

	"github.com/af83/edwig/model"
)

type SituationUpdateRequestId string

type SituationUpdateRequest struct {
	id        SituationUpdateRequestId
	lineId    model.LineId
	createdAt time.Time
}

func NewSituationUpdateRequest(lineId model.LineId) *SituationUpdateRequest {
	return &SituationUpdateRequest{
		id:        SituationUpdateRequestId(model.DefaultUUIDGenerator().NewUUID()),
		lineId:    lineId,
		createdAt: model.DefaultClock().Now(),
	}
}

func (situationUpdateRequest *SituationUpdateRequest) Id() SituationUpdateRequestId {
	return situationUpdateRequest.id
}

func (situationUpdateRequest *SituationUpdateRequest) LineId() model.LineId {
	return situationUpdateRequest.lineId
}

func (situationUpdateRequest *SituationUpdateRequest) CreatedAt() time.Time {
	return situationUpdateRequest.createdAt
}
