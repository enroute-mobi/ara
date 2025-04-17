package model

import (
	"slices"
)

type ControlId string

type controller func(ModelInstance) error

// type controllerAttributes interface{}
// type controllerFactory func(controllerAttributes) (controller, error)

const (
	Dummy      = "Dummy"
	Unexpected = "Unexpected"
)

var allControllers = []string{
	Dummy,
	Unexpected,
}

type Control struct {
	ctx         Context
	controllers []controller

	Criticity    string
	InternalCode string
}

func NewControl() *Control {
	return &Control{}
}

func NewControlWithContext(ctx Context) *Control {
	return &Control{
		ctx: ctx,
	}
}

func (c *Control) AddController(controllers controller) {
	c.controllers = append(c.controllers, controllers)
}

func (c *Control) AddContext(ctx Context) {
	if c.ctx == nil {
		c.ctx = ctx
		return
	}
	c.ctx = func(mi ModelInstance) bool {
		if !c.ctx(mi) {
			return false
		}
		return ctx(mi)
	}
}

func (c *Control) Control(mi ModelInstance) (ok bool, err error) {
	if c.ctx != nil && !c.ctx(mi) {
		return false, nil
	}
	for i := range c.controllers {
		err := c.controllers[i](mi)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func NewControllerFromDatabase(sc *SelectControl) (controller, error) {
	switch {
	case sc.Type == Dummy:
		return NewDummyController(sc)
	case sc.Type == Unexpected:
		return NewUnexpectedController(sc)
	}

	return nil, nil
}

func IsController(c string) bool {
	return slices.Contains(allControllers, c)
}
