package core

import (
	"fmt"
	"strconv"
	"strings"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
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

	siriConnector

	stopAreaUpdateSubscriber UpdateSubscriber
}

type SIRIStopPointsDiscoveryRequestCollectorFactory struct{}

func (factory *SIRIStopPointsDiscoveryRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopPointsDiscoveryRequestCollector(partner)
}

func (factory *SIRIStopPointsDiscoveryRequestCollectorFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting(REMOTE_OBJECTID_KIND)
	ok = ok && apiPartner.ValidatePresenceOfSetting(REMOTE_URL)
	ok = ok && apiPartner.ValidatePresenceOfSetting(REMOTE_CREDENTIAL)
	return ok
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
	startTime := connector.Clock().Now()

	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	request := &siri.SIRIStopPointsDiscoveryRequest{
		MessageIdentifier: connector.Partner().IdentifierGenerator(MESSAGE_IDENTIFIER).NewMessageIdentifier(),
		RequestorRef:      connector.SIRIPartner().RequestorRef(),
		RequestTimestamp:  startTime,
	}

	logSIRIStopPointsDiscoveryRequest(logStashEvent, request)

	response, err := connector.SIRIPartner().SOAPClient().StopDiscovery(request)
	logStashEvent["responseTime"] = connector.Clock().Since(startTime).String()
	if err != nil {
		logStashEvent["status"] = "false"
		logStashEvent["errorDescription"] = fmt.Sprintf("Error during StopDiscovery: %v", err)
		return
	}

	logXMLStopPointsDiscoveryResponse(logStashEvent, response)

	if !response.Status() {
		return
	}

	stopPointRefs := []string{}
	idKind := connector.partner.Setting(REMOTE_OBJECTID_KIND)
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

	logStashEvent["stopPointRefs"] = strings.Join(stopPointRefs, ",")
}

func (connector *SIRIStopPointsDiscoveryRequestCollector) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "StopPointsDiscoveryRequestCollector"
	return event
}

func logSIRIStopPointsDiscoveryRequest(logStashEvent audit.LogStashEvent, request *siri.SIRIStopPointsDiscoveryRequest) {
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
}

func logXMLStopPointsDiscoveryResponse(logStashEvent audit.LogStashEvent, response *siri.XMLStopPointsDiscoveryResponse) {
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp().String()
	logStashEvent["responseXML"] = response.RawXML()
	logStashEvent["status"] = strconv.FormatBool(response.Status())
	if !response.Status() {
		logStashEvent["errorType"] = response.ErrorType()
		if response.ErrorType() == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(response.ErrorNumber())
		}
		logStashEvent["errorText"] = response.ErrorText()
		logStashEvent["errorDescription"] = response.ErrorDescription()
	}
}
