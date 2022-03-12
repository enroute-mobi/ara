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

	remoteObjectidKind        string
	vehicleRemoteObjectidKind string
}

type SIRILiteVehicleMonitoringRequestBroadcasterFactory struct{}

func NewSIRILiteVehicleMonitoringRequestBroadcaster(partner *Partner) *SIRILiteVehicleMonitoringRequestBroadcaster {
	siriVehicleMonitoringRequestBroadcaster := &SIRILiteVehicleMonitoringRequestBroadcaster{
		remoteObjectidKind:        partner.RemoteObjectIDKind(SIRI_LITE_VEHICLE_MONITORING_REQUEST_BROADCASTER),
		vehicleRemoteObjectidKind: partner.VehicleRemoteObjectIDKind(SIRI_LITE_VEHICLE_MONITORING_REQUEST_BROADCASTER),
	}
	siriVehicleMonitoringRequestBroadcaster.partner = partner
	return siriVehicleMonitoringRequestBroadcaster
}

func (connector *SIRILiteVehicleMonitoringRequestBroadcaster) RequestVehicles(url string, filters url.Values, message *audit.BigQueryMessage) (siriLiteResponse *siri.SiriLiteResponse) {
	tx := connector.Partner().Referential().NewTransaction()

	logStashEvent := connector.newLogStashEvent()
	defer func() {
		tx.Close()
		audit.CurrentLogStash().WriteEvent(logStashEvent)
	}()

	lineRef := filters.Get("LineRef")

	logStashEvent["RequestURL"] = url
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
	line, ok := tx.Model().Lines().FindByObjectId(objectid)
	if !ok {
		response.ErrorCondition = &siri.ErrorCondition{
			ErrorType: "InvalidDataReferencesError",
			ErrorText: fmt.Sprintf("Line %v not found", objectid.Value()),
		}
		logSIRILiteVehicleMonitoringResponse(logStashEvent, siriLiteResponse)
		message.Status = "Error"
		message.ErrorDetails = response.ErrorCondition.ErrorText
		return
	}

	response.Status = true

	var vehicleIds []string

	for _, vehicle := range tx.Model().Vehicles().FindByLineId(line.Id()) {
		vehicleId, ok := vehicle.ObjectID(connector.vehicleRemoteObjectidKind)
		if !ok {
			continue
		}

		vj := vehicle.VehicleJourney()
		if vj == nil {
			continue
		}
		dvj, ok := connector.datedVehicleJourneyRef(vj)
		if !ok {
			continue
		}

		activity := siri.NewSiriLiteVehicleActivity()
		activity.RecordedAtTime = vehicle.RecordedAtTime
		activity.ValidUntilTime = vehicle.RecordedAtTime
		activity.VehicleMonitoringRef = vehicleId.Value()
		activity.MonitoredVehicleJourney.LineRef = lineRef
		activity.MonitoredVehicleJourney.PublishedLineName = line.Name
		activity.MonitoredVehicleJourney.DirectionName = vj.Attributes["DirectionName"]
		activity.MonitoredVehicleJourney.OriginName = vj.OriginName
		activity.MonitoredVehicleJourney.DestinationName = vj.DestinationName
		activity.MonitoredVehicleJourney.Monitored = vj.Monitored
		activity.MonitoredVehicleJourney.Bearing = vehicle.Bearing

		refs := vj.References.Copy()
		activity.MonitoredVehicleJourney.OriginRef = connector.handleRef(tx, "OriginRef", vj.Origin, refs)
		activity.MonitoredVehicleJourney.DestinationRef = connector.handleRef(tx, "DestinationRef", vj.Origin, refs)

		modelDate := tx.Model().Date()
		activity.MonitoredVehicleJourney.FramedVehicleJourneyRef.DataFrameRef =
			connector.Partner().IdentifierGenerator(idgen.DATA_FRAME_IDENTIFIER).NewIdentifier(idgen.IdentifierAttributes{Id: modelDate.String()})
		activity.MonitoredVehicleJourney.FramedVehicleJourneyRef.DatedVehicleJourneyRef = dvj

		activity.MonitoredVehicleJourney.VehicleLocation.Longitude = vehicle.Longitude
		activity.MonitoredVehicleJourney.VehicleLocation.Latitude = vehicle.Latitude

		// Delay                   *time.Time `json:",omitempty"`

		response.VehicleActivity = append(response.VehicleActivity, activity)

		vehicleIds = append(vehicleIds, vehicleId.Value())
	}

	message.Vehicles = vehicleIds

	logSIRILiteVehicleMonitoringResponse(logStashEvent, siriLiteResponse)

	return siriLiteResponse
}

