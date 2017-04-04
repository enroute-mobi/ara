package api

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
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

func (server *Server) ListenAndServe() error {
	http.HandleFunc("/", server.APIHandler)

	logger.Log.Debugf("Starting server on %s", server.bind)
	return http.ListenAndServe(server.bind, nil)
}

func (server *Server) APIHandler(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	pathRegexp := "/([0-9a-zA-Z-_]+)(?:/([0-9a-zA-Z-_]+))?(?:/([0-9a-zA-Z-]+(?:&[0-9a-zA-Z-]+)?))?"
	pattern := regexp.MustCompile(pathRegexp)
	foundStrings := pattern.FindStringSubmatch(path)
	if foundStrings == nil || foundStrings[1] == "" {
		http.Error(response, "Invalid request", 400)
		return
	}

	response.Header().Set("Content-Type", "application/json")

	if foundStrings[2] == "siri" {
		server.handleSIRI(response, request, foundStrings[1])
	} else if strings.HasPrefix(foundStrings[1], "_") {
		server.handleControllers(response, request, foundStrings[1], foundStrings[2])
	} else {
		server.handleWithReferentialControllers(response, request, foundStrings[1], foundStrings[2], foundStrings[3])
	}
}

func (server *Server) handleControllers(response http.ResponseWriter, request *http.Request, ressource, value string) {
	newController, ok := newControllerMap[ressource]
	if !ok {
		http.Error(response, "Invalid ressource", 500)
		return
	}

	logger.Log.Debugf("%s controller request: %s", ressource[1:], request)

	controller := newController(server)
	controller.serve(response, request, value)
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

func (server *Server) handleSIRI(response http.ResponseWriter, request *http.Request, referential string) {
	foundReferential := server.CurrentReferentials().FindBySlug(core.ReferentialSlug(referential))

	logger.Log.Debugf("SIRI request: %s", request)

	siriHandler := NewSIRIHandler(foundReferential)
	siriHandler.serve(response, request)
}
