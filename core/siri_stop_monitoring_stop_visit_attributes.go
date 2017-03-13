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
		RecordedAt:             attributes.response.RecordedAt(),

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
	stopVisitAttributes.Attributes = attributes.FillStopVisitAttributes()
	stopVisitAttributes.References = attributes.FillStopVisitReferences()
	return stopVisitAttributes
}

func (attributes *SIRIStopVisitUpdateAttributes) FillVehicleJourneyAttributes() map[string]string {
	attrMap := make(map[string]string)
	tmpattrMap := make(map[string]string)

	tmpattrMap["Delay"] = attributes.response.Delay()
	tmpattrMap["Bearing"] = attributes.response.Bearing()
	tmpattrMap["InPanic"] = attributes.response.InPanic()
	tmpattrMap["InCongestion"] = attributes.response.InCongestion()
	tmpattrMap["SituationRef"] = attributes.response.SituationRef()
	tmpattrMap["DirectionName"] = attributes.response.DirectionName()
	tmpattrMap["DestinationName"] = attributes.response.DestinationName()
	tmpattrMap["DirectionRef"] = attributes.response.DirectionRef()
	tmpattrMap["FirstOrLastJourney"] = attributes.response.FirstOrLastJourney()
	tmpattrMap["HeadwayService"] = attributes.response.HeadwayService()
	tmpattrMap["JourneyNote"] = attributes.response.JourneyNote()
	tmpattrMap["JourneyPatternName"] = attributes.response.JourneyPatternName()
	tmpattrMap["Monitored"] = attributes.response.Monitored()
	tmpattrMap["MonitoringError"] = attributes.response.MonitoringError()
	tmpattrMap["Occupancy"] = attributes.response.Occupancy()
	tmpattrMap["OriginAimedDepartureTime"] = attributes.response.OriginAimedDepartureTime()
	tmpattrMap["DestinationAimedArrivalTime"] = attributes.response.DestinationAimedArrivalTime()
	tmpattrMap["OriginName"] = attributes.response.OriginName()
	tmpattrMap["ProductCategoryRef"] = attributes.response.ProductCategoryRef()
	tmpattrMap["ServiceFeatureRef"] = attributes.response.ServiceFeatureRef()
	tmpattrMap["TrainNumberRef"] = attributes.response.TrainNumberRef()
	tmpattrMap["VehicleFeature"] = attributes.response.VehicleFeature()
	tmpattrMap["VehicleMode"] = attributes.response.VehicleMode()
	tmpattrMap["ViaPlaceName"] = attributes.response.ViaPlaceName()
	tmpattrMap["VehicleJourneyName"] = attributes.response.VehicleJourneyName()

	for key, value := range tmpattrMap {
		if value != "" {
			attrMap[key] = value
		}
	}

	return attrMap
}

func (attributes *SIRIStopVisitUpdateAttributes) FillVehicleJourneyReferences() map[string]model.Reference {
	refMap := make(map[string]model.Reference)

	if attributes.response.PlaceRef() != "" {
		placeRefObjId := model.NewObjectID(attributes.objectid_kind, attributes.response.PlaceRef())
		refMap["PlaceRef"] = model.Reference{ObjectId: &placeRefObjId, Id: ""}
	}

	if attributes.response.OriginRef() != "" {
		originRefObjId := model.NewObjectID(attributes.objectid_kind, attributes.response.OriginRef())
		refMap["OriginRef"] = model.Reference{ObjectId: &originRefObjId, Id: ""}
	}

	if attributes.response.DestinationRef() != "" {
		destinationRefObjId := model.NewObjectID(attributes.objectid_kind, attributes.response.DestinationRef())
		refMap["DestinationRef"] = model.Reference{ObjectId: &destinationRefObjId, Id: ""}
	}

	if attributes.response.JourneyPatternRef() != "" {
		journeyPatternRefObjId := model.NewObjectID(attributes.objectid_kind, attributes.response.JourneyPatternRef())
		refMap["JourneyPatternRef"] = model.Reference{ObjectId: &journeyPatternRefObjId, Id: ""}
	}

	if attributes.response.RouteRef() != "" {
		routeRefObjId := model.NewObjectID(attributes.objectid_kind, attributes.response.RouteRef())
		refMap["RouteRef"] = model.Reference{ObjectId: &routeRefObjId, Id: ""}
	}

	return refMap
}

func (attributes *SIRIStopVisitUpdateAttributes) FillStopVisitAttributes() map[string]string {
	attrMap := make(map[string]string)
	tmpattrMap := make(map[string]string)

	tmpattrMap["Delay"] = attributes.response.Delay()
	tmpattrMap["ActualQuayName"] = attributes.response.ActualQuayName()
	tmpattrMap["AimedHeadwayInterval"] = attributes.response.AimedHeadwayInterval()
	tmpattrMap["ArrivalPlatformName"] = attributes.response.ArrivalPlatformName()
	tmpattrMap["ArrivalProximyTest"] = attributes.response.ArrivalProximyTest()
	tmpattrMap["DepartureBoardingActivity"] = attributes.response.DepartureBoardingActivity()
	tmpattrMap["DeparturePlatformName"] = attributes.response.DeparturePlatformName()
	tmpattrMap["DestinationDisplay"] = attributes.response.DestinationDisplay()
	tmpattrMap["DistanceFromStop"] = attributes.response.DistanceFromStop()
	tmpattrMap["ExpectedHeadwayInterval"] = attributes.response.ExpectedHeadwayInterval()
	tmpattrMap["NumberOfStopsAway"] = attributes.response.NumberOfStopsAway()
	tmpattrMap["PlatformTraversal"] = attributes.response.PlatformTraversal()

	for key, value := range tmpattrMap {
		if value != "" {
			attrMap[key] = value
		}
	}

	return attrMap
}

func (attributes *SIRIStopVisitUpdateAttributes) FillStopVisitReferences() map[string]model.Reference {
	refMap := make(map[string]model.Reference)

	if attributes.response.OperatorRef() != "" {
		OperatorRefObjId := model.NewObjectID(attributes.objectid_kind, attributes.response.OperatorRef())
		refMap["OperatorRef"] = model.Reference{ObjectId: &OperatorRefObjId, Id: ""}
	}
	return refMap
}

func (attributes *SIRIStopVisitUpdateAttributes) VehicleJourneyAttributes() *model.VehicleJourneyAttributes {
	objectid := model.NewObjectID(attributes.objectid_kind, attributes.response.DatedVehicleJourneyRef())
	lineObjectId := model.NewObjectID(attributes.objectid_kind, attributes.response.LineRef())

	vehicleJourneyAttributes := &model.VehicleJourneyAttributes{
		ObjectId:     objectid,
		LineObjectId: lineObjectId,
	}

	vehicleJourneyAttributes.Attributes = attributes.FillVehicleJourneyAttributes()
	vehicleJourneyAttributes.References = attributes.FillVehicleJourneyReferences()
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
