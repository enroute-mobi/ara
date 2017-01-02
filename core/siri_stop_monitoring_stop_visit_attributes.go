package core

import (
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SIRIStopVisitUpdateAttributes struct {
	objectid_kind string
	response      *siri.XMLMonitoredStopVisit
}

func NewSIRIStopVisitUpdateAttributes(response *siri.XMLMonitoredStopVisit, objectid_kind string) *SIRIStopVisitUpdateAttributes {
	return &SIRIStopVisitUpdateAttributes{
		objectid_kind: objectid_kind,
		response:      response,
	}
}

// WIP
func (attributes *SIRIStopVisitUpdateAttributes) StopVisitAttributes() *model.StopVisitAttributes {
	objectid := model.NewObjectID(attributes.objectid_kind, attributes.response.ItemIdentifier())
	stopAreaObjectid := model.NewObjectID(attributes.objectid_kind, attributes.response.StopPointRef())
	vehicleJourneyObjectId := model.NewObjectID(attributes.objectid_kind, attributes.response.DatedVehicleJourneyRef())

	stopVisitAttributes := &model.StopVisitAttributes{
		ObjectId:         objectid,
		StopAreaObjectId: stopAreaObjectid,

		VehicleJourneyObjectId: vehicleJourneyObjectId,
		PassageOrder:           attributes.response.Order(),

		DepartureStatus: model.StopVisitDepartureStatus(attributes.response.DepartureStatus()),
		ArrivalStatus:   model.StopVisitArrivalStatus(attributes.response.ArrivalStatus()),
	}
	stopVisitAttributes.Schedules = model.NewStopVisitSchedules()
	if !attributes.response.AimedDepartureTime().IsZero() || !attributes.response.AimedArrivalTime().IsZero() {
		stopVisitAttributes.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_AIMED, attributes.response.AimedDepartureTime(), attributes.response.AimedArrivalTime())
	}
	if !attributes.response.ExpectedDepartureTime().IsZero() || !attributes.response.ExpectedArrivalTime().IsZero() {
		stopVisitAttributes.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_EXPECTED, attributes.response.ExpectedDepartureTime(), attributes.response.ExpectedArrivalTime())
	}
	if !attributes.response.ActualDepartureTime().IsZero() || !attributes.response.ActualArrivalTime().IsZero() {
		stopVisitAttributes.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_ACTUAL, attributes.response.ActualDepartureTime(), attributes.response.ActualArrivalTime())
	}
	return stopVisitAttributes
}

func (attributes *SIRIStopVisitUpdateAttributes) VehicleJourneyAttributes() *model.VehicleJourneyAttributes {
	objectid := model.NewObjectID(attributes.objectid_kind, attributes.response.DatedVehicleJourneyRef())
	lineObjectId := model.NewObjectID(attributes.objectid_kind, attributes.response.LineRef())

	vehicleJourneyAttributes := &model.VehicleJourneyAttributes{
		ObjectId:     objectid,
		LineObjectId: lineObjectId,
	}

	return vehicleJourneyAttributes
}

func (attributes *SIRIStopVisitUpdateAttributes) LineAttributes() *model.LineAttributes {
	objectid := model.NewObjectID(attributes.objectid_kind, attributes.response.LineRef())

	lineAttributes := &model.LineAttributes{
		ObjectId: objectid,
		Name:     attributes.response.PublishedLineName(),
	}

	return lineAttributes
}

func (attributes *SIRIStopVisitUpdateAttributes) StopAreaAttributes() *model.StopAreaAttributes {
	objectid := model.NewObjectID(attributes.objectid_kind, attributes.response.StopPointRef())

	stopAreaAttributes := &model.StopAreaAttributes{
		ObjectId: objectid,
		Name:     attributes.response.StopPointName(),
	}

	return stopAreaAttributes
}