func (connector *SIRILiteVehicleMonitoringRequestBroadcaster) datedVehicleJourneyRef(vehicleJourney *model.VehicleJourney) (string, bool) {
	vehicleJourneyId, ok := vehicleJourney.ObjectID(connector.remoteObjectidKind)

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

func (connector *SIRILiteVehicleMonitoringRequestBroadcaster) handleRef(tx *model.Transaction, refType, origin string, references model.References) string {
	reference, ok := references.Get(refType)
	if !ok || reference.ObjectId == nil || (refType == "DestinationRef" && connector.noDestinationRefRewritingFrom(origin)) {
		return ""
	}
	return connector.resolveStopAreaRef(tx, reference)
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

func (connector *SIRILiteVehicleMonitoringRequestBroadcaster) resolveStopAreaRef(tx *model.Transaction, reference model.Reference) string {
	stopArea, ok := tx.Model().StopAreas().FindByObjectId(*reference.ObjectId)
	if ok {
		obj, ok := stopArea.ReferentOrSelfObjectId(connector.remoteObjectidKind)
		if ok {
			return obj.Value()
		}
	}
	return connector.Partner().IdentifierGenerator(idgen.REFERENCE_STOP_AREA_IDENTIFIER).NewIdentifier(idgen.IdentifierAttributes{Id: reference.GetSha1()})
}

func (connector *SIRILiteVehicleMonitoringRequestBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "VehicleMonitoringRequestBroadcaster"
	return event
}

func (factory *SIRILiteVehicleMonitoringRequestBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *SIRILiteVehicleMonitoringRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRILiteVehicleMonitoringRequestBroadcaster(partner)
}

func logSIRILiteVehicleMonitoringResponse(logStashEvent audit.LogStashEvent, siriLiteResponse *siri.SiriLiteResponse) {

}

// func logXMLVehicleMonitoringRequest(logStashEvent audit.LogStashEvent, request *siri.XMLVehicleMonitoringRequest) {
// 	logStashEvent["siriType"] = "VehicleMonitoringResponse"
// 	logStashEvent["messageIdentifier"] = request.MessageIdentifier()
// 	logStashEvent["monitoringRef"] = request.MonitoringRef()
// 	logStashEvent["stopVisitTypes"] = request.StopVisitTypes()
// 	logStashEvent["lineRef"] = request.LineRef()
// 	logStashEvent["maximumStopVisits"] = strconv.Itoa(request.MaximumStopVisits())
// 	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
// 	logStashEvent["startTime"] = request.StartTime().String()
// 	logStashEvent["previewInterval"] = request.PreviewInterval().String()
// 	logStashEvent["requestXML"] = request.RawXML()
// }

// func logSIRIVehicleMonitoringDelivery(logStashEvent audit.LogStashEvent, delivery siri.SIRIVehicleMonitoringDelivery) {
// 	logStashEvent["requestMessageRef"] = delivery.RequestMessageRef
// 	logStashEvent["responseTimestamp"] = delivery.ResponseTimestamp.String()
// 	logStashEvent["status"] = strconv.FormatBool(delivery.Status)
// 	if !delivery.Status {
// 		logStashEvent["errorType"] = delivery.ErrorType
// 		if delivery.ErrorType == "OtherError" {
// 			logStashEvent["errorNumber"] = strconv.Itoa(delivery.ErrorNumber)
// 		}
// 		logStashEvent["errorText"] = delivery.ErrorText
// 	}
// }

// func logSIRIVehicleMonitoringResponse(logStashEvent audit.LogStashEvent, response *siri.SIRIVehicleMonitoringResponse) {
// 	logStashEvent["address"] = response.Address
// 	logStashEvent["producerRef"] = response.ProducerRef
// 	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier
// 	xml, err := response.BuildXML()
// 	if err != nil {
// 		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
// 		return
// 	}
// 	logStashEvent["responseXML"] = xml
// }
