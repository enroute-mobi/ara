package model

type StopVisitArrivalStatus string

const (
	STOP_VISIT_ARRIVAL_ARRIVED      StopVisitArrivalStatus = "arrived"
	STOP_VISIT_ARRIVAL_ONTIME       StopVisitArrivalStatus = "onTime"
	STOP_VISIT_ARRIVAL_EARLY        StopVisitArrivalStatus = "early"
	STOP_VISIT_ARRIVAL_DELAYED      StopVisitArrivalStatus = "delayed"
	STOP_VISIT_ARRIVAL_CANCELLED    StopVisitArrivalStatus = "cancelled"
	STOP_VISIT_ARRIVAL_NOREPORT     StopVisitArrivalStatus = "noreport"
	STOP_VISIT_ARRIVAL_MISSED       StopVisitArrivalStatus = "missed"
	STOP_VISIT_ARRIVAL_NOT_EXPECTED StopVisitArrivalStatus = "notExpected"
	STOP_VISIT_ARRIVAL_UNDEFINED    StopVisitArrivalStatus = ""
)

func (svas StopVisitArrivalStatus) Arrived() bool {
	return svas == STOP_VISIT_ARRIVAL_CANCELLED ||
		svas == STOP_VISIT_ARRIVAL_ARRIVED
}

func (svds StopVisitDepartureStatus) Departed() bool {
	return svds == STOP_VISIT_DEPARTURE_CANCELLED ||
		svds == STOP_VISIT_DEPARTURE_DEPARTED
}
