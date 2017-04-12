package core

import (
	"time"

	"github.com/af83/edwig/model"
)

type SituationUpdateRequestId string

type SituationUpdateRequest struct {
	id        SituationUpdateRequestId
	createdAt time.Time
}

func NewSituationUpdateRequest(id SituationUpdateRequestId) *SituationUpdateRequest {
	return &SituationUpdateRequest{
		id:        id,
		createdAt: model.DefaultClock().Now(),
	}
}

func (situationUpdateRequest *SituationUpdateRequest) Id() SituationUpdateRequestId {
	return situationUpdateRequest.id
}

func (situationUpdateRequest *SituationUpdateRequest) CreatedAt() time.Time {
	return situationUpdateRequest.createdAt
}
