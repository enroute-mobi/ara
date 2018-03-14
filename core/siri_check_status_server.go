package core

import (
	"fmt"
	"strconv"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type CheckStatusServer interface {
	CheckStatus(*siri.XMLCheckStatusRequest) (*siri.SIRICheckStatusResponse, error)
}

type SIRICheckStatusServer struct {
	model.ClockConsumer

	siriConnector
}

type SIRICheckStatusServerFactory struct{}

func NewSIRICheckStatusServer(partner *Partner) *SIRICheckStatusServer {
	siriCheckStatusServer := &SIRICheckStatusServer{}
	siriCheckStatusServer.partner = partner
	return siriCheckStatusServer
}

func (connector *SIRICheckStatusServer) CheckStatus(request *siri.XMLCheckStatusRequest) (*siri.SIRICheckStatusResponse, error) {
	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logXMLCheckStatusRequest(logStashEvent, request)

	response := &siri.SIRICheckStatusResponse{
		Address:                   connector.Partner().Address(),
		ProducerRef:               connector.Partner().ProducerRef(),
		RequestMessageRef:         request.MessageIdentifier(),
		ResponseMessageIdentifier: connector.SIRIPartner().IdentifierGenerator("response_message_identifier").NewMessageIdentifier(),
		Status:             true,
		ResponseTimestamp:  connector.Clock().Now(),
		ServiceStartedTime: connector.Partner().Referential().StartedAt(),
	}

	logSIRICheckStatusResponse(logStashEvent, response)

	return response, nil
}

func (connector *SIRICheckStatusServer) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "CheckStatusServer"
	return event
}

func (factory *SIRICheckStatusServerFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("local_credential")
	return ok
}

func (factory *SIRICheckStatusServerFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRICheckStatusServer(partner)
}

func logXMLCheckStatusRequest(logStashEvent audit.LogStashEvent, request *siri.XMLCheckStatusRequest) {
	logStashEvent["siriType"] = "CheckStatusResponse"
	logStashEvent["messageIdentifier"] = request.MessageIdentifier()
	logStashEvent["requestorRef"] = request.RequestorRef()
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
	logStashEvent["requestXML"] = request.RawXML()
}

func logSIRICheckStatusResponse(logStashEvent audit.LogStashEvent, response *siri.SIRICheckStatusResponse) {
	logStashEvent["address"] = response.Address
	logStashEvent["producerRef"] = response.ProducerRef
	logStashEvent["requestMessageRef"] = response.RequestMessageRef
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier
	logStashEvent["status"] = strconv.FormatBool(response.Status)
	if !response.Status {
		logStashEvent["errorType"] = response.ErrorType
		if response.ErrorType == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(response.ErrorNumber)
		}
		logStashEvent["errorText"] = response.ErrorText
	}
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp.String()
	logStashEvent["serviceStartedTime"] = response.ServiceStartedTime.String()
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}
