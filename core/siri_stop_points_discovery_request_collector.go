package core

import (
	"fmt"
	"strconv"
	"strings"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	ig "bitbucket.org/enroute-mobi/ara/core/identifier_generator"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type StopPointsDiscoveryRequestCollector interface {
	RequestStopPoints()
}

type SIRIStopPointsDiscoveryRequestCollector struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	connector

	stopAreaUpdateSubscriber UpdateSubscriber
}

type SIRIStopPointsDiscoveryRequestCollectorFactory struct{}

func (factory *SIRIStopPointsDiscoveryRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopPointsDiscoveryRequestCollector(partner)
}

func (factory *SIRIStopPointsDiscoveryRequestCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfRemoteCredentials()
}

func NewSIRIStopPointsDiscoveryRequestCollector(partner *Partner) *SIRIStopPointsDiscoveryRequestCollector {
	connector := &SIRIStopPointsDiscoveryRequestCollector{}
	connector.partner = partner
	manager := partner.Referential().CollectManager()
	connector.stopAreaUpdateSubscriber = manager.BroadcastUpdateEvent

	return connector
}

func (connector *SIRIStopPointsDiscoveryRequestCollector) SetSubscriber(subscriber UpdateSubscriber) {
	connector.stopAreaUpdateSubscriber = subscriber
}

func (connector *SIRIStopPointsDiscoveryRequestCollector) broadcastUpdateEvent(event model.UpdateEvent) {
	if connector.stopAreaUpdateSubscriber != nil {
		connector.stopAreaUpdateSubscriber(event)
	}
}

func (connector *SIRIStopPointsDiscoveryRequestCollector) RequestStopPoints() {
	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	message := connector.newBQEvent()
	defer audit.CurrentBigQuery(string(connector.Partner().Referential().Slug())).WriteEvent(message)

	startTime := connector.Clock().Now()

	request := &siri.SIRIStopPointsDiscoveryRequest{
		MessageIdentifier: connector.Partner().IdentifierGenerator(ig.MESSAGE_IDENTIFIER).NewMessageIdentifier(),
		RequestorRef:      connector.Partner().RequestorRef(),
		RequestTimestamp:  startTime,
	}

	logSIRIStopPointsDiscoveryRequest(logStashEvent, message, request)

	response, err := connector.Partner().SOAPClient().StopDiscovery(request)
	logStashEvent["responseTime"] = connector.Clock().Since(startTime).String()
	message.ProcessingTime = connector.Clock().Since(startTime).Seconds()
	if err != nil {
		e := fmt.Sprintf("Error during StopDiscovery: %v", err)
		logStashEvent["status"] = "false"
		logStashEvent["errorDescription"] = e

		message.Status = "Error"
		message.ErrorDetails = e
		return
	}

	logXMLStopPointsDiscoveryResponse(logStashEvent, message, response)

	if !response.Status() {
		return
	}

	stopPointRefs := []string{}
	idKind := connector.partner.RemoteObjectIDKind()
	partner := string(connector.Partner().Slug())

	for _, annotatedStopPoint := range response.AnnotatedStopPointRefs() {
		stopPointRefs = append(stopPointRefs, annotatedStopPoint.StopPointRef())
		event := model.NewStopAreaUpdateEvent()

		event.Origin = partner
		event.ObjectId = model.NewObjectID(idKind, annotatedStopPoint.StopPointRef())
		event.Name = annotatedStopPoint.StopName()
		event.CollectedAlways = true

		connector.broadcastUpdateEvent(event)
	}

	connector.partner.RegisterDiscoveredStopAreas(stopPointRefs)
	logStashEvent["stopPointRefs"] = strings.Join(stopPointRefs, ",")
	message.StopAreas = stopPointRefs
}

func (connector *SIRIStopPointsDiscoveryRequestCollector) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      "StopDiscoveryRequest",
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(connector.partner.Slug()),
		Status:    "OK",
	}
}

func (connector *SIRIStopPointsDiscoveryRequestCollector) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "StopPointsDiscoveryRequestCollector"
	return event
}

func logSIRIStopPointsDiscoveryRequest(logStashEvent audit.LogStashEvent, message *audit.BigQueryMessage, request *siri.SIRIStopPointsDiscoveryRequest) {
	logStashEvent["siriType"] = "StopPointsDiscoveryRequest"
	logStashEvent["messageIdentifier"] = request.MessageIdentifier
	logStashEvent["requestorRef"] = request.RequestorRef
	logStashEvent["requestTimestamp"] = request.RequestTimestamp.String()
	xml, err := request.BuildXML()
	if err != nil {
		logStashEvent["requestXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["requestXML"] = xml

	message.RequestIdentifier = request.MessageIdentifier
	message.RequestRawMessage = xml
	message.RequestSize = int64(len(xml))
}

func logXMLStopPointsDiscoveryResponse(logStashEvent audit.LogStashEvent, message *audit.BigQueryMessage, response *siri.XMLStopPointsDiscoveryResponse) {
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp().String()
	logStashEvent["responseXML"] = response.RawXML()
	logStashEvent["status"] = strconv.FormatBool(response.Status())
	if !response.Status() {
		message.Status = "Error"
		logStashEvent["errorType"] = response.ErrorType()
		if response.ErrorType() == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(response.ErrorNumber())
		}
		logStashEvent["errorText"] = response.ErrorText()
		logStashEvent["errorDescription"] = response.ErrorDescription()
		message.ErrorDetails = response.ErrorString()
	}

	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
}
