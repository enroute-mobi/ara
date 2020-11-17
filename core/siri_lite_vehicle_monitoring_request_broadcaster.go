package core

import (
	"fmt"
	"net/url"
	"strings"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type VehicleMonitoringRequestBroadcaster interface {
	RequestVehicles(string, url.Values) *siri.SiriLiteResponse
}

type SIRILiteVehicleMonitoringRequestBroadcaster struct {
	clock.ClockConsumer

	BaseConnector

	remoteObjectidKind string
}

type SIRILiteVehicleMonitoringRequestBroadcasterFactory struct{}

func NewSIRILiteVehicleMonitoringRequestBroadcaster(partner *Partner) *SIRILiteVehicleMonitoringRequestBroadcaster {
	siriVehicleMonitoringRequestBroadcaster := &SIRILiteVehicleMonitoringRequestBroadcaster{
		remoteObjectidKind: partner.RemoteObjectIDKind(SIRI_LITE_VEHICLE_MONITORING_REQUEST_BROADCASTER),
	}
	siriVehicleMonitoringRequestBroadcaster.partner = partner
	return siriVehicleMonitoringRequestBroadcaster
}

func (connector *SIRILiteVehicleMonitoringRequestBroadcaster) RequestVehicles(url string, filters url.Values) (siriLiteResponse *siri.SiriLiteResponse) {
	tx := connector.Partner().Referential().NewTransaction()

	logStashEvent := connector.newLogStashEvent()
	defer func() {
		tx.Close()
		audit.CurrentLogStash().WriteEvent(logStashEvent)
	}()

	logStashEvent["RequestURL"] = url

	siriLiteResponse = siri.NewSiriLiteResponse()
	siriLiteResponse.Siri.ServiceDelivery.ResponseTimestamp = connector.Clock().Now()
	siriLiteResponse.Siri.ServiceDelivery.ProducerRef = connector.Partner().ProducerRef()
	siriLiteResponse.Siri.ServiceDelivery.ResponseMessageIdentifier = connector.Partner().IdentifierGenerator(RESPONSE_MESSAGE_IDENTIFIER).NewMessageIdentifier()
	siriLiteResponse.Siri.ServiceDelivery.RequestMessageRef = filters.Get("MessageIdentifier")

	response := siri.NewSiriLiteVehicleMonitoringDelivery()
	response.ResponseTimestamp = connector.Clock().Now()
	response.RequestMessageRef = filters.Get("MessageIdentifier")
	siriLiteResponse.Siri.ServiceDelivery.VehicleMonitoringDelivery = response

	lineRef := filters.Get("LineRef")
	objectid := model.NewObjectID(connector.remoteObjectidKind, lineRef)
	line, ok := tx.Model().Lines().FindByObjectId(objectid)
	if !ok {
		response.ErrorCondition = &siri.ErrorCondition{
			ErrorType: "InvalidDataReferencesError",
			ErrorText: fmt.Sprintf("Line %v not found", objectid.Value()),
		}
		logSIRILiteVehicleMonitoringResponse(logStashEvent, siriLiteResponse)
		return
	}

	response.Status = true

	for _, vehicle := range tx.Model().Vehicles().FindByLineId(line.Id()) {
		vehicleId, ok := vehicle.ObjectID(connector.remoteObjectidKind)
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
			connector.Partner().IdentifierGenerator(DATA_FRAME_IDENTIFIER).NewIdentifier(IdentifierAttributes{Id: modelDate.String()})
		activity.MonitoredVehicleJourney.FramedVehicleJourneyRef.DatedVehicleJourneyRef = dvj

		activity.MonitoredVehicleJourney.VehicleLocation.Longitude = vehicle.Longitude
		activity.MonitoredVehicleJourney.VehicleLocation.Latitude = vehicle.Latitude

		// Delay                   *time.Time `json:",omitempty"`

		response.VehicleActivity = append(response.VehicleActivity, activity)
	}

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
			connector.Partner().IdentifierGenerator(REFERENCE_IDENTIFIER).NewIdentifier(IdentifierAttributes{Type: "VehicleJourney", Default: defaultObjectID.Value()})
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
	return connector.Partner().IdentifierGenerator(REFERENCE_STOP_AREA_IDENTIFIER).NewIdentifier(IdentifierAttributes{Default: reference.GetSha1()})
}

func (connector *SIRILiteVehicleMonitoringRequestBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "VehicleMonitoringRequestBroadcaster"
	return event
}

func (factory *SIRILiteVehicleMonitoringRequestBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting(REMOTE_OBJECTID_KIND)
	ok = ok && apiPartner.ValidatePresenceOfLocalCredentials()
	return ok
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
