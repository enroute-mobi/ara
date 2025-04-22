package model

import (
	"slices"
)

type MacroId string
type updaterAttributes interface{}
type updater func(ModelInstance) error
type updaterFactory func(updaterAttributes) (updater, error)

const (
	SetAttribute              = "SetAttribute"
	DefineAimedScheduledTimes = "DefineAimedScheduledTimes"
	DefineSituationAffects    = "DefineSituationAffects"
	CreateCode                = "CreateCode"
)

var updaters = []string{SetAttribute, DefineAimedScheduledTimes}

type Macro struct {
	c Context
	u []updater
}

func NewMacro() *Macro {
	return &Macro{}
}

func NewMacroWithContext(c Context) *Macro {
	return &Macro{
		c: c,
	}
}

func (m *Macro) AddUpdater(u updater) {
	m.u = append(m.u, u)
}

func (m *Macro) AddContext(c Context) {
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

func NewUpdaterFromDatabase(m *MemoryModel, sm *SelectMacro) (updater, error) {
	switch {
	case sm.ModelType.String == "VehicleJourney" && sm.Type == SetAttribute:
		return NewVehicleJourneySetAttributeUpdater(sm)
	case sm.ModelType.String == "StopVisit" && sm.Type == DefineAimedScheduledTimes:
		return NewStopVisitDefineAimedScheduledTimesUpdater(sm)
	case sm.ModelType.String == "Situation" && sm.Type == DefineSituationAffects:
		return NewDefineSituationAffectsUpdater(sm)
	case sm.Type == CreateCode:
		return NewCreateCodeUpdater(m, sm)
	}
	return nil, nil
}

func IsUpdater(u string) bool {
	return slices.Contains(updaters, u)
}
