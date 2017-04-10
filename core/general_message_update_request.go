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

func NewSituationUpdateEvent(id SituationUpdateRequestId) *SituationUpdateRequest {
	return &SituationUpdateRequest{
		id:        id,
		CreatedAt: model.DefaultClock().Now(),
	}
}

func (situationUpdateRequest *SituationUpdateRequest) Id() SituationUpdateRequestId {
	return situationUpdateRequest.id
}

func (situationUpdateRequest *SituationUpdateRequest) CreatedAt() time.Time {
	return situationUpdateRequest.createdAt
}
