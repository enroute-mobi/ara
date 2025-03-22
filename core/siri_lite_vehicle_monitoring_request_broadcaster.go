package core

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/model/schedules"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type VehicleMonitoringLiteRequestBroadcaster interface {
	RequestVehicles(url.Values, *audit.BigQueryMessage) *siri.SiriLiteResponse
}

type SIRILiteVehicleMonitoringRequestBroadcaster struct {
	connector

	vjRemoteCodeSpaces      []string
	vehicleRemoteCodeSpaces []string
}

type SIRILiteVehicleMonitoringRequestBroadcasterFactory struct{}

func NewSIRILiteVehicleMonitoringRequestBroadcaster(partner *Partner) *SIRILiteVehicleMonitoringRequestBroadcaster {
	connector := &SIRILiteVehicleMonitoringRequestBroadcaster{
		vjRemoteCodeSpaces:      partner.VehicleJourneyRemoteCodeSpaceWithFallback(SIRI_LITE_VEHICLE_MONITORING_REQUEST_BROADCASTER),
		vehicleRemoteCodeSpaces: partner.VehicleRemoteCodeSpaceWithFallback(SIRI_LITE_VEHICLE_MONITORING_REQUEST_BROADCASTER),
	}
	connector.remoteCodeSpace = partner.RemoteCodeSpace(SIRI_LITE_VEHICLE_MONITORING_REQUEST_BROADCASTER)
	connector.partner = partner
	return connector
}

func (connector *SIRILiteVehicleMonitoringRequestBroadcaster) RequestVehicles(filters url.Values, message *audit.BigQueryMessage) (siriLiteResponse *siri.SiriLiteResponse) {
	lineRef := filters.Get("LineRef")

	message.RequestIdentifier = filters.Get("MessageIdentifier")
	message.Lines = []string{lineRef}

	siriLiteResponse = siri.NewSiriLiteResponse()
	siriLiteResponse.Siri.ServiceDelivery.ResponseTimestamp = connector.Clock().Now()
	siriLiteResponse.Siri.ServiceDelivery.ProducerRef = connector.Partner().ProducerRef()
	siriLiteResponse.Siri.ServiceDelivery.ResponseMessageIdentifier = connector.Partner().NewResponseMessageIdentifier()
	siriLiteResponse.Siri.ServiceDelivery.RequestMessageRef = filters.Get("MessageIdentifier")

	response := siri.NewSiriLiteVehicleMonitoringDelivery()
	response.ResponseTimestamp = connector.Clock().Now()
	response.RequestMessageRef = filters.Get("MessageIdentifier")
	siriLiteResponse.Siri.ServiceDelivery.VehicleMonitoringDelivery = response

	code := model.NewCode(connector.remoteCodeSpace, lineRef)
	requestedLine, ok := connector.partner.Model().Lines().FindByCode(code)
	if !ok {
		response.ErrorCondition = &siri.ErrorCondition{
			ErrorType: "InvalidDataReferencesError",
			ErrorText: fmt.Sprintf("Line %v not found", code.Value()),
		}
		message.Status = "Error"
		message.ErrorDetails = response.ErrorCondition.ErrorText
		return
	}

	response.Status = true

	var vehicleIds []string
	vehicleJourneyRefs := make(map[string]struct{})
	lineRefs := make(map[string]struct{})

	lineIds := connector.partner.Model().Lines().FindFamily(requestedLine.Id())

	for j := range lineIds {
		vs := connector.partner.Model().Vehicles().FindByLineId(lineIds[j])
		for i := range vs {
			vehicleId, ok := vs[i].CodeWithFallback(connector.vehicleRemoteCodeSpaces)
			if !ok {
				continue
			}

			vj := vs[i].VehicleJourney()
			if vj == nil {
				continue
			}

			dvj, ok := connector.datedVehicleJourneyRef(vj)
			if !ok {
				continue
			}

			activity := siri.NewSiriLiteVehicleActivity()
			activity.RecordedAtTime = vs[i].RecordedAtTime
			activity.ValidUntilTime = vs[i].ValidUntilTime
			activity.VehicleMonitoringRef = vehicleId.Value()
			activity.MonitoredVehicleJourney.LineRef = lineRef
			activity.MonitoredVehicleJourney.PublishedLineName = requestedLine.Name
			activity.MonitoredVehicleJourney.DirectionName = vj.Attributes[siri_attributes.DirectionName]
			activity.MonitoredVehicleJourney.OriginName = vj.OriginName
			activity.MonitoredVehicleJourney.DestinationName = vj.DestinationName
			activity.MonitoredVehicleJourney.Monitored = vj.Monitored
			activity.MonitoredVehicleJourney.Bearing = vs[i].Bearing
			activity.MonitoredVehicleJourney.DriverRef = vs[i].DriverRef

			refs := vj.References.Copy()
			activity.MonitoredVehicleJourney.OriginRef = connector.handleRef("OriginRef", vj.Origin, refs)
			activity.MonitoredVehicleJourney.DestinationRef = connector.handleRef("DestinationRef", vj.Origin, refs)

			if vs[i].NextStopVisitId != model.StopVisitId("") {
				nextStopVisit, ok := connector.Partner().Model().StopVisits().Find(vs[i].NextStopVisitId)
				if ok {
					stopArea, stopAreaCode, ok := connector.stopPointRef(nextStopVisit.StopArea().Id())
					if ok {
						monitoredCall := &siri.MonitoredCall{}
						monitoredCall.StopPointRef = stopAreaCode
						monitoredCall.StopPointName = stopArea.Name
						monitoredCall.VehicleAtStop = nextStopVisit.VehicleAtStop
						monitoredCall.DestinationDisplay = nextStopVisit.Attributes["DestinationDisplay"]
						monitoredCall.ExpectedArrivalTime = nextStopVisit.Schedules.DepartureTimeFromKind([]schedules.StopVisitScheduleType{schedules.Expected})
						monitoredCall.ExpectedDepartureTime = nextStopVisit.Schedules.ArrivalTimeFromKind([]schedules.StopVisitScheduleType{schedules.Expected})
						monitoredCall.DepartureStatus = string(nextStopVisit.DepartureStatus)
						monitoredCall.Order = &nextStopVisit.PassageOrder
						monitoredCall.AimedArrivalTime = nextStopVisit.Schedules.ArrivalTimeFromKind([]schedules.StopVisitScheduleType{schedules.Aimed})
						monitoredCall.AimedDepartureTime = nextStopVisit.Schedules.DepartureTimeFromKind([]schedules.StopVisitScheduleType{schedules.Aimed})
						monitoredCall.ArrivalStatus = string(nextStopVisit.ArrivalStatus)
						monitoredCall.ActualArrivalTime = nextStopVisit.Schedules.ArrivalTimeFromKind([]schedules.StopVisitScheduleType{schedules.Actual})
						monitoredCall.ActualDepartureTime = nextStopVisit.Schedules.DepartureTimeFromKind([]schedules.StopVisitScheduleType{schedules.Actual})

						activity.MonitoredVehicleJourney.MonitoredCall = monitoredCall
					}

				}
			}

			modelDate := connector.partner.Model().Date()
			activity.MonitoredVehicleJourney.FramedVehicleJourneyRef.DataFrameRef =
				connector.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "DataFrame", Id: modelDate.String()})
			activity.MonitoredVehicleJourney.FramedVehicleJourneyRef.DatedVehicleJourneyRef = dvj

			activity.MonitoredVehicleJourney.VehicleLocation.Longitude = vs[i].Longitude
			activity.MonitoredVehicleJourney.VehicleLocation.Latitude = vs[i].Latitude

			activity.MonitoredVehicleJourney.Occupancy = vs[i].Occupancy

			response.VehicleActivity = append(response.VehicleActivity, activity)
			vehicleIds = append(vehicleIds, vehicleId.Value())
			vehicleJourneyRefs[dvj] = struct{}{}

		}
	}

	lineRefs[lineRef] = struct{}{}

	if connector.partner.PartnerSettings.SortPaylodForTest() {
		sort.Sort(siri.SortByVehicleMonitoringRef{VehicleActivities: response.VehicleActivity})
	}

	message.Lines = GetModelReferenceSlice(lineRefs)
	message.Vehicles = vehicleIds
	message.VehicleJourneys = GetModelReferenceSlice(vehicleJourneyRefs)

	return siriLiteResponse
}

