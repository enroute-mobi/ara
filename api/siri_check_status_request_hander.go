package api

import (
	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SIRICheckStatusRequestHandler struct {
	xmlRequest *siri.XMLCheckStatusRequest
}

func (handler *SIRICheckStatusRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRICheckStatusRequestHandler) ConnectorType() string {
	return "siri-check-status-client"
}

func (handler *SIRICheckStatusRequestHandler) XMLResponse(connector core.Connector) string {
	logger.Log.Debugf("CheckStatus %s\n", handler.xmlRequest.MessageIdentifier())

	response := new(siri.SIRICheckStatusResponse)
	response.Address = connector.(*core.SIRICheckStatusClient).Partner().Setting("Address")
	response.ProducerRef = "Edwig"
	response.RequestMessageRef = handler.xmlRequest.MessageIdentifier()
	response.ResponseMessageIdentifier = connector.(*core.SIRICheckStatusClient).SIRIPartner().NewMessageIdentifier()
	response.Status = true
	response.ResponseTimestamp = model.DefaultClock().Now()
	response.ServiceStartedTimeStartedTime = connector.(*core.SIRICheckStatusClient).Partner().StartedAt()

	return response.BuildXML()
	// return ""
}

// func (server *Server) handleCheckStatus(w http.ResponseWriter, r *http.Request, referential string) {
// 	// Create XMLCheckStatusResponse
// 	envelope, err := siri.NewSOAPEnvelope(r.Body)
// 	if err != nil {
// 		http.Error(w, "Invalid request: can't read content", 400)
// 		return
// 	}
// 	if envelope.BodyType() != "CheckStatus" {
// 		http.Error(w, "Invalid request: not a checkstatus", 400)
// 		return
// 	}
// 	xmlRequest := siri.NewXMLCheckStatusRequest(envelope.Body())

// 	logger.Log.Debugf("CheckStatus %s\n", xmlRequest.MessageIdentifier())

// 	// Set Content-Type header and create a SIRICheckStatusResponse
// 	w.Header().Set("Content-Type", "text/xml")

// 	response := new(siri.SIRICheckStatusResponse)
// 	response.Address = strings.Join([]string{r.URL.Host, r.URL.Path}, "")
// 	response.ProducerRef = "Edwig"
// 	response.RequestMessageRef = xmlRequest.MessageIdentifier()
// 	response.ResponseMessageIdentifier = fmt.Sprintf("Edwig:ResponseMessage::%s:LOC", server.NewUUID())
// 	response.Status = true // Temp
// 	response.ResponseTimestamp = server.Clock().Now()
// 	response.ServiceStartedTime = server.startedTime

// 	// Wrap soap and send response
// 	soapEnvelope := siri.NewSOAPEnvelopeBuffer()
// 	soapEnvelope.WriteXML(response.BuildXML())

// 	_, err = soapEnvelope.WriteTo(w)
// 	if err != nil {
// 		http.Error(w, "Service internal error", 500)
// 	}
// }
