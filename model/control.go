package model

import "slices"

type ControlId string

type controler func(ModelInstance) error

// type controlerAttributes interface{}
// type controlerFactory func(controlerAttributes) (controler, error)

const (
	Dummy = "Dummy"
)

var controlers = []string{Dummy}

type Control struct {
	ctx        Context
	controlers []controler

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

func (c *Control) AddControler(controlers controler) {
	c.controlers = append(c.controlers, controlers)
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
	for i := range c.controlers {
		err := c.controlers[i](mi)
		if err != nil {
			return false, err
		}
	}
	return true, nil
}

func NewControlerFromDatabase(sc *SelectControl) (controler, error) {
	if sc.Type == Dummy {
		return NewDummyControler(sc)
	}

	return nil, nil
}

func IsControler(c string) bool {
	return slices.Contains(controlers, c)
}
