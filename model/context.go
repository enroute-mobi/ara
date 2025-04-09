package model

import (
	"errors"
	"slices"

	"bitbucket.org/enroute-mobi/ara/model/hooks"
	"bitbucket.org/enroute-mobi/ara/model/model_types"
)

type ContextAttributes interface{}
type Context func(ModelInstance) bool
type ContexFactory func(ContextAttributes) (Context, error)

const (
	IfAttribute = "IfAttribute"
)

var Contexes = []string{IfAttribute}

func IsContext(c string) bool {
	return slices.Contains(Contexes, c)
}

func HookAndModelType(s DatabaseStructureWithContext) (h hooks.Type, t model_types.Model, errs []error) {
	if !s.GetModelType().Valid {
		errs = append(errs, errors.New("Control with invalid Type"))
		return
	}
	t = model_types.Type[s.GetModelType().String]

	h, ok := hooks.Hook[s.GetHook().String]
	if !ok {
		h = hooks.AfterSave
	}
	return
}

func NewContexFromDatabase(s DatabaseStructureWithContext) (Context, error) {
	if s.GetModelType().String == "VehicleJourney" && s.GetType() == IfAttribute {
		return NewVehicleJourneyIfAttributeContext(s.GetAttributes())
	}
	return nil, nil
}
