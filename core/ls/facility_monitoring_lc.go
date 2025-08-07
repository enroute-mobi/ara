package ls

import "bitbucket.org/enroute-mobi/ara/model"

type FacilityMonitoringLastChange struct {
	lastState

	Status model.FacilityStatus
}

func (flc *FacilityMonitoringLastChange) UpdateState(f *model.Facility) {
	flc.Status = f.Status
}

func NewFacilityMonitoringLastChange(v *model.Facility, sub subscription) *FacilityMonitoringLastChange {
	flc := &FacilityMonitoringLastChange{}
	flc.SetSubscription(sub)
	flc.UpdateState(v)
	return flc
}

func (flc *FacilityMonitoringLastChange) HasChanged(f *model.Facility) bool {
	// Check Status
	if flc.Status != f.Status {
		return true
	}

	return false
}
