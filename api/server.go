package api

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/config"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/monitoring"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"bitbucket.org/enroute-mobi/ara/version"
)

type Server struct {
	uuid.UUIDConsumer
	clock.ClockConsumer
	core.ReferentialsConsumer

	srv         *http.Server
	bind        string
	startedTime time.Time
	apiKey      string
}

type SIRIRequestData struct {
	Filters     url.Values
	Referential string
	Request     string
	Url         string
}

type RequestData struct {
	Filters     url.Values
	Body        []byte
	Method      string
	Referential string
	Resource    string
	Id          string
	Action      string
	Url         string
}

func NewRequestDataFromContent(params []string) *RequestData {
	requestFiller := make([]string, 15)

	copy(requestFiller, params)

	pathRegexp := "([0-9a-zA-Z-]+(?::[0-9a-zA-Z-:]+)?)?(?:/([0-9a-zA-Z-_]+))?"
	pattern := regexp.MustCompile(pathRegexp)
	foundStrings := pattern.FindStringSubmatch(requestFiller[3])

	return &RequestData{
		Referential: requestFiller[1],
		Resource:    requestFiller[2],
		Id:          foundStrings[1],
		Action:      foundStrings[2],
	}
}

func NewSIRIRequestDataFromContent(params []string) (*SIRIRequestData, bool) {
	requestFiller := make([]string, 15)

	copy(requestFiller, params)

	requestData := &SIRIRequestData{
		Referential: requestFiller[1],
	}

	if requestFiller[3] != "" {
		pathRegexp := "v2.0/([a-z-]+).json"
		pattern := regexp.MustCompile(pathRegexp)
		foundStrings := pattern.FindStringSubmatch(requestFiller[3])
		if len(foundStrings) == 0 {
			return nil, false
		}
		requestData.Request = foundStrings[1]
	}

	return requestData, true
}

func NewServer(bind string) *Server {
	server := Server{bind: bind}
	server.startedTime = server.Clock().Now()

	server.apiKey = config.Config.ApiKey

	return &server
}

func (server *Server) ListenAndServe() error {
	server.srv = &http.Server{
		Addr:         server.bind,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 60 * time.Second,
	}
	http.HandleFunc("/", server.HandleFlow)

	logger.Log.Debugf("Starting server on %s", server.bind)
	return server.srv.ListenAndServe()
}

func (server *Server) handleControllers(response http.ResponseWriter, request *http.Request, requestData *RequestData) {
	newController, ok := newControllerMap[requestData.Referential]
	if !ok {
		http.Error(response, "Invalid ressource", http.StatusBadRequest)
		return
	}

	logger.Log.Debugf("%s controller request: %v", requestData.Resource, request)

	controller := newController(server)
	controller.serve(response, request, requestData)
}

func (server *Server) getToken(r *http.Request) string {
	const prefix = "Token"
	auth := r.Header.Get("Authorization")

	if !strings.HasPrefix(auth, prefix) {
		return ""
	}
	s := strings.IndexByte(auth, '=')

	return auth[s+1:]
}

func (server *Server) isAdmin(r *http.Request) bool {
	return server.getToken(r) == server.apiKey
}

func (server *Server) isAuth(referential *core.Referential, request *http.Request) bool {
	authToken := server.getToken(request)

	if authToken == "" {
		return false
	}

	for _, token := range referential.Tokens {
		if authToken == token {
			return true
		}
	}
	return false
}

