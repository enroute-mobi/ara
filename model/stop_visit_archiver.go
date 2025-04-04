package model

import (
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/model/schedules"
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
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
		PassageOrder:       sv.PassageOrder,
		AimedArrivalTime:   sva.setArrivalTimeEventFromKind(sv, schedules.Aimed),
		AimedDepartureTime: sva.setDepartureTimeEventFromKind(sv, schedules.Aimed),

		ExpectedArrivalTime:   sva.setArrivalTimeEventFromKind(sv, schedules.Expected),
		ExpectedDepartureTime: sva.setDepartureTimeEventFromKind(sv, schedules.Expected),

		ActualArrivalTime:   sva.setArrivalTimeEventFromKind(sv, schedules.Actual),
		ActualDepartureTime: sva.setDepartureTimeEventFromKind(sv, schedules.Actual),

		StopAreaName:        sa.Name,
		StopAreaCoordinates: fmt.Sprintf("POINT(%f %f)", sa.Longitude, sa.Latitude),
		ArrivalStatus:       string(sv.ArrivalStatus),
		DepartureStatus:     string(sv.DepartureStatus),
	}

	if vj != nil {
		longTermStopVisitEvent.VehicleJourneyDirectionType = vj.DirectionType
		longTermStopVisitEvent.VehicleJourneyOriginName = vj.OriginName
		longTermStopVisitEvent.VehicleJourneyDestinationName = vj.DestinationName

		for _, obj := range vj.codes {
			code := &audit.Code{
				CodeSpace: obj.codeSpace,
				Value:     obj.value,
			}
			longTermStopVisitEvent.VehicleJourneyCodes = append(longTermStopVisitEvent.VehicleJourneyCodes, *code)
		}

		transportMode, ok := vj.Attribute(siri_attributes.VehicleMode)
		if ok {
			longTermStopVisitEvent.TransportMode = transportMode
		}

		if vj.Line() != nil {
			longTermStopVisitEvent.LineName = vj.Line().Name
			longTermStopVisitEvent.LineNumber = vj.Line().Number
			for _, obj := range vj.Line().codes {
				code := &audit.Code{
					CodeSpace: obj.codeSpace,
					Value:     obj.value,
				}
				longTermStopVisitEvent.LineCodes = append(longTermStopVisitEvent.LineCodes, *code)
			}
		}

		if vj.Occupancy != "" {
			longTermStopVisitEvent.VehicleOccupancy = vj.Occupancy
		}
	}

	for _, obj := range sa.codes {
		code := &audit.Code{
			CodeSpace: obj.codeSpace,
			Value:     obj.value,
		}
		longTermStopVisitEvent.StopAreaCodes = append(longTermStopVisitEvent.StopAreaCodes, *code)
	}
	vehicle, ok := sva.Model.Vehicles().FindByNextStopVisitId(sv.Id())
	if ok {
		longTermStopVisitEvent.VehicleDriverRef = vehicle.DriverRef
		longTermStopVisitEvent.VehicleOccupancy = vehicle.Occupancy
	}
	audit.CurrentBigQuery(sva.Model.Referential()).WriteEvent(longTermStopVisitEvent)
}

func (sva *StopVisitArchiver) setArrivalTimeEventFromKind(sv *StopVisit, kind schedules.StopVisitScheduleType) bigquery.NullTimestamp {
	t := bigquery.NullTimestamp{}
	arrivalTime := sv.Schedules.ArrivalTimeFromKind([]schedules.StopVisitScheduleType{kind})
	if arrivalTime == (time.Time{}) {
		t.Valid = false
	} else {
		t.Timestamp = arrivalTime
		t.Valid = true
	}

	return t
}

func (sva *StopVisitArchiver) setDepartureTimeEventFromKind(sv *StopVisit, kind schedules.StopVisitScheduleType) bigquery.NullTimestamp {
	t := bigquery.NullTimestamp{}
	departureTime := sv.Schedules.DepartureTimeFromKind([]schedules.StopVisitScheduleType{kind})
	if departureTime == (time.Time{}) {
		t.Valid = false
	} else {
		t.Timestamp = departureTime
		t.Valid = true
	}

	return t
}
