package ls

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/model"
)

type GeneralMessageLastChange struct {
	lastState

	recordedAt time.Time
}

func NewGeneralMessageLastChange(situation *model.Situation, sub subscription) *GeneralMessageLastChange {
	gmlc := &GeneralMessageLastChange{}
	gmlc.SetSubscription(sub)
	gmlc.UpdateState(situation)
	return gmlc
}

func (gmlc *GeneralMessageLastChange) UpdateState(situation *model.Situation) bool {
	gmlc.recordedAt = situation.RecordedAt
	return true
}

func (gmlc *GeneralMessageLastChange) Haschanged(situation *model.Situation) bool {
	return !situation.RecordedAt.Equal(gmlc.recordedAt)
}
