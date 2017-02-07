package model

import "time"

type StopVisitUpdateAttributes interface {
	StopVisitAttributes() *StopVisitAttributes
	VehicleJourneyAttributes() *VehicleJourneyAttributes
	LineAttributes() *LineAttributes
	StopAreaAttributes() *StopAreaAttributes
}

type TestStopVisitUpdateAttributes struct{}

func (t *TestStopVisitUpdateAttributes) StopVisitAttributes() *StopVisitAttributes {
	objectid := NewObjectID("kind", "value")
	return &StopVisitAttributes{
		ObjectId:               objectid,
		StopAreaObjectId:       objectid,
		VehicleJourneyObjectId: objectid,
		RecordedAt:             time.Time{},
		PassageOrder:           1,
		DepartureStatus:        STOP_VISIT_DEPARTURE_CANCELLED,
		ArrivalStatus:          STOP_VISIT_ARRIVAL_CANCELLED,
		// 	Schedules       StopVisitSchedules
	}
}

func (t *TestStopVisitUpdateAttributes) VehicleJourneyAttributes() *VehicleJourneyAttributes {
	objectid := NewObjectID("kind", "value")
	return &VehicleJourneyAttributes{
		ObjectId:     objectid,
		LineObjectId: objectid,
	}
}

func (t *TestStopVisitUpdateAttributes) LineAttributes() *LineAttributes {
	objectid := NewObjectID("kind", "value")
	return &LineAttributes{
		ObjectId: objectid,
		Name:     "line",
	}
}

func (t *TestStopVisitUpdateAttributes) StopAreaAttributes() *StopAreaAttributes {
	objectid := NewObjectID("kind", "value")
	return &StopAreaAttributes{
		ObjectId: objectid,
		Name:     "StopArea",
	}
}
