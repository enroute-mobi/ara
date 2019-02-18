package api

import (
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/af83/edwig/config"
	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/version"
)

type Server struct {
	model.UUIDConsumer
	model.ClockConsumer
	core.ReferentialsConsumer

	bind        string
	startedTime time.Time
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

	return &RequestData{
		Referential: requestFiller[1],
		Resource:    requestFiller[2],
		Id:          requestFiller[3],
		Action:      requestFiller[4],
	}
}

func NewServer(bind string) *Server {
	server := Server{bind: bind}
	server.startedTime = server.Clock().Now()

	return &server
}

func (server *Server) ListenAndServe() error {
	http.HandleFunc("/", server.HandleFlow)

	logger.Log.Debugf("Starting server on %s", server.bind)
	return http.ListenAndServe(server.bind, nil)
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

func (server *Server) parse(response http.ResponseWriter, request *http.Request) (*RequestData, bool) {
	path := request.URL.RequestURI()

	pathRegexp := "/([0-9a-zA-Z-_]+)(?:/([0-9a-zA-Z-_]+))?(?:/([0-9a-zA-Z-]+(?::[0-9a-zA-Z-:]+)?))?/?([0-9a-zA-Z-_]+)?"
	pattern := regexp.MustCompile(pathRegexp)
	foundStrings := pattern.FindStringSubmatch(path)
	if foundStrings == nil || foundStrings[1] == "" {
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return nil, false
	}

	requestData := NewRequestDataFromContent(foundStrings)
	requestData.Method = request.Method
	requestData.Url = request.URL.Path
	requestData.Filters = request.URL.Query()

	response.Header().Set("Content-Type", "application/json")
	response.Header().Set("Server", version.ApplicationName())
	return requestData, true
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
	return server.getToken(r) == config.Config.ApiKey
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

func (server *Server) handleRoutes(response http.ResponseWriter, request *http.Request, requestData *RequestData) {
	if requestData.Resource == "siri" {
		server.handleSIRI(response, request, requestData)
	} else if strings.HasPrefix(requestData.Referential, "_") {
		if !server.isAdmin(request) {
			http.Error(response, "Unauthorized request", http.StatusUnauthorized)
			logger.Log.Debugf("Tried to access ressource admin without autorization token \n%v", request)
			return
		}
		if requestData.Referential == "_referentials" {
			requestData.Id = requestData.Resource
		}
		server.handleControllers(response, request, requestData)
	} else {
		server.handleWithReferentialControllers(response, request, requestData)
	}
}

func (server *Server) HandleFlow(response http.ResponseWriter, request *http.Request) {
	requestData, ok := server.parse(response, request)
	if !ok {
		return
	}

	server.handleRoutes(response, request, requestData)
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

func (server *Server) handleSIRI(response http.ResponseWriter, request *http.Request, requestData *RequestData) {
	foundReferential := server.CurrentReferentials().FindBySlug(core.ReferentialSlug(requestData.Referential))

	logger.Log.Debugf("SIRI request: %v", request)

	siriHandler := NewSIRIHandler(foundReferential)
	siriHandler.serve(response, request)
}
