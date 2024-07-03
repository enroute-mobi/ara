package model

import (
	"errors"
	"fmt"
	"sync"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type hook uint8
type ModelType uint8
type macros [][][]Macro

type Macros interface {
	GetMacros(hook, ModelType) []Macro
}

const (
	// Warning: Hooks needs to be sorted
	AfterCreate hook = iota
	AfterSave

	totalHookNumber = 2
)

const (
	StopAreaType ModelType = iota
	LineType
	VehicleJourneyType
	StopVisitType
	VehicleType
	MacroSituationType

	totalModelTypes = 6
)

var modelType = map[string]ModelType{
	"StopArea":       StopAreaType,
	"Line":           LineType,
	"VehicleJourney": VehicleJourneyType,
	"StopVisit":      StopVisitType,
	"Vehicle":        VehicleType,
	"Situation":      MacroSituationType,
}

var hooks = map[string]hook{
	"AfterCreate": AfterCreate,
	"AfterSave":   AfterSave,
}

type MacroManager struct {
	mutex *sync.RWMutex

	macros macros
}

func NewMacroManager() *MacroManager {
	m := &MacroManager{
		mutex: &sync.RWMutex{},
	}
	m.reset()
	return m
}

func (mm *MacroManager) Reset() {
	mm.mutex.Lock()
	mm.reset()
	mm.mutex.Unlock()
}

func (mm *MacroManager) reset() {
	mm.macros = make([][][]Macro, totalHookNumber)
	for i := range totalHookNumber {
		mm.macros[i] = make([][]Macro, totalModelTypes)
		for j := range totalModelTypes {
			mm.macros[i][j] = make([]Macro, 5)
		}
	}
}

/* Unused but I keep the methods here for now

func (mm *MacroManager) SetMacro(h hook, t ModelType, m Macro) {
	mm.mutex.Lock()
	mm.setMacro(h, t, m)
	mm.mutex.Unlock()
}

func (mm *MacroManager) setMacro(h hook, t ModelType, m Macro) {
	mm.macros[h][t] = append(mm.macros[h][t], m)
}
*/

// If we ask for AfterCreate, we'll also get AfterSave Macros
func (mm MacroManager) GetMacros(h hook, t ModelType) (m []Macro) {
	for i := 0; i < totalHookNumber; i++ {
		m = append(m, mm.macros[i][t]...)
	}
	return
}

type macroBuilder struct {
	manager        *MacroManager
	initialContext []*contextBuilder
	contexes       map[string]*contextBuilder
}

type contextBuilder struct {
	childrenId string
	macro      *SelectMacro
	updaters   []*SelectMacro
}

func (b *macroBuilder) buildMacros() []error {
	e := []error{}
	for _, c := range b.initialContext {
		if c.macro != nil { // We are handling a context
			e = append(e, b.buildContext(c)...)
		} else { // We should only have one Updater at a time
			for i := range c.updaters {
				e = append(e, b.buildUpdater(c.updaters[i])...)
			}
		}
	}
	return e
}

func (b *macroBuilder) buildContext(c *contextBuilder) []error {
	h, mt, errs := hookAndModelType(c.macro)
	if len(errs) != 0 {
		return errs
	}

	e := []error{}

	m := NewMacro()
	e = append(e, b.handleContexes(c, m)...)
	b.manager.macros[h][mt] = append(b.manager.macros[h][mt], *m)

	return e
}

func (b *macroBuilder) buildUpdater(sm *SelectMacro) []error {
	h, mt, errs := hookAndModelType(sm)
	if len(errs) != 0 {
		return errs
	}

	e := []error{}

	m := NewMacro()
	updater, err := NewUpdaterFromDatabase(sm)
	if err != nil {
		e = append(e, err)
		return e
	}
	m.AddUpdater(updater)
	b.manager.macros[h][mt] = append(b.manager.macros[h][mt], *m)

	return e
}

func hookAndModelType(sm *SelectMacro) (h hook, mt ModelType, errs []error) {
	if !sm.ModelType.Valid {
		errs = append(errs, errors.New("Macro with invalid Type"))
		return
	}
	mt = modelType[sm.ModelType.String]

	h, ok := hooks[sm.Hook.String]
	if !ok {
		h = AfterSave
	}
	return
}

func (b *macroBuilder) handleContexes(c *contextBuilder, m *Macro) []error {
	e := []error{}
	if c.macro != nil {
		context, err := NewContexFromDatabase(c.macro)
		if err != nil {
			e = append(e, err)
			return e
		}
		m.AddContext(context)
	}
	for _, u := range c.updaters {
		updater, err := NewUpdaterFromDatabase(u)
		if err != nil {
			e = append(e, err)
			continue
		}
		m.AddUpdater(updater)
	}
	if c.childrenId != "" {
		b.handleContexes(b.contexes[c.childrenId], m)
	}
	return e
}

func (manager *MacroManager) Load(referentialSlug string) error {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	manager.reset()

	builder := &macroBuilder{
		manager:  manager,
		contexes: make(map[string]*contextBuilder),
	}

	var selectMacros []SelectMacro

	sqlQuery := fmt.Sprintf("select * from macros where referential_slug = '%s' order by context_id nulls first, position", referentialSlug)
	_, err := Database.Select(&selectMacros, sqlQuery)
	if err != nil {
		return err
	}

	for _, sm := range selectMacros {
		if !sm.ContextId.Valid {
			context := &contextBuilder{
				updaters: make([]*SelectMacro, 0),
			}
			if IsContext(sm.Type) {
				context.macro = &sm
			} else {
				context.updaters = append(context.updaters, &sm)
			}
			builder.initialContext = append(builder.initialContext, context)
			builder.contexes[sm.Id] = context
			continue
		}

		parent := builder.contexes[sm.ContextId.String]
		if IsContext(sm.Type) {
			parent.childrenId = sm.Id
			context := &contextBuilder{}
			context.macro = &sm
			builder.contexes[sm.Id] = context
			continue
		}

		parent.updaters = append(parent.updaters, &sm)
	}

	errs := builder.buildMacros()
	if len(errs) != 0 {
		logger.Log.Debugf("errors while loading Macros: %v", errs)
		return fmt.Errorf("errors while loading Macros: %v", errs)
	}
	return nil
}
