package ls

import (
	"bitbucket.org/enroute-mobi/ara/model"
)

type ProductionTimetableLastChange struct {
	lastState
}

func NewProductionTimetableLastChange(sv *model.StopVisit, sub subscription) *ProductionTimetableLastChange {
	pttlc := &ProductionTimetableLastChange{}
	pttlc.SetSubscription(sub)
	pttlc.UpdateState(sv)
	return pttlc
}

func (pttlc *ProductionTimetableLastChange) UpdateState(sv *model.StopVisit) bool {
	return true
}

func (pttlc *ProductionTimetableLastChange) Haschanged(sv *model.StopVisit) bool {
	return false
}
