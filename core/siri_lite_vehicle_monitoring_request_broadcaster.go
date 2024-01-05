package core

import (
	"fmt"
	"net/url"
	"sort"
	"strings"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
)

type VehicleMonitoringLiteRequestBroadcaster interface {
	RequestVehicles(string, url.Values, *audit.BigQueryMessage) *siri.SiriLiteResponse
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

func (connector *SIRILiteVehicleMonitoringRequestBroadcaster) RequestVehicles(url string, filters url.Values, message *audit.BigQueryMessage) (siriLiteResponse *siri.SiriLiteResponse) {
	lineRef := filters.Get("LineRef")

	message.RequestIdentifier = filters.Get("MessageIdentifier")
	message.Lines = []string{lineRef}

	siriLiteResponse = siri.NewSiriLiteResponse()
	siriLiteResponse.Siri.ServiceDelivery.ResponseTimestamp = connector.Clock().Now()
	siriLiteResponse.Siri.ServiceDelivery.ProducerRef = connector.Partner().ProducerRef()
	siriLiteResponse.Siri.ServiceDelivery.ResponseMessageIdentifier = connector.Partner().ResponseMessageIdentifierGenerator().NewMessageIdentifier()
	siriLiteResponse.Siri.ServiceDelivery.RequestMessageRef = filters.Get("MessageIdentifier")

	response := siri.NewSiriLiteVehicleMonitoringDelivery()
	response.ResponseTimestamp = connector.Clock().Now()
	response.RequestMessageRef = filters.Get("MessageIdentifier")
	siriLiteResponse.Siri.ServiceDelivery.VehicleMonitoringDelivery = response

	code := model.NewCode(connector.remoteCodeSpace, lineRef)
	line, ok := connector.partner.Model().Lines().FindByCode(code)
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

	vs := connector.partner.Model().Vehicles().FindByLineId(line.Id())
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
		activity.MonitoredVehicleJourney.PublishedLineName = line.Name
		activity.MonitoredVehicleJourney.DirectionName = vj.Attributes["DirectionName"]
		activity.MonitoredVehicleJourney.OriginName = vj.OriginName
		activity.MonitoredVehicleJourney.DestinationName = vj.DestinationName
		activity.MonitoredVehicleJourney.Monitored = vj.Monitored
		activity.MonitoredVehicleJourney.Bearing = vs[i].Bearing
		activity.MonitoredVehicleJourney.DriverRef = vs[i].DriverRef

		refs := vj.References.Copy()
		activity.MonitoredVehicleJourney.OriginRef = connector.handleRef("OriginRef", vj.Origin, refs)
		activity.MonitoredVehicleJourney.DestinationRef = connector.handleRef("DestinationRef", vj.Origin, refs)

		modelDate := connector.partner.Model().Date()
		activity.MonitoredVehicleJourney.FramedVehicleJourneyRef.DataFrameRef =
			connector.Partner().DataFrameIdentifierGenerator().NewIdentifier(idgen.IdentifierAttributes{Id: modelDate.String()})
		activity.MonitoredVehicleJourney.FramedVehicleJourneyRef.DatedVehicleJourneyRef = dvj

		activity.MonitoredVehicleJourney.VehicleLocation.Longitude = vs[i].Longitude
		activity.MonitoredVehicleJourney.VehicleLocation.Latitude = vs[i].Latitude

		activity.MonitoredVehicleJourney.Occupancy = vs[i].Occupancy

		response.VehicleActivity = append(response.VehicleActivity, activity)

		vehicleIds = append(vehicleIds, vehicleId.Value())
		vehicleJourneyRefs[dvj] = struct{}{}
		lineRefs[lineRef] = struct{}{}
	}

	if connector.partner.PartnerSettings.SortPaylodForTest() {
		sort.Sort(siri.SortByVehicleMonitoringRef{VehicleActivities: response.VehicleActivity})
	}

	message.Lines = GetModelReferenceSlice(lineRefs)
	message.Vehicles = vehicleIds
	message.VehicleJourneys = GetModelReferenceSlice(vehicleJourneyRefs)

	return siriLiteResponse
}

func (connector *SIRILiteVehicleMonitoringRequestBroadcaster) datedVehicleJourneyRef(vehicleJourney *model.VehicleJourney) (string, bool) {
	vehicleJourneyId, ok := vehicleJourney.CodeWithFallback(connector.vjRemoteCodeSpaces)

	var dataVehicleJourneyRef string
	if ok {
		dataVehicleJourneyRef = vehicleJourneyId.Value()
	} else {
		defaultCode, ok := vehicleJourney.Code("_default")
		if !ok {
			return "", false
		}
		dataVehicleJourneyRef =
			connector.Partner().ReferenceIdentifierGenerator().NewIdentifier(idgen.IdentifierAttributes{Type: "VehicleJourney", Id: defaultCode.Value()})
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
	return connector.Partner().ReferenceStopAreaIdentifierGenerator().NewIdentifier(idgen.IdentifierAttributes{Id: reference.GetSha1()})
}

func (factory *SIRILiteVehicleMonitoringRequestBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *SIRILiteVehicleMonitoringRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRILiteVehicleMonitoringRequestBroadcaster(partner)
}
