package model

import (
	"encoding/json"
	"errors"
	"fmt"
)

var vehicleJourneySetAttributeUpdaterFactories = map[string]updaterFactory{
	"DirectionType": newVehicleJourneyDirectionTypeUpdater,
}

type vehicleJourneySetAttributeUpdaterAttributes struct {
	AttributeName string `json:"attribute_name"`
	Value         string `json:"value"`
}

func NewVehicleJourneySetAttributeUpdater(sm *SelectMacro) (updater, error) {
	if !sm.Attributes.Valid {
		return nil, errors.New("empty Attributes")
	}

	var attrs vehicleJourneySetAttributeUpdaterAttributes
	err := json.Unmarshal([]byte(sm.Attributes.String), &attrs)
	if err != nil {
		return nil, fmt.Errorf("can't parse Attributes: %v", err)
	}

	f, ok := vehicleJourneySetAttributeUpdaterFactories[attrs.AttributeName]
	if !ok {
		return nil, errors.New("unknown Attribute")
	}

	return f(attrs.Value)
}

func newVehicleJourneyDirectionTypeUpdater(d updaterAttributes) (updater, error) {
	return func(mi ModelInstance) error {
		vj := mi.(*VehicleJourney)
		vj.DirectionType = d.(string)
		return nil
	}, nil
}
