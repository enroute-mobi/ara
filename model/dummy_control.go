package model

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/audit"
)

func NewDummyControler(sc *SelectControl) (controler, error) {
	return func(mi ModelInstance) error {
		var messageAttribute string
		switch sc.ModelType.String {
		case "StopArea":
			messageAttribute = mi.(*StopArea).Name
		case "Line":
			messageAttribute = mi.(*Line).Name
		case "VehicleJourney":
			messageAttribute = mi.(*VehicleJourney).Name
		case "Situation":
			messageAttribute = mi.(*Situation).Summary.DefaultValue
		default:
			messageAttribute = fmt.Sprintf("Don't know how to handle model type %s", sc.ModelType.String)
		}

		m := &audit.BigQueryControlEvent{
			Criticity:                        sc.Criticity.String,
			ControlType:                      "Dummy",
			InternalCode:                     sc.InternalCode.String,
			TargetModelClass:                 sc.ModelType.String,
			TargetModelUUID:                  string(mi.modelId()),
			TranslationInfoMessageKey:        fmt.Sprintf("dummy_%s", sc.ModelType.String),
			TranslationInfoMessageAttributes: messageAttribute,
		}

		audit.CurrentBigQuery(sc.ReferentialSlug).WriteEvent(m)

		return nil
	}, nil
}
