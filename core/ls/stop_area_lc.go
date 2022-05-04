package ls

import (
	"bitbucket.org/enroute-mobi/ara/model"
)

type StopAreaLastChange struct {
	lastState

	origins *model.StopAreaOrigins
}

func NewStopAreaLastChange(sa *model.StopArea, sub subscription) *StopAreaLastChange {
	salc := &StopAreaLastChange{}
	salc.SetSubscription(sub)
	salc.UpdateState(sa)
	return salc
}

func (salc *StopAreaLastChange) UpdateState(stopArea *model.StopArea) bool {
	salc.origins = stopArea.Origins.Copy()

	return true
}

func (salc *StopAreaLastChange) Haschanged(stopArea *model.StopArea) ([]string, bool) {
	return salc.origins.PartnersLost(stopArea.Origins)
}
