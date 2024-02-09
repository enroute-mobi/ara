package core

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type SituationExchangeRequestCollector interface {
	RequestSituationUpdate(kind, requestedId string)
}

type SIRISituationExchangeRequestCollectorFactory struct{}

type SIRISituationExchangeRequestCollector struct {
	connector

	situationUpdateSubscriber SituationUpdateSubscriber
}

func NewSIRISituationExchangeRequestCollector(partner *Partner) *SIRISituationExchangeRequestCollector {
	siriSituationExchangeRequestCollector := &SIRISituationExchangeRequestCollector{}
	siriSituationExchangeRequestCollector.partner = partner
	manager := partner.Referential().CollectManager()
	siriSituationExchangeRequestCollector.situationUpdateSubscriber = manager.BroadcastSituationUpdateEvent

	return siriSituationExchangeRequestCollector
}

func (connector *SIRISituationExchangeRequestCollector) RequestAllSituationsUpdate() {}

func (connector *SIRISituationExchangeRequestCollector) RequestSituationUpdate(kind, requestedId string) {
	message := connector.newBQEvent()
	defer audit.CurrentBigQuery(string(connector.Partner().Referential().Slug())).WriteEvent(message)

	startTime := connector.Clock().Now()

	siriSituationExchangeRequest := siri.NewSIRISituationExchangeRequest(
		connector.Partner().NewMessageIdentifier(),
		connector.Partner().RequestorRef(),
		connector.Clock().Now(),
	)

	// Check the request filter
	switch kind {
	case SITUATION_UPDATE_REQUEST_LINE:
		siriSituationExchangeRequest.LineRef = []string{requestedId}
	case SITUATION_UPDATE_REQUEST_STOP_AREA:
		siriSituationExchangeRequest.StopPointRef = []string{requestedId}
	}

	connector.logSIRISituationExchangeRequest(message, siriSituationExchangeRequest)

	xmlSituationExchangeResponse, err := connector.Partner().SIRIClient().SituationExchangeMonitoring(siriSituationExchangeRequest)
	message.ProcessingTime = connector.Clock().Since(startTime).Seconds()
	if err != nil {
		e := fmt.Sprintf("Error during GetSituationExchange: %v", err)
		message.Status = "Error"
		message.ErrorDetails = e
		return
	}

	logXMLSituationExchangeResponse(message, xmlSituationExchangeResponse)
	situationUpdateEvents := []*model.SituationUpdateEvent{}

	builder := NewSituationExchangeUpdateEventBuilder(connector.partner)
	builder.SetSituationExchangeDeliveryUpdateEvents(&situationUpdateEvents, xmlSituationExchangeResponse)

	// Log models
	message.Lines = GetModelReferenceSlice(builder.LineRefs)
	message.StopAreas = GetModelReferenceSlice(builder.MonitoringRefs)

	connector.broadcastSituationUpdateEvent(situationUpdateEvents)
}

func (connector *SIRISituationExchangeRequestCollector) SetSituationUpdateSubscriber(situationUpdateSubscriber SituationUpdateSubscriber) {
	connector.situationUpdateSubscriber = situationUpdateSubscriber
}

func (connector *SIRISituationExchangeRequestCollector) broadcastSituationUpdateEvent(event []*model.SituationUpdateEvent) {
	if connector.situationUpdateSubscriber != nil {
		connector.situationUpdateSubscriber(event)
	}
}

func (connector *SIRISituationExchangeRequestCollector) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      audit.SITUATION_EXCHANGE_REQUEST,
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(connector.partner.Slug()),
		Status:    "OK",
	}
}

func (factory *SIRISituationExchangeRequestCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfRemoteCredentials()
}

func (factory *SIRISituationExchangeRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRISituationExchangeRequestCollector(partner)
}

func (connector *SIRISituationExchangeRequestCollector) logSIRISituationExchangeRequest(message *audit.BigQueryMessage, request *siri.SIRIGetSituationExchangeRequest) {
	xml, err := request.BuildXML(connector.Partner().SIRIEnvelopeType())
	if err != nil {
		return
	}

	message.RequestIdentifier = request.MessageIdentifier
	message.RequestRawMessage = xml
	message.RequestSize = int64(len(xml))
}

func logXMLSituationExchangeResponse(message *audit.BigQueryMessage, response *sxml.XMLSituationExchangeResponse) {
	message.ResponseIdentifier = response.ResponseMessageIdentifier()

	if !response.Status() {
		message.Status = "Error"
		message.ErrorDetails = response.ErrorString()
	}
	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
}
