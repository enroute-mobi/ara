package core

import (
	"fmt"
	"strconv"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type GeneralMessageRequestCollector interface {
	RequestSituationUpdate(kind, requestedId string)
}

type SIRIGeneralMessageRequestCollectorFactory struct{}

type SIRIGeneralMessageRequestCollector struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	connector

	situationUpdateSubscriber SituationUpdateSubscriber
}

func NewSIRIGeneralMessageRequestCollector(partner *Partner) *SIRIGeneralMessageRequestCollector {
	siriGeneralMessageRequestCollector := &SIRIGeneralMessageRequestCollector{}
	siriGeneralMessageRequestCollector.partner = partner
	manager := partner.Referential().CollectManager()
	siriGeneralMessageRequestCollector.situationUpdateSubscriber = manager.BroadcastSituationUpdateEvent

	return siriGeneralMessageRequestCollector
}

func (connector *SIRIGeneralMessageRequestCollector) RequestAllSituationsUpdate() {}

func (connector *SIRIGeneralMessageRequestCollector) RequestSituationUpdate(kind, requestedId string) {
	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	message := connector.newBQEvent()
	defer audit.CurrentBigQuery(string(connector.Partner().Referential().Slug())).WriteEvent(message)

	startTime := connector.Clock().Now()

	siriGeneralMessageRequest := &siri.SIRIGetGeneralMessageRequest{
		RequestorRef: connector.Partner().RequestorRef(),
	}
	siriGeneralMessageRequest.MessageIdentifier = connector.Partner().NewMessageIdentifier()
	siriGeneralMessageRequest.RequestTimestamp = connector.Clock().Now()

	// Check the request filter
	switch kind {
	case SITUATION_UPDATE_REQUEST_LINE:
		siriGeneralMessageRequest.LineRef = []string{requestedId}
		logStashEvent["lineRef"] = requestedId
		message.Lines = []string{requestedId}
	case SITUATION_UPDATE_REQUEST_STOP_AREA:
		siriGeneralMessageRequest.StopPointRef = []string{requestedId}
		logStashEvent["stopPointRef"] = requestedId
		message.StopAreas = []string{requestedId}
	}

	// Check the request version
	if connector.partner.GeneralMessageRequestVersion22() {
		siriGeneralMessageRequest.XsdInWsdl = true
	}

	logSIRIGeneralMessageRequest(logStashEvent, message, siriGeneralMessageRequest)

	xmlGeneralMessageResponse, err := connector.Partner().SIRIClient().SituationMonitoring(siriGeneralMessageRequest)
	logStashEvent["responseTime"] = connector.Clock().Since(startTime).String()
	message.ProcessingTime = connector.Clock().Since(startTime).Seconds()
	if err != nil {
		e := fmt.Sprintf("Error during GetGeneralMessage: %v", err)
		logStashEvent["status"] = "false"
		logStashEvent["errorDescription"] = e

		message.Status = "Error"
		message.ErrorDetails = e
		return
	}

	logXMLGeneralMessageResponse(logStashEvent, message, xmlGeneralMessageResponse)
	situationUpdateEvents := []*model.SituationUpdateEvent{}
	connector.setSituationUpdateEvents(&situationUpdateEvents, xmlGeneralMessageResponse)

	connector.broadcastSituationUpdateEvent(situationUpdateEvents)
}

func (connector *SIRIGeneralMessageRequestCollector) setSituationUpdateEvents(situationEvents *[]*model.SituationUpdateEvent, xmlResponse *siri.XMLGeneralMessageResponse) {
	builder := NewGeneralMessageUpdateEventBuilder(connector.partner)
	builder.SetGeneralMessageResponseUpdateEvents(situationEvents, xmlResponse)
}

func (connector *SIRIGeneralMessageRequestCollector) SetSituationUpdateSubscriber(situationUpdateSubscriber SituationUpdateSubscriber) {
	connector.situationUpdateSubscriber = situationUpdateSubscriber
}

func (connector *SIRIGeneralMessageRequestCollector) broadcastSituationUpdateEvent(event []*model.SituationUpdateEvent) {
	if connector.situationUpdateSubscriber != nil {
		connector.situationUpdateSubscriber(event)
	}
}

func (connector *SIRIGeneralMessageRequestCollector) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      "GeneralMessageRequest",
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(connector.partner.Slug()),
		Status:    "OK",
	}
}

func (connector *SIRIGeneralMessageRequestCollector) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "GeneralMessageRequestCollector"
	return event
}

func (factory *SIRIGeneralMessageRequestCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfRemoteCredentials()
}

func (factory *SIRIGeneralMessageRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIGeneralMessageRequestCollector(partner)
}

func logSIRIGeneralMessageRequest(logStashEvent audit.LogStashEvent, message *audit.BigQueryMessage, request *siri.SIRIGetGeneralMessageRequest) {
	logStashEvent["siriType"] = "GeneralMessageRequest"
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

func logXMLGeneralMessageResponse(logStashEvent audit.LogStashEvent, message *audit.BigQueryMessage, response *siri.XMLGeneralMessageResponse) {
	message.ResponseIdentifier = response.ResponseMessageIdentifier()

	logStashEvent["address"] = response.Address()
	logStashEvent["producerRef"] = response.ProducerRef()
	logStashEvent["requestMessageRef"] = response.RequestMessageRef()
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier()
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp().String()
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
	logStashEvent["responseXML"] = response.RawXML()
	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
}
