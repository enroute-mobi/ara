package model

import (
	"slices"
)

type MacroId string
type contextAttributes interface{}
type updaterAttributes interface{}
type context func(ModelInstance) bool
type updater func(ModelInstance) error
type contexFactory func(contextAttributes) (context, error)
type updaterFactory func(updaterAttributes) (updater, error)

const (
	IfAttribute               = "IfAttribute"
	SetAttribute              = "SetAttribute"
	DefineAimedScheduledTimes = "DefineAimedScheduledTimes"
	DefineSituationAffects    = "DefineSituationAffects"
)

var contexes = []string{IfAttribute}
var updaters = []string{SetAttribute, DefineAimedScheduledTimes}

type Macro struct {
	c context
	u []updater
}

func NewMacro() *Macro {
	return &Macro{}
}

func NewMacroWithContext(c context) *Macro {
	return &Macro{
		c: c,
	}
}

func (m *Macro) AddUpdater(u updater) {
	m.u = append(m.u, u)
}

func (m *Macro) AddContext(c context) {
	if m.c == nil {
		m.c = c
		return
	}
	m.c = func(mi ModelInstance) bool {
		if !m.c(mi) {
			return false
		}
		return c(mi)
	}
}

func (m *Macro) Update(mi ModelInstance) (ok bool, err error) {
	if m.c != nil && !m.c(mi) {
		return false, nil
	}
	for i := range m.u {
		err := m.u[i](mi)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func NewContexFromDatabase(sm *SelectMacro) (context, error) {
	if sm.ModelType.String == "VehicleJourney" && sm.Type == IfAttribute {
		return NewVehicleJourneyIfAttributeContext(sm)
	}
	return nil, nil
}

func NewUpdaterFromDatabase(sm *SelectMacro) (updater, error) {
	if sm.ModelType.String == "VehicleJourney" && sm.Type == SetAttribute {
		return NewVehicleJourneySetAttributeUpdater(sm)
	} else if sm.ModelType.String == "StopVisit" && sm.Type == DefineAimedScheduledTimes {
		return NewStopVisitDefineAimedScheduledTimesUpdater(sm)
	} else if sm.ModelType.String == "Situation" && sm.Type == DefineSituationAffects {
		return NewDefineSituationAffectsUpdater(sm)
	}
	return nil, nil
}

func IsContext(c string) bool {
	return slices.Contains(contexes, c)
}

func IsUpdater(u string) bool {
	return slices.Contains(updaters, u)
}
