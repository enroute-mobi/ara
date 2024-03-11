package ls

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/model"
)

type SituationLastChange struct {
	lastState

	recordedAt time.Time
}

func NewSituationLastChange(situation *model.Situation, sub subscription) *SituationLastChange {
	gmlc := &SituationLastChange{}
	gmlc.SetSubscription(sub)
	gmlc.UpdateState(situation)
	return gmlc
}

func (slc *SituationLastChange) UpdateState(situation *model.Situation) bool {
	slc.recordedAt = situation.RecordedAt
	return true
}

func (slc *SituationLastChange) Haschanged(situation *model.Situation) bool {
	return !situation.RecordedAt.Equal(slc.recordedAt)
}
