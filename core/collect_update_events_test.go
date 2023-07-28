package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetLines(t *testing.T) {
	assert := assert.New(t)
	updateEvents := NewCollectUpdateEvents()
	updateEvents.LineRefs["NINOXE:Line:A:LOC"] = struct{}{}
	updateEvents.LineRefs["NINOXE:Line:B:LOC"] = struct{}{}

	expected := []string{"NINOXE:Line:A:LOC", "NINOXE:Line:B:LOC"}
	assert.ElementsMatch(expected, updateEvents.GetLines())
}

func Test_GetStopAreas(t *testing.T) {
	assert := assert.New(t)
	updateEvents := NewCollectUpdateEvents()
	updateEvents.MonitoringRefs["NINOXE:Stop:A:LOC"] = struct{}{}
	updateEvents.MonitoringRefs["NINOXE:Stop:B:LOC"] = struct{}{}

	expected := []string{"NINOXE:Stop:A:LOC", "NINOXE:Stop:B:LOC"}
	assert.ElementsMatch(expected, updateEvents.GetStopAreas())
}

func Test_GetVehicleJourneys(t *testing.T) {
	assert := assert.New(t)
	updateEvents := NewCollectUpdateEvents()
	updateEvents.VehicleJourneyRefs["NINOXE:VehicleJourney:A:LOC"] = struct{}{}
	updateEvents.VehicleJourneyRefs["NINOXE:VehicleJourney:B:LOC"] = struct{}{}

	expected := []string{"NINOXE:VehicleJourney:A:LOC", "NINOXE:VehicleJourney:B:LOC"}
	assert.ElementsMatch(expected, updateEvents.GetVehicleJourneys())
}
