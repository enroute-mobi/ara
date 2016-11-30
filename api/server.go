package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type Server struct {
	model.UUIDConsumer
	model.ClockConsumer
	core.ReferentialsConsumer

	bind        string
	startedTime time.Time
}

func NewServer(bind string) *Server {
	server := Server{bind: bind}
	server.startedTime = server.Clock().Now()

	return &server
}

func (server *Server) ListenAndServe(slug core.ReferentialSlug) error {
	// Temp #1852: Create a default referential
	referential := server.CurrentReferentials().New(slug)
	referential.Save()
	referential.Start()

	http.HandleFunc("/", server.APIHandler)

	logger.Log.Debugf("Starting server on %s", server.bind)
	return http.ListenAndServe(server.bind, nil)
}

func (server *Server) APIHandler(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	pathRegexp := "/([0-9a-zA-Z-_]+)(?:/([0-9a-zA-Z-_]+))?(?:/([0-9a-zA-Z-]+))?"
	pattern := regexp.MustCompile(pathRegexp)
	foundStrings := pattern.FindStringSubmatch(path)
	if foundStrings[1] == "" {
		http.Error(response, "Invalid request", 400)
		return
	}

	response.Header().Set("Content-Type", "application/json")

	if foundStrings[2] == "siri" {
		server.handleCheckStatus(response, request, foundStrings[1])
	} else if strings.HasPrefix(foundStrings[1], "_") {
		server.handleControllers(response, request, foundStrings[1], foundStrings[2])
	} else {
		server.handleWithReferentialControllers(response, request, foundStrings[1], foundStrings[2], foundStrings[3])
	}
}

func (server *Server) handleControllers(response http.ResponseWriter, request *http.Request, ressource, id string) {
	newController, ok := newControllerMap[ressource]
	if !ok {
		http.Error(response, "Invalid ressource", 500)
		return
	}

	logger.Log.Debugf("%s controller request: %s", ressource[1:], request)

	controller := newController(server)
	controller.serve(response, request, id)
}

func (server *Server) handleWithReferentialControllers(response http.ResponseWriter, request *http.Request, referential, ressource, id string) {
	foundReferential := server.CurrentReferentials().FindBySlug(core.ReferentialSlug(referential))
	if foundReferential == nil {
		http.Error(response, "Referential not found", 500)
		return
	}
	newController, ok := newWithReferentialControllerMap[ressource]
	if !ok {
		http.Error(response, "Invalid ressource", 500)
		return
	}

	logger.Log.Debugf("%s controller request: %s", ressource, request)

	controller := newController(foundReferential)
	controller.serve(response, request, id)
}

func (server *Server) handleCheckStatus(w http.ResponseWriter, r *http.Request, referential string) {
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
