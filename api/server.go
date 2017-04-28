package api

import (
	"errors"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/af83/edwig/config"
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
	http.HandleFunc("/", server.HandleFlow)

	logger.Log.Debugf("Starting server on %s", server.bind)
	return http.ListenAndServe(server.bind, nil)
}

// func (server *Server) HandleFlow(response http.ResponseWriter, request *http.Request) {
// 	path := request.URL.Path
// 	pathRegexp := "/([0-9a-zA-Z-_]+)(?:/([0-9a-zA-Z-_]+))?(?:/([0-9a-zA-Z-]+(?::[0-9a-zA-Z-:]+)?))?"
// 	pattern := regexp.MustCompile(pathRegexp)
// 	foundStrings := pattern.FindStringSubmatch(path)
// 	if foundStrings == nil || foundStrings[1] == "" {
// 		http.Error(response, "Invalid request", 400)
// 		return
// 	}
// 	response.Header().Set("Content-Type", "application/json")
//
// 	if foundStrings[2] == "siri" {
// 		server.handleSIRI(response, request, foundStrings[1])
// 	} else if strings.HasPrefix(foundStrings[1], "_") {
// 		server.handleControllers(response, request, foundStrings[1], foundStrings[2])
// 	} else {
// 		server.handleWithReferentialControllers(response, request, foundStrings[1], foundStrings[2], foundStrings[3])
// 	}
// }

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

func (server *Server) checkFormat(response http.ResponseWriter, request *http.Request) ([]string, error) {
	path := request.URL.Path
	pathRegexp := "/([0-9a-zA-Z-_]+)(?:/([0-9a-zA-Z-_]+))?(?:/([0-9a-zA-Z-]+(?::[0-9a-zA-Z-:]+)?))?"
	pattern := regexp.MustCompile(pathRegexp)
	foundStrings := pattern.FindStringSubmatch(path)
	if foundStrings == nil || foundStrings[1] == "" {
		http.Error(response, "Invalid request", 400)
		return nil, errors.New("Invalid request")
	}

	response.Header().Set("Content-Type", "application/json")
	return foundStrings, nil
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
	if server.getToken(r) != config.Config.ApiKey {
		return false
	}
	return true
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

func (server *Server) handleRoutes(response http.ResponseWriter, request *http.Request, foundStrings []string) {
	if foundStrings[2] == "siri" {
		server.handleSIRI(response, request, foundStrings[1])
	} else if strings.HasPrefix(foundStrings[1], "_") {
		if !server.isAdmin(request) {
			http.Error(response, "Unauthorized request", 401)
			logger.Log.Debugf("Tried to access ressource admin without autorization token \n%s", request)
			return
		}
		server.handleControllers(response, request, foundStrings[1], foundStrings[2])
	} else {
		server.handleWithReferentialControllers(response, request, foundStrings[1], foundStrings[2], foundStrings[3])
	}
}

func (server *Server) HandleFlow(response http.ResponseWriter, request *http.Request) {
	foundStrings := []string{}

	f, _ := os.OpenFile("/tmp/salut", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
	f.WriteString(server.getToken(request) + "\n")

	foundStrings, err := server.checkFormat(response, request)

	if err != nil {
		return
	}
	server.handleRoutes(response, request, foundStrings)
}

func (server *Server) handleWithReferentialControllers(response http.ResponseWriter, request *http.Request, referential, ressource, id string) {
	foundReferential := server.CurrentReferentials().FindBySlug(core.ReferentialSlug(referential))
	if foundReferential == nil {
		http.Error(response, "Referential not found", 500)
		return
	}
	if !server.isAuth(foundReferential, request) {
		http.Error(response, "Unauthorized request", 401)
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
