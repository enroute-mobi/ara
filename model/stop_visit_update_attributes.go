package model

type StopVisitUpdateAttributes interface {
	StopVisitAttributes() *StopVisitAttributes
	VehiculeJourneyAttributes() *VehicleJourneyAttributes
	LineAttributes() *LineAttributes
	StopAreaAttributes() *StopAreaAttributes
}
