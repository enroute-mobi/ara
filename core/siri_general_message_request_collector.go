package core

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type GeneralMessageRequestCollector interface {
	RequestSituationUpdate(kind, requestedId string)
}

type SIRIGeneralMessageRequestCollectorFactory struct{}

type SIRIGeneralMessageRequestCollector struct {
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
		message.Lines = []string{requestedId}
	case SITUATION_UPDATE_REQUEST_STOP_AREA:
		siriGeneralMessageRequest.StopPointRef = []string{requestedId}
		message.StopAreas = []string{requestedId}
	}

	// Check the request version
	if connector.partner.GeneralMessageRequestVersion22() {
		siriGeneralMessageRequest.XsdInWsdl = true
	}

	connector.logSIRIGeneralMessageRequest(message, siriGeneralMessageRequest)

	xmlGeneralMessageResponse, err := connector.Partner().SIRIClient().SituationMonitoring(siriGeneralMessageRequest)
	message.ProcessingTime = connector.Clock().Since(startTime).Seconds()
	if err != nil {
		e := fmt.Sprintf("Error during GetGeneralMessage: %v", err)

		message.Status = "Error"
		message.ErrorDetails = e
		return
	}

	logXMLGeneralMessageResponse(message, xmlGeneralMessageResponse)
	situationUpdateEvents := []*model.SituationUpdateEvent{}
	connector.setSituationUpdateEvents(&situationUpdateEvents, xmlGeneralMessageResponse)

	connector.broadcastSituationUpdateEvent(situationUpdateEvents)
}

func (connector *SIRIGeneralMessageRequestCollector) setSituationUpdateEvents(situationEvents *[]*model.SituationUpdateEvent, xmlResponse *sxml.XMLGeneralMessageResponse) {
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

func (factory *SIRIGeneralMessageRequestCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfRemoteCredentials()
}

func (factory *SIRIGeneralMessageRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIGeneralMessageRequestCollector(partner)
}

func (connector *SIRIGeneralMessageRequestCollector) logSIRIGeneralMessageRequest(message *audit.BigQueryMessage, request *siri.SIRIGetGeneralMessageRequest) {
	xml, err := request.BuildXML(connector.Partner().SIRIEnvelopeType())
	if err != nil {
		return
	}

	message.RequestIdentifier = request.MessageIdentifier
	message.RequestRawMessage = xml
	message.RequestSize = int64(len(xml))
}

func logXMLGeneralMessageResponse(message *audit.BigQueryMessage, response *sxml.XMLGeneralMessageResponse) {
	message.ResponseIdentifier = response.ResponseMessageIdentifier()

	if !response.Status() {
		message.Status = "Error"
		message.ErrorDetails = response.ErrorString()
	}
	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
}
