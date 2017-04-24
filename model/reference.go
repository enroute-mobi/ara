package model

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
)

type Reference struct {
	ObjectId *ObjectID `json:",omitempty"`
	Id       string    `json:",omitempty"`
	Type     string    `json:",omitempty"`
}

func (reference *Reference) GetSha1() string {
	hasher := sha1.New() // oui, on sait
	hasher.Write([]byte(reference.ObjectId.Value()))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}

func (reference *Reference) Getformat(ref, value string) string {
	allRef := make(map[string]string)

	allRef["PlaceRef"] = "StopPoint:Q:"
	allRef["OriginRef"] = "StopPoint:Q:"
	allRef["DestinationRef"] = "StopPoint:Q:"
	allRef["JourneyPatternRef"] = "JourneyPattern::"
	allRef["RouteRef"] = "Route::"
	allRef["DatedVehicleJourneyRef"] = "VehiculeJourney::"
	allRef["OperatorRef"] = "Operator::"

	formated := fmt.Sprintf("RATPDev:%s%s:LOC", allRef[ref], value)

	return formated
}

func (reference *Reference) UnmarshalJSON(data []byte) error {
	type Alias Reference
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(reference),
	}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	return nil
}
