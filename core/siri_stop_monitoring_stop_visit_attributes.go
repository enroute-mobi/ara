package core

import (
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SIRIStopMonitoringStopVisitAttributes struct {
	objectid_kind string
	response      *siri.XMLMonitoredStopVisit
}

func NewSIRIStopMonitoringStopVisitAttributes(response *siri.XMLMonitoredStopVisit, objectid_kind string) *SIRIStopMonitoringStopVisitAttributes {
	return &SIRIStopMonitoringStopVisitAttributes{
		objectid_kind: objectid_kind,
		response:      response,
	}
}

func (attributes *SIRIStopMonitoringStopVisitAttributes) StopVisitAttributes() *model.StopVisitAttributes {
	objectid := model.NewObjectID(attributes.objectid_kind, attributes.response.ItemIdentifier())
	stopAreaObjectid := model.NewObjectID(attributes.objectid_kind, attributes.response.MonitoringRef())
	vehicleJourneyObjectId := model.NewObjectID(attributes.objectid_kind, attributes.response.DatedVehicleJourneyRef())

	stopVisitAttributes := &model.StopVisitAttributes{
		ObjectId:         &objectid,
		StopAreaObjectId: &stopAreaObjectid,

		VehicleJourneyObjectId: &vehicleJourneyObjectId,
		PassageOrder:           attributes.response.Order(),

		DepartureStatus: model.StopVisitDepartureStatus(attributes.response.DepartureStatus()),
		ArrivalStatus:   model.StopVisitArrivalStatus(attributes.response.ArrivalStatus()),
	}
	stopVisitAttributes.Schedules = model.NewStopVisitSchedules()
	stopVisitAttributes.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_AIMED, attributes.response.AimedDepartureTime(), attributes.response.AimedArrivalTime())
	stopVisitAttributes.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_EXPECTED, attributes.response.ExpectedDepartureTime(), attributes.response.ExpectedArrivalTime())
	stopVisitAttributes.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_ACTUAL, attributes.response.ActualDepartureTime(), attributes.response.ActualArrivalTime())
	return stopVisitAttributes
}

// func (attributes *SIRIStopMonitoringStopVisitAttributes) VehiculeJourneyAttributes() *model.VehiculeJourneyAttributes {
// }

// func (attributes *SIRIStopMonitoringStopVisitAttributes) LineAttributes() *model.LineAttributes {
// }

// func (attributes *SIRIStopMonitoringStopVisitAttributes) StopAreaAttributes() *model.StopAreaAttributes {
// }
