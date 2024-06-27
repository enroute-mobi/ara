package model

import (
	"encoding/json"
	"errors"
	"fmt"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

var vehicleJourneyIfAttributeContextFactories = map[string]contexFactory{
	"DirectionName": newVehicleJourneyDirectionNameContext,
}

type vehicleJourneyIfAttributeContextAttributes struct {
	AttributeName string `json:"attribute_name"`
	Value         string `json:"value"`
}

func NewVehicleJourneyIfAttributeContext(sm *SelectMacro) (context, error) {
	if !sm.Attributes.Valid {
		return nil, errors.New("empty Attributes")
	}

	var attrs vehicleJourneyIfAttributeContextAttributes
	err := json.Unmarshal([]byte(sm.Attributes.String), &attrs)
	if err != nil {
		return nil, fmt.Errorf("can't parse Attributes: %v", err)
	}

	f, ok := vehicleJourneyIfAttributeContextFactories[attrs.AttributeName]
	if !ok {
		return nil, errors.New("unknown Attribute")
	}

	return f(attrs.Value)
}

func newVehicleJourneyDirectionNameContext(d contextAttributes) (context, error) {
	return func(mi ModelInstance) bool {
		vj := mi.(*VehicleJourney)
		return vj.Attributes[siri_attributes.DirectionName] == d.(string)
	}, nil
}
