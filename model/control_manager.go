package model

import (
	"errors"
	"fmt"
	"sync"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model/hooks"
	"bitbucket.org/enroute-mobi/ara/model/model_types"
)

type controls [][][]Control

type Controls interface {
	Loadable

	GetSimpleControls(hooks.Type, model_types.Model) []Control
	GetComplexControls(model_types.Model) map[hooks.Type][]Control
}

type ControlManager struct {
	mutex *sync.RWMutex

	controls controls
}

func NewControlManager() Controls {
	m := &ControlManager{
		mutex: &sync.RWMutex{},
	}
	m.reset()
	return m
}

func (mm *ControlManager) Reset() {
	mm.mutex.Lock()
	mm.reset()
	mm.mutex.Unlock()
}

func (mm *ControlManager) reset() {
	mm.controls = make([][][]Control, hooks.Total)
	for i := range hooks.Total {
		mm.controls[i] = make([][]Control, model_types.Total)
	}
}

/* Unused but I keep the methods here for now

func (mm *ControlManager) SetControl(h hooks.Type, t ModelType, m Control) {
	mm.mutex.Lock()
	mm.setControl(h, t, m)
	mm.mutex.Unlock()
}

func (mm *ControlManager) setControl(h hooks.Type, t ModelType, m Control) {
	mm.controls[h][t] = append(mm.controls[h][t], m)
}
*/

// If we ask for AfterCreate, we'll also get AfterSave Controls
func (mm ControlManager) GetSimpleControls(h hooks.Type, t model_types.Model) (m []Control) {
	for i := h; i < hooks.TotalSimpleHooks; i++ {
		m = append(m, mm.controls[i][t]...)
	}
	return
}

func (mm ControlManager) GetComplexControls(t model_types.Model) (m map[hooks.Type][]Control) {
	m = make(map[hooks.Type][]Control)
	for i := hooks.TotalSimpleHooks; i < hooks.Total; i++ {
		m[hooks.Type(i)] = append(m[hooks.Type(i)], mm.controls[hooks.Type(i)][t]...)
	}
	return
}

type controlBuilder struct {
	manager        *ControlManager
	initialContext []*controlContextBuilder
	contexes       map[string]*controlContextBuilder
}

type controlContextBuilder struct {
	childrenId  string
	control     *SelectControl
	controllers []*SelectControl
}

func (b *controlBuilder) buildControls() []error {
	e := []error{}
	for _, c := range b.initialContext {
		if c.control != nil { // We are handling a context
			e = append(e, b.buildContext(c)...)
		} else { // We should only have one Controller at a time
			for i := range c.controllers {
				e = append(e, b.buildController(c.controllers[i])...)
			}
		}
	}
	return e
}

func (b *controlBuilder) buildContext(c *controlContextBuilder) []error {
	h, mt, errs := HookAndModelType(c.control)
	if len(errs) != 0 {
		return errs
	}

	e := []error{}

	m := NewControl()
	m.Id = c.control.Id
	e = append(e, b.handleContexes(c, m)...)
	b.manager.controls[h][mt] = append(b.manager.controls[h][mt], *m)

	return e
}

func (b *controlBuilder) buildController(sm *SelectControl) []error {
	h, mt, errs := HookAndModelType(sm)
	if len(errs) != 0 {
		return errs
	}

	e := []error{}

	m := NewControl()
	m.Id = sm.Id
	updater, err := NewControllerFromDatabase(sm)
	if err != nil {
		e = append(e, err)
		return e
	}
	m.AddController(updater)
	b.manager.controls[h][mt] = append(b.manager.controls[h][mt], *m)

	return e
}

func (b *controlBuilder) handleContexes(c *controlContextBuilder, m *Control) []error {
	e := []error{}
	if c.control != nil {
		context, err := NewContexFromDatabase(c.control)
		if err != nil {
			e = append(e, err)
			return e
		}
		m.AddContext(context)
	}
	for _, u := range c.controllers {
		updater, err := NewControllerFromDatabase(u)
		if err != nil {
			e = append(e, err)
			continue
		}
		m.AddController(updater)
	}
	if c.childrenId != "" {
		b.handleContexes(b.contexes[c.childrenId], m)
	}
	return e
}

func (manager *ControlManager) Load(referentialSlug string) error {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	manager.reset()

	builder := &controlBuilder{
		manager:  manager,
		contexes: make(map[string]*controlContextBuilder),
	}

	var selectControls []SelectControl

	sqlQuery := fmt.Sprintf("select * from controls where referential_slug = '%s' order by context_id nulls first, position", referentialSlug)
	_, err := Database.Select(&selectControls, sqlQuery)
	if err != nil {
		return err
	}

	for _, sm := range selectControls {
		if !sm.ContextId.Valid {
			context := &controlContextBuilder{
				controllers: make([]*SelectControl, 0),
			}
			if IsContext(sm.Type) {
				context.control = &sm
			} else {
				context.controllers = append(context.controllers, &sm)
			}
			builder.initialContext = append(builder.initialContext, context)
			builder.contexes[sm.Id] = context
			continue
		}

		parent := builder.contexes[sm.ContextId.String]
		if IsContext(sm.Type) {
			parent.childrenId = sm.Id
			context := &controlContextBuilder{}
			context.control = &sm
			builder.contexes[sm.Id] = context
			continue
		}

		parent.controllers = append(parent.controllers, &sm)
	}

	errs := builder.buildControls()
	if len(errs) != 0 {
		e := fmt.Sprintf("errors while loading Controls: %v", errs)
		logger.Log.Debug(e)
		return errors.New(e)
	}
	return nil
}
