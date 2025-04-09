package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

var vehicleJourneyIfAttributeContextFactories = map[string]ContexFactory{
	"DirectionName": newVehicleJourneyDirectionNameContext,
}

type vehicleJourneyIfAttributeContextAttributes struct {
	AttributeName string `json:"attribute_name"`
	Value         string `json:"value"`
}

func NewVehicleJourneyIfAttributeContext(attributes sql.NullString) (Context, error) {
	if !attributes.Valid {
		return nil, errors.New("empty Attributes")
	}

	var attrs vehicleJourneyIfAttributeContextAttributes
	err := json.Unmarshal([]byte(attributes.String), &attrs)
	if err != nil {
		return nil, fmt.Errorf("can't parse Attributes: %v", err)
	}

	f, ok := vehicleJourneyIfAttributeContextFactories[attrs.AttributeName]
	if !ok {
		return nil, errors.New("unknown Attribute")
	}

	return f(attrs.Value)
}

func newVehicleJourneyDirectionNameContext(d ContextAttributes) (Context, error) {
	return func(mi ModelInstance) bool {
		vj := mi.(*VehicleJourney)
		return vj.Attributes[siri_attributes.DirectionName] == d.(string)
	}, nil
}
