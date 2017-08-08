package core

import (
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type lastState interface {
}

type stopMonitoringLastChange struct {
	lastChange siri.SIRIMonitoredStopVisit
}

func (smlc *stopMonitoringLastChange) Haschanged(stopVisit model.StopVisit) bool {

	return true
}

func (smlc *stopMonitoringLastChange) UpdateState(sm siri.SIRIMonitoredStopVisit) bool {
	smlc.lastChange = sm
	return true
}

type generalMessageLastChange struct {
	lastChange siri.SIRIGeneralMessage
}

func (smlc *generalMessageLastChange) Haschanged(situation model.Situation) bool {

	return true
}

func (smlc *generalMessageLastChange) UpdateState(sm siri.SIRIGeneralMessage) bool {
	smlc.lastChange = sm
	return true
}
