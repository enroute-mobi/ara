package ls

import (
	"bitbucket.org/enroute-mobi/ara/model"
)

type ProductionTimeTableLastChange struct {
	lastState
}

func NewProductionTimeTableLastChange(sv *model.StopVisit, sub subscription) *ProductionTimeTableLastChange {
	pttlc := &ProductionTimeTableLastChange{}
	pttlc.SetSubscription(sub)
	pttlc.UpdateState(sv)
	return pttlc
}

func (pttlc *ProductionTimeTableLastChange) UpdateState(sv *model.StopVisit) bool {
	return true
}

func (pttlc *ProductionTimeTableLastChange) Haschanged(sv *model.StopVisit) bool {
	return false
}
