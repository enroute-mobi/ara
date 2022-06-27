package core

import (
	"fmt"
	"net/url"
	"strings"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type VehicleMonitoringRequestBroadcaster interface {
	RequestVehicles(string, url.Values, *audit.BigQueryMessage) *siri.SiriLiteResponse
}

type SIRILiteVehicleMonitoringRequestBroadcaster struct {
	clock.ClockConsumer

	connector

	vjRemoteObjectidKinds      []string
	vehicleRemoteObjectidKinds []string
}

type SIRILiteVehicleMonitoringRequestBroadcasterFactory struct{}

func NewSIRILiteVehicleMonitoringRequestBroadcaster(partner *Partner) *SIRILiteVehicleMonitoringRequestBroadcaster {
	connector := &SIRILiteVehicleMonitoringRequestBroadcaster{
		vjRemoteObjectidKinds:      partner.VehicleJourneyRemoteObjectIDKindWithFallback(SIRI_LITE_VEHICLE_MONITORING_REQUEST_BROADCASTER),
		vehicleRemoteObjectidKinds: partner.VehicleRemoteObjectIDKindWithFallback(SIRI_LITE_VEHICLE_MONITORING_REQUEST_BROADCASTER),
	}
	connector.remoteObjectidKind = partner.RemoteObjectIDKind(SIRI_LITE_VEHICLE_MONITORING_REQUEST_BROADCASTER)
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
	siriLiteResponse.Siri.ServiceDelivery.ResponseMessageIdentifier = connector.Partner().IdentifierGenerator(idgen.RESPONSE_MESSAGE_IDENTIFIER).NewMessageIdentifier()
	siriLiteResponse.Siri.ServiceDelivery.RequestMessageRef = filters.Get("MessageIdentifier")

	response := siri.NewSiriLiteVehicleMonitoringDelivery()
	response.ResponseTimestamp = connector.Clock().Now()
	response.RequestMessageRef = filters.Get("MessageIdentifier")
	siriLiteResponse.Siri.ServiceDelivery.VehicleMonitoringDelivery = response

	objectid := model.NewObjectID(connector.remoteObjectidKind, lineRef)
	line, ok := connector.partner.Model().Lines().FindByObjectId(objectid)
	if !ok {
		response.ErrorCondition = &siri.ErrorCondition{
			ErrorType: "InvalidDataReferencesError",
			ErrorText: fmt.Sprintf("Line %v not found", objectid.Value()),
		}
		message.Status = "Error"
		message.ErrorDetails = response.ErrorCondition.ErrorText
		return
	}

	response.Status = true

	var vehicleIds []string

	vs := connector.partner.Model().Vehicles().FindByLineId(line.Id())
	for i := range vs {
		vehicleId, ok := vs[i].ObjectIDWithFallback(connector.vehicleRemoteObjectidKinds)
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

		refs := vj.References.Copy()
		activity.MonitoredVehicleJourney.OriginRef = connector.handleRef("OriginRef", vj.Origin, refs)
		activity.MonitoredVehicleJourney.DestinationRef = connector.handleRef("DestinationRef", vj.Origin, refs)

		modelDate := connector.partner.Model().Date()
		activity.MonitoredVehicleJourney.FramedVehicleJourneyRef.DataFrameRef =
			connector.Partner().IdentifierGenerator(idgen.DATA_FRAME_IDENTIFIER).NewIdentifier(idgen.IdentifierAttributes{Id: modelDate.String()})
		activity.MonitoredVehicleJourney.FramedVehicleJourneyRef.DatedVehicleJourneyRef = dvj

		activity.MonitoredVehicleJourney.VehicleLocation.Longitude = vs[i].Longitude
		activity.MonitoredVehicleJourney.VehicleLocation.Latitude = vs[i].Latitude

		response.VehicleActivity = append(response.VehicleActivity, activity)

		vehicleIds = append(vehicleIds, vehicleId.Value())
	}

	message.Vehicles = vehicleIds

	return siriLiteResponse
}

func (connector *SIRILiteVehicleMonitoringRequestBroadcaster) datedVehicleJourneyRef(vehicleJourney *model.VehicleJourney) (string, bool) {
	vehicleJourneyId, ok := vehicleJourney.ObjectIDWithFallback(connector.vjRemoteObjectidKinds)

	var dataVehicleJourneyRef string
	if ok {
		dataVehicleJourneyRef = vehicleJourneyId.Value()
	} else {
		defaultObjectID, ok := vehicleJourney.ObjectID("_default")
		if !ok {
			return "", false
		}
		dataVehicleJourneyRef =
			connector.Partner().IdentifierGenerator(idgen.REFERENCE_IDENTIFIER).NewIdentifier(idgen.IdentifierAttributes{Type: "VehicleJourney", Id: defaultObjectID.Value()})
	}
	return dataVehicleJourneyRef, true
}

func (connector *SIRILiteVehicleMonitoringRequestBroadcaster) handleRef(refType, origin string, references model.References) string {
	reference, ok := references.Get(refType)
	if !ok || reference.ObjectId == nil || (refType == "DestinationRef" && connector.noDestinationRefRewritingFrom(origin)) {
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
	stopArea, ok := connector.partner.Model().StopAreas().FindByObjectId(*reference.ObjectId)
	if ok {
		obj, ok := stopArea.ReferentOrSelfObjectId(connector.remoteObjectidKind)
		if ok {
			return obj.Value()
		}
	}
	return connector.Partner().IdentifierGenerator(idgen.REFERENCE_STOP_AREA_IDENTIFIER).NewIdentifier(idgen.IdentifierAttributes{Id: reference.GetSha1()})
}

func (factory *SIRILiteVehicleMonitoringRequestBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *SIRILiteVehicleMonitoringRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRILiteVehicleMonitoringRequestBroadcaster(partner)
}
