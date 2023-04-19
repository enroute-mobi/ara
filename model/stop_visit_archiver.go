package model

import (
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"cloud.google.com/go/bigquery"
)

type StopVisitArchiver struct {
	Model     Model
	StopVisit *StopVisit
}

func (sva *StopVisitArchiver) StopArea() *StopArea {
	stopArea, _ := sva.Model.StopAreas().Find(sva.StopVisit.StopAreaId)
	return stopArea
}

func (sva *StopVisitArchiver) VehicleJourney() *VehicleJourney {
	vehicleJourney, ok := sva.Model.VehicleJourneys().Find(sva.StopVisit.VehicleJourneyId)
	if !ok {
		return nil
	}
	return vehicleJourney
}

func (sva *StopVisitArchiver) Archive() {
	sv := sva.StopVisit
	sa := sva.StopArea()
	vj := sva.VehicleJourney()
	longTermStopVisitEvent := &audit.BigQueryLongTermStopVisitEvent{
		StopVisitUUID:      string(sv.Id()),
		AimedArrivalTime:   sva.setArrivalTimeEventFromKind(sv, STOP_VISIT_SCHEDULE_AIMED),
		AimedDepartureTime: sva.setDepartureTimeEventFromKind(sv, STOP_VISIT_SCHEDULE_AIMED),

		ExpectedArrivalTime:   sva.setArrivalTimeEventFromKind(sv, STOP_VISIT_SCHEDULE_EXPECTED),
		ExpectedDepartureTime: sva.setDepartureTimeEventFromKind(sv, STOP_VISIT_SCHEDULE_EXPECTED),

		ActualArrivalTime:   sva.setArrivalTimeEventFromKind(sv, STOP_VISIT_SCHEDULE_ACTUAL),
		ActualDepartureTime: sva.setDepartureTimeEventFromKind(sv, STOP_VISIT_SCHEDULE_ACTUAL),

		StopAreaName:        sa.Name,
		StopAreaCoordinates: fmt.Sprintf("POINT(%f %f)", sa.Longitude, sa.Latitude),
		ArrivalStatus:       string(sv.ArrivalStatus),
		DepartureStatus:     string(sv.DepartureStatus),
	}

	if vj != nil {
		longTermStopVisitEvent.VehicleJourneyDirectionType = vj.DirectionType
		longTermStopVisitEvent.VehicleJourneyOriginName = vj.OriginName
		longTermStopVisitEvent.VehicleJourneyDestinationName = vj.DestinationName

		for _, obj := range vj.objectids {
			code := &audit.Code{
				Kind:  obj.kind,
				Value: obj.value,
			}
			longTermStopVisitEvent.VehicleJourneyCodes = append(longTermStopVisitEvent.VehicleJourneyCodes, *code)
		}

		transportMode, ok := vj.Attribute("VehicleMode")
		if ok {
			longTermStopVisitEvent.TransportMode = transportMode
		}

		if vj.Line() != nil {
			longTermStopVisitEvent.LineName = vj.Line().Name
			longTermStopVisitEvent.LineNumber = vj.Line().Number
			for _, obj := range vj.Line().objectids {
				code := &audit.Code{
					Kind:  obj.kind,
					Value: obj.value,
				}
				longTermStopVisitEvent.LineCodes = append(longTermStopVisitEvent.LineCodes, *code)
			}
		}

		if vj.Occupancy != "" {
			longTermStopVisitEvent.VehicleOccupancy = vj.Occupancy
		}
	}

	for _, obj := range sa.objectids {
		code := &audit.Code{
			Kind:  obj.kind,
			Value: obj.value,
		}
		longTermStopVisitEvent.StopAreaCodes = append(longTermStopVisitEvent.StopAreaCodes, *code)
	}

	audit.CurrentBigQuery(sva.Model.Referential()).WriteEvent(longTermStopVisitEvent)
}

func (sva *StopVisitArchiver) setArrivalTimeEventFromKind(sv *StopVisit, kind StopVisitScheduleType) bigquery.NullTimestamp {
	t := bigquery.NullTimestamp{}
	arrivalTime := sv.Schedules.ArrivalTimeFromKind([]StopVisitScheduleType{kind})
	if arrivalTime == (time.Time{}) {
		t.Valid = false
	} else {
		t.Timestamp = arrivalTime
		t.Valid = true
	}

	return t
}

func (sva *StopVisitArchiver) setDepartureTimeEventFromKind(sv *StopVisit, kind StopVisitScheduleType) bigquery.NullTimestamp {
	t := bigquery.NullTimestamp{}
	departureTime := sv.Schedules.DepartureTimeFromKind([]StopVisitScheduleType{kind})
	if departureTime == (time.Time{}) {
		t.Valid = false
	} else {
		t.Timestamp = departureTime
		t.Valid = true
	}

	return t
}