func (server *Server) HandleFlow(response http.ResponseWriter, request *http.Request) {
	defer monitoring.HandleHttpPanic(response)

	path := request.URL.RequestURI()
	pathRegexp := "/([0-9a-zA-Z-_]+)(?:/([0-9a-zA-Z-_]+))?(?:/([/0-9a-zA-Z-_.:]+))?"
	pattern := regexp.MustCompile(pathRegexp)
	foundStrings := pattern.FindStringSubmatch(path)
	if foundStrings == nil || foundStrings[1] == "" {
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}
	response.Header().Set("Server", version.ApplicationName())

	if foundStrings[2] == "siri" {
		requestData, ok := NewSIRIRequestDataFromContent(foundStrings)
		if !ok {
			http.Error(response, "Invalid request", http.StatusBadRequest)
			return
		}
		requestData.Filters = request.URL.Query()
		requestData.Url = request.URL.RequestURI()
		server.handleSIRI(response, request, requestData)
		return
	}

	if foundStrings[2] == "push" {
		server.handlePush(response, request, foundStrings[1])
		return
	}

	if foundStrings[2] == "gtfs" {
		server.handleGtfs(response, request, foundStrings[1], foundStrings[3])
		return
	}

	requestData := NewRequestDataFromContent(foundStrings)
	requestData.Method = request.Method
	requestData.Url = request.URL.Path
	requestData.Filters = request.URL.Query()

	response.Header().Set("Content-Type", "application/json")

	if strings.HasPrefix(requestData.Referential, "_") {
		if requestData.Referential != "_status" && !server.isAdmin(request) {
			http.Error(response, "Unauthorized request", http.StatusUnauthorized)
			logger.Log.Debugf("Tried to access ressource admin without autorization token:\n%v", request)
			return
		}
		if requestData.Referential == "_referentials" {
			requestData.Action = requestData.Id
			requestData.Id = requestData.Resource
		}
		server.handleControllers(response, request, requestData)
		return
	}

	server.handleWithReferentialControllers(response, request, requestData)
}

func (server *Server) handleWithReferentialControllers(response http.ResponseWriter, request *http.Request, requestData *RequestData) {

	foundReferential := server.CurrentReferentials().FindBySlug(core.ReferentialSlug(requestData.Referential))
	if foundReferential == nil {
		http.Error(response, "Referential not found", http.StatusNotFound)
		return
	}
	if !server.isAuth(foundReferential, request) {
		http.Error(response, "Unauthorized request", http.StatusUnauthorized)
		return
	}
	newController, ok := newWithReferentialControllerMap[requestData.Resource]
	if !ok {
		http.Error(response, "Invalid ressource", http.StatusBadRequest)
		return
	}

	logger.Log.Debugf("%s controller request: %v", requestData.Resource, request)

	controller := newController(foundReferential)
	controller.serve(response, request, requestData)
}

func (server *Server) handleSIRI(response http.ResponseWriter, request *http.Request, requestData *SIRIRequestData) {
	foundReferential := server.CurrentReferentials().FindBySlug(core.ReferentialSlug(requestData.Referential))

	logger.Log.Debugf("SIRI request: %v", request)

	if requestData.Request == "" {
		siriHandler := NewSIRIHandler(foundReferential)
		siriHandler.serve(response, request)
		return
	}

	siriHandler := NewSIRILiteHandler(foundReferential, server.getToken(request))
	siriHandler.serve(response, request, requestData)
}

func (server *Server) handlePush(response http.ResponseWriter, request *http.Request, referential string) {
	foundReferential := server.CurrentReferentials().FindBySlug(core.ReferentialSlug(referential))
	if foundReferential == nil {
		http.Error(response, "Referential not found", http.StatusNotFound)
		return
	}

	logger.Log.Debugf("Push request: %v", request)

	pushHandler := NewPushHandler(foundReferential, server.getToken(request))
	pushHandler.serve(response, request)
}

func (server *Server) handleGtfs(response http.ResponseWriter, request *http.Request, referential, resource string) {
	foundReferential := server.CurrentReferentials().FindBySlug(core.ReferentialSlug(referential))
	if foundReferential == nil {
		http.Error(response, "Referential not found", http.StatusNotFound)
		return
	}

	logger.Log.Debugf("Gtfs request: %v", request)

	gtfsHandler := NewGtfsHandler(foundReferential, server.getToken(request))
	gtfsHandler.serve(response, request, resource)
}
