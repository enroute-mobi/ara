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
	stopVisitAttributes.Attributes = attributes.FillStopVisitAttributes()
	return stopVisitAttributes
}

func (attributes *SIRIStopVisitUpdateAttributes) FillVehicleJourneyAttributes() map[string]string {
	attrMap := make(map[string]string)

	attrMap["Delay"] = attributes.response.Delay()
	attrMap["DirectionName"] = attributes.response.DirectionName()
	attrMap["DestinationName"] = attributes.response.DestinationName()
	attrMap["DirectionRef"] = attributes.response.DirectionRef()
	attrMap["FirstOrLastJourney"] = attributes.response.FirstOrLastJourney()
	attrMap["HeadwayService"] = attributes.response.HeadwayService()
	attrMap["JourneyNote"] = attributes.response.JourneyNote()
	attrMap["JourneyPatternName"] = attributes.response.JourneyPatternName()
	attrMap["Monitored"] = attributes.response.Monitored()
	attrMap["MonitoringError"] = attributes.response.MonitoringError()
	attrMap["Occupancy"] = attributes.response.Occupancy()
	attrMap["OriginAimedDepartureTime"] = attributes.response.OriginAimedDepartureTime()
	attrMap["DestinationAimedArrivalTime"] = attributes.response.DestinationAimedArrivalTime()
	attrMap["OriginName"] = attributes.response.OriginName()
	attrMap["ProductCategoryRef"] = attributes.response.ProductCategoryRef()
	attrMap["ServiceFeatureRef"] = attributes.response.ServiceFeatureRef()
	attrMap["TrainNumberRef"] = attributes.response.TrainNumberRef()
	attrMap["VehicleFeature"] = attributes.response.VehicleFeature()
	attrMap["VehicleMode"] = attributes.response.VehicleMode()
	attrMap["ViaPlaceName"] = attributes.response.ViaPlaceName()
	attrMap["VehicleJourneyName"] = attributes.response.VehicleJourneyName()

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

	return refMap
}

func (attributes *SIRIStopVisitUpdateAttributes) FillStopVisitAttributes() map[string]string {
	attrMap := make(map[string]string)

	attrMap["Delay"] = attributes.response.Delay()
	attrMap["ActualQuayName"] = attributes.response.ActualQuayName()
	attrMap["AimedHeadwayInterval"] = attributes.response.AimedHeadwayInterval()
	attrMap["ArrivalPlatformName"] = attributes.response.ArrivalPlatformName()
	attrMap["ArrivalProximyTest"] = attributes.response.ArrivalProximyTest()
	attrMap["DepartureBoardingActivity"] = attributes.response.DepartureBoardingActivity()
	attrMap["DeparturePlatformName"] = attributes.response.DeparturePlatformName()
	attrMap["DestinationDisplay"] = attributes.response.DestinationDisplay()
	attrMap["DistanceFromStop"] = attributes.response.DistanceFromStop()
	attrMap["ExpectedHeadwayInterval"] = attributes.response.ExpectedHeadwayInterval()
	attrMap["NumberOfStopsAway"] = attributes.response.NumberOfStopsAway()
	attrMap["PlatformTraversal"] = attributes.response.PlatformTraversal()

	return attrMap
}

func (attributes *SIRIStopVisitUpdateAttributes) VehicleJourneyAttributes() *model.VehicleJourneyAttributes {
	objectid := model.NewObjectID(attributes.objectid_kind, attributes.response.DatedVehicleJourneyRef())
	lineObjectId := model.NewObjectID(attributes.objectid_kind, attributes.response.LineRef())

	vehicleJourneyAttributes := &model.VehicleJourneyAttributes{
		ObjectId:     objectid,
		LineObjectId: lineObjectId,
		Attributes:   make(map[string]string),
		References:   make(map[string]model.Reference),
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
