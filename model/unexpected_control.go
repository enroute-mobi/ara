package model

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"bitbucket.org/enroute-mobi/ara/audit"
)

func NewUnexpectedController(sc *SelectControl) (controller, error) {
	if sc.Hook.String != "AfterCreate" {
		return nil, errors.New("'unexpected' controller must be defined AfterCreate")
	}
	if !slices.Contains([]string{"StopArea", "Line", "VehicleJourney"}, sc.ModelType.String) {
		return nil, fmt.Errorf("don't know how to handle model type %s in 'unexpected' controller", sc.ModelType.String)
	}

	return func(mi ModelInstance) error {
		var messageAttribute string
		switch sc.ModelType.String {
		case "StopArea":
			messageAttribute = mi.(*StopArea).Name
		case "Line":
			messageAttribute = mi.(*Line).Name
		case "VehicleJourney":
			messageAttribute = mi.(*VehicleJourney).Name
		}

		m := &audit.BigQueryControlEvent{
			Criticity:                        sc.Criticity.String,
			ControlType:                      "Unexpected",
			InternalCode:                     sc.InternalCode.String,
			TargetModelClass:                 sc.ModelType.String,
			TargetModelUUID:                  string(mi.ModelId()),
			TranslationInfoMessageKey:        fmt.Sprintf("unexpected_%s", strings.ToLower(sc.ModelType.String)),
			TranslationInfoMessageAttributes: messageAttribute,
		}

		audit.CurrentBigQuery(sc.ReferentialSlug).WriteEvent(m)

		return nil
	}, nil
}
