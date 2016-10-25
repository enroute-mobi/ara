package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type Server struct {
	model.UUIDConsumer
	model.ClockConsumer

	bind        string
	startedTime time.Time

	stopAreaController *StopAreaController
}

func NewServer(bind string) *Server {
	server := Server{bind: bind}
	server.startedTime = server.Clock().Now()
	server.stopAreaController = NewStopAreaController()
	return &server
}

func (server *Server) ListenAndServe(slug model.ReferentialSlug) error {
	// Temp #1852: Create a default referential
	referential := model.CurrentReferentials().New(slug)
	referential.Save()
	server.stopAreaController.SetReferential(&referential)

	http.HandleFunc(fmt.Sprintf("/%s/siri", slug), server.checkStatusHandler)
	http.HandleFunc(fmt.Sprintf("/%s/stop_areas", slug), server.stopAreaController.ServeHTTP)
	logger.Log.Debugf("Starting server on %s\n", server.bind)
	return http.ListenAndServe(server.bind, nil)
}

func (server *Server) checkStatusHandler(w http.ResponseWriter, r *http.Request) {
	// Create XMLCheckStatusResponse
	envelope, err := siri.NewSOAPEnvelope(r.Body)
	if err != nil {
		http.Error(w, "Invalid request: can't read content", 400)
		return
	}
	if envelope.BodyType() != "CheckStatus" {
		http.Error(w, "Invalid request: not a checkstatus", 400)
		return
	}
	xmlRequest := siri.NewXMLCheckStatusRequest(envelope.Body())

	logger.Log.Debugf("CheckStatus %s\n", xmlRequest.MessageIdentifier())

	// Set Content-Type header and create a SIRICheckStatusResponse
	w.Header().Set("Content-Type", "text/xml")

	response := new(siri.SIRICheckStatusResponse)
	response.Address = strings.Join([]string{r.URL.Host, r.URL.Path}, "")
	response.ProducerRef = "Edwig"
	response.RequestMessageRef = xmlRequest.MessageIdentifier()
	response.ResponseMessageIdentifier = fmt.Sprintf("Edwig:ResponseMessage::%s:LOC", server.NewUUID())
	response.Status = true // Temp
	response.ResponseTimestamp = server.Clock().Now()
	response.ServiceStartedTime = server.startedTime

	// Wrap soap and send response
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()
	soapEnvelope.WriteXML(response.BuildXML())

	_, err = soapEnvelope.WriteTo(w)
	if err != nil {
		http.Error(w, "Service internal error", 500)
	}
}
