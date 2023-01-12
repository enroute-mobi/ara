package model

type StopVisitDepartureStatus string

const (
	STOP_VISIT_DEPARTURE_ONTIME    StopVisitDepartureStatus = "onTime"
	STOP_VISIT_DEPARTURE_EARLY     StopVisitDepartureStatus = "early"
	STOP_VISIT_DEPARTURE_DELAYED   StopVisitDepartureStatus = "delayed"
	STOP_VISIT_DEPARTURE_CANCELLED StopVisitDepartureStatus = "cancelled"
	STOP_VISIT_DEPARTURE_NOREPORT  StopVisitDepartureStatus = "noreport"
	STOP_VISIT_DEPARTURE_DEPARTED  StopVisitDepartureStatus = "departed"
)

func SetStopVisitDepartureStatus(departureStatus string) StopVisitDepartureStatus {
	switch departureStatus {
	case "":
		return STOP_VISIT_DEPARTURE_ONTIME
	default:
		return StopVisitDepartureStatus(departureStatus)
	}
}
