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
	logStashEvent := make(audit.LogStashEvent)
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logXMLCheckStatusRequest(logStashEvent, request)

	response := new(siri.SIRICheckStatusResponse)
	response.Address = connector.Partner().Setting("address")
	response.ProducerRef = "Edwig"
	response.RequestMessageRef = request.MessageIdentifier()
	response.ResponseMessageIdentifier = connector.SIRIPartner().NewMessageIdentifier()
	response.Status = true
	response.ResponseTimestamp = connector.Clock().Now()
	response.ServiceStartedTime = connector.Partner().Referential().StartedAt()

	logSIRICheckStatusResponse(logStashEvent, response)

	return response, nil
}

func (factory *SIRICheckStatusServerFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("local_credential")
	return ok
}

func (factory *SIRICheckStatusServerFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRICheckStatusServer(partner)
}

func logXMLCheckStatusRequest(logStashEvent audit.LogStashEvent, request *siri.XMLCheckStatusRequest) {
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
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp.String()
	logStashEvent["serviceStartedTime"] = response.ServiceStartedTime.String()
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}