func (connector *SIRILiteVehicleMonitoringRequestBroadcaster) stopPointRef(stopAreaId model.StopAreaId) (*model.StopArea, string, bool) {
	stopPointRef, ok := connector.partner.Model().StopAreas().Find(stopAreaId)
	if !ok {
		return &model.StopArea{}, "", false
	}

	stopPointRefCode, ok := stopPointRef.Code(connector.remoteCodeSpace)
	if ok {
		return stopPointRef, stopPointRefCode.Value(), true
	}
	referent, ok := stopPointRef.Referent()
	if ok {
		referentCode, ok := referent.Code(connector.remoteCodeSpace)
		if ok {
			return referent, referentCode.Value(), true
		}
	}
	return &model.StopArea{}, "", false
}

func (connector *SIRILiteVehicleMonitoringRequestBroadcaster) datedVehicleJourneyRef(vehicleJourney *model.VehicleJourney) (string, bool) {
	vehicleJourneyId, ok := vehicleJourney.CodeWithFallback(connector.vjRemoteCodeSpaces)

	var dataVehicleJourneyRef string
	if ok {
		dataVehicleJourneyRef = vehicleJourneyId.Value()
	} else {
		defaultCode, ok := vehicleJourney.Code(model.Default)
		if !ok {
			return "", false
		}
		dataVehicleJourneyRef =
			connector.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "VehicleJourney", Id: defaultCode.Value()})
	}
	return dataVehicleJourneyRef, true
}

func (connector *SIRILiteVehicleMonitoringRequestBroadcaster) handleRef(refType, origin string, references model.References) string {
	reference, ok := references.Get(refType)
	if !ok || reference.Code == nil || (refType == "DestinationRef" && connector.noDestinationRefRewritingFrom(origin)) {
		return ""
	}
	return connector.resolveStopAreaRef(reference)
}

func (connector *SIRILiteVehicleMonitoringRequestBroadcaster) noDestinationRefRewritingFrom(origin string) bool {
	ndrrf := connector.Partner().NoDestinationRefRewritingFrom()
	for _, o := range ndrrf {
		if origin == strings.TrimSpace(o) {
			return true
		}
	}
	return false
}

func (connector *SIRILiteVehicleMonitoringRequestBroadcaster) resolveStopAreaRef(reference model.Reference) string {
	stopArea, ok := connector.partner.Model().StopAreas().FindByCode(*reference.Code)
	if ok {
		obj, ok := stopArea.ReferentOrSelfCode(connector.remoteCodeSpace)
		if ok {
			return obj.Value()
		}
	}
	return connector.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "StopArea", Id: reference.GetSha1()})
}

func (factory *SIRILiteVehicleMonitoringRequestBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *SIRILiteVehicleMonitoringRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRILiteVehicleMonitoringRequestBroadcaster(partner)
}
