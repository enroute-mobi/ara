package siri

import "time"

type StopVisitScheduleType int

const (
	AIMED StopVisitScheduleType = iota
	EXPECTED
	ARRIVAL
)

type StopVisitSchedule struct {
	kind          StopVisitScheduleType
	departureTime time.Time
	arrivalTime   time.Time
}
