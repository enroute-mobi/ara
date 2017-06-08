package model

import "time"

type StopVisitDepartureStatus string

const (
	STOP_VISIT_DEPARTURE_ONTIME    StopVisitDepartureStatus = "onTime"
	STOP_VISIT_DEPARTURE_EARLY     StopVisitDepartureStatus = "early"
	STOP_VISIT_DEPARTURE_DELAYED   StopVisitDepartureStatus = "delayed"
	STOP_VISIT_DEPARTURE_CANCELLED StopVisitDepartureStatus = "cancelled"
	STOP_VISIT_DEPARTURE_NOREPORT  StopVisitDepartureStatus = "noreport"
	STOP_VISIT_DEPARTURE_DEPARTED  StopVisitDepartureStatus = "departed"
	STOP_VISIT_DEPARTURE_UNDEFINED StopVisitDepartureStatus = ""
)

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

type StopVisitUpdateEvent struct {
	Id                string
	Created_at        time.Time
	RecordedAt        time.Time
	StopVisitObjectid ObjectID
	StopAreaObjectId  ObjectID
	Schedules         StopVisitSchedules
	DepartureStatus   StopVisitDepartureStatus
	ArrivalStatuts    StopVisitArrivalStatus
	VehicleAtStop     bool

	// Vehicle Journey
	DatedVehicleJourneyRef string
	DestinationRef         string
	OriginRef              string
	DestinationName        string
	OriginName             string
	Attributes             StopVisitUpdateAttributes
}
