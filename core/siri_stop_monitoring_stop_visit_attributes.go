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

func (attributes *SIRIStopVisitUpdateAttributes) FillVehicleJourneyAttributes() model.Attributes {
	vehiculeJourneyAttributes := model.NewAttributes()

	vehiculeJourneyAttributes.Set("Delay", attributes.response.Delay())
	vehiculeJourneyAttributes.Set("Bearing", attributes.response.Bearing())
	vehiculeJourneyAttributes.Set("InPanic", attributes.response.InPanic())
	vehiculeJourneyAttributes.Set("InCongestion", attributes.response.InCongestion())
	vehiculeJourneyAttributes.Set("SituationRef", attributes.response.SituationRef())
	vehiculeJourneyAttributes.Set("DirectionName", attributes.response.DirectionName())
	vehiculeJourneyAttributes.Set("DestinationName", attributes.response.DestinationName())
	vehiculeJourneyAttributes.Set("DirectionRef", attributes.response.DirectionRef())
	vehiculeJourneyAttributes.Set("FirstOrLastJourney", attributes.response.FirstOrLastJourney())
	vehiculeJourneyAttributes.Set("HeadwayService", attributes.response.HeadwayService())
	vehiculeJourneyAttributes.Set("JourneyNote", attributes.response.JourneyNote())
	vehiculeJourneyAttributes.Set("JourneyPatternName", attributes.response.JourneyPatternName())
	vehiculeJourneyAttributes.Set("Monitored", attributes.response.Monitored())
	vehiculeJourneyAttributes.Set("MonitoringError", attributes.response.MonitoringError())
	vehiculeJourneyAttributes.Set("Occupancy", attributes.response.Occupancy())
	vehiculeJourneyAttributes.Set("OriginAimedDepartureTime", attributes.response.OriginAimedDepartureTime())
	vehiculeJourneyAttributes.Set("DestinationAimedArrivalTime", attributes.response.DestinationAimedArrivalTime())
	vehiculeJourneyAttributes.Set("OriginName", attributes.response.OriginName())
	vehiculeJourneyAttributes.Set("ProductCategoryRef", attributes.response.ProductCategoryRef())
	vehiculeJourneyAttributes.Set("ServiceFeatureRef", attributes.response.ServiceFeatureRef())
	vehiculeJourneyAttributes.Set("TrainNumberRef", attributes.response.TrainNumberRef())
	vehiculeJourneyAttributes.Set("VehicleFeature", attributes.response.VehicleFeature())
	vehiculeJourneyAttributes.Set("VehicleMode", attributes.response.VehicleMode())
	vehiculeJourneyAttributes.Set("ViaPlaceName", attributes.response.ViaPlaceName())
	vehiculeJourneyAttributes.Set("VehicleJourneyName", attributes.response.VehicleJourneyName())

	return vehiculeJourneyAttributes
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

func (attributes *SIRIStopVisitUpdateAttributes) FillStopVisitAttributes() model.Attributes {
	stopVisitAttributes := model.NewAttributes()

	stopVisitAttributes.Set("Delay", attributes.response.Delay())
	stopVisitAttributes.Set("ActualQuayName", attributes.response.ActualQuayName())
	stopVisitAttributes.Set("AimedHeadwayInterval", attributes.response.AimedHeadwayInterval())
	stopVisitAttributes.Set("ArrivalPlatformName", attributes.response.ArrivalPlatformName())
	stopVisitAttributes.Set("ArrivalProximyTest", attributes.response.ArrivalProximyTest())
	stopVisitAttributes.Set("DepartureBoardingActivity", attributes.response.DepartureBoardingActivity())
	stopVisitAttributes.Set("DeparturePlatformName", attributes.response.DeparturePlatformName())
	stopVisitAttributes.Set("DestinationDisplay", attributes.response.DestinationDisplay())
	stopVisitAttributes.Set("DistanceFromStop", attributes.response.DistanceFromStop())
	stopVisitAttributes.Set("ExpectedHeadwayInterval", attributes.response.ExpectedHeadwayInterval())
	stopVisitAttributes.Set("NumberOfStopsAway", attributes.response.NumberOfStopsAway())
	stopVisitAttributes.Set("PlatformTraversal", attributes.response.PlatformTraversal())

	return stopVisitAttributes
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
