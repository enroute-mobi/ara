package model_types

type Model uint8

const (
	StopArea Model = iota
	Line
	VehicleJourney
	StopVisit
	Vehicle
	Situation

	Total = 6
)

var Type = map[string]Model{
	"StopArea":       StopArea,
	"Line":           Line,
	"VehicleJourney": VehicleJourney,
	"StopVisit":      StopVisit,
	"Vehicle":        Vehicle,
	"Situation":      Situation,
}
