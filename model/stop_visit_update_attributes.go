package model

type StopVisitUpdateAttributes interface {
	StopVisitAttributes() *StopVisitAttributes
	VehicleJourneyAttributes() *VehicleJourneyAttributes
	LineAttributes() *LineAttributes
	StopAreaAttributes() *StopAreaAttributes
}
