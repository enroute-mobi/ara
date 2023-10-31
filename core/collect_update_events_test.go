package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetLines(t *testing.T) {
	assert := assert.New(t)
	updatedRefs := NewCollectedRefs()
	updatedRefs.LineRefs["NINOXE:Line:A:LOC"] = struct{}{}
	updatedRefs.LineRefs["NINOXE:Line:B:LOC"] = struct{}{}

	expected := []string{"NINOXE:Line:A:LOC", "NINOXE:Line:B:LOC"}
	assert.ElementsMatch(expected, updatedRefs.GetLines())
}

func Test_GetStopAreas(t *testing.T) {
	assert := assert.New(t)
	updatedRefs := NewCollectedRefs()
	updatedRefs.MonitoringRefs["NINOXE:Stop:A:LOC"] = struct{}{}
	updatedRefs.MonitoringRefs["NINOXE:Stop:B:LOC"] = struct{}{}

	expected := []string{"NINOXE:Stop:A:LOC", "NINOXE:Stop:B:LOC"}
	assert.ElementsMatch(expected, updatedRefs.GetStopAreas())
}

func Test_GetVehicleJourneys(t *testing.T) {
	assert := assert.New(t)
	updatedRefs := NewCollectedRefs()
	updatedRefs.VehicleJourneyRefs["NINOXE:VehicleJourney:A:LOC"] = struct{}{}
	updatedRefs.VehicleJourneyRefs["NINOXE:VehicleJourney:B:LOC"] = struct{}{}

	expected := []string{"NINOXE:VehicleJourney:A:LOC", "NINOXE:VehicleJourney:B:LOC"}
	assert.ElementsMatch(expected, updatedRefs.GetVehicleJourneys())
}

func Test_GetVehicles(t *testing.T) {
	assert := assert.New(t)
	updatedRefs := NewCollectedRefs()
	updatedRefs.VehicleRefs["RLA290"] = struct{}{}
	updatedRefs.VehicleRefs["RLA800"] = struct{}{}

	expected := []string{"RLA290", "RLA800"}
	assert.ElementsMatch(expected, updatedRefs.GetVehicles())
}
