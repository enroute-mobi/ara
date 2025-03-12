package api

import (
	"net/http"
	"slices"
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

func NewServer(bind string) *Server {
	server := Server{bind: bind}
	server.startedTime = server.Clock().Now()

	server.apiKey = config.Config.ApiKey

	return &server
}

func (server *Server) ListenAndServe() error {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /{referential_slug}/graphql", server.handleGraphql)
	mux.HandleFunc("POST /{referential_slug}/push", server.handlePush)

	mux.HandleFunc("GET /{referential_slug}/gtfs/{resource}", server.handleGtfs)

	mux.HandleFunc("POST /{referential_slug}/siri", server.HandleSIRI)
	mux.HandleFunc("GET /{referential_slug}/siri/v2.0/{resource}", server.handleSIRILite)

	mux.HandleFunc("GET /_status", server.handleStatus)
	mux.HandleFunc("GET /_time", server.handleTimeGet)
	mux.HandleFunc("POST /_time/advance", server.handleTimeAdvance)

	mux.HandleFunc("GET /_referentials", server.handleReferentialIndex)
	mux.HandleFunc("POST /_referentials", server.handleReferentialCreate)

	// To avoid overlap between /{referential_slug}/gtfs and  /_referentials/{id} and /{referential_slug}
	mux.HandleFunc("GET /{referential_slug}/{model}", server.handleReferentialGet)

	mux.HandleFunc("PUT /_referentials/{id}", server.handleReferentialUpdate)
	mux.HandleFunc("DELETE /_referentials/{id}", server.handleReferentialDelete)
	mux.HandleFunc("POST /_referentials/save", server.handleReferentialSave)
	mux.HandleFunc("POST /_referentials/reload/{id}", server.handleReferentialReload)

	mux.HandleFunc("GET /{referential_slug}/{model}/{id}", server.handleReferentialModelShow)
	mux.HandleFunc("POST /{referential_slug}/{model}", server.handleReferentialModelCreate)
	mux.HandleFunc("PUT /{referential_slug}/{model}/{id}", server.handleReferentialModelUpdate)
	mux.HandleFunc("DELETE /{referential_slug}/{model}/{id}", server.handleReferentialModelDelete)
	mux.HandleFunc("POST /{referential_slug}/import", server.handleReferentialImport)

	mux.HandleFunc("POST /{referential_slug}/partners/save", server.handleReferentialPartnerSave)

	mux.HandleFunc("GET /{referential_slug}/partners/{id}/subscriptions", server.handlePartnerSubscriptionsIndex)
	mux.HandleFunc("POST /{referential_slug}/partners/{id}/subscriptions", server.handlePartnerSubscriptionsCreate)
	mux.HandleFunc("DELETE /{referential_slug}/partners/{id}/subscriptions/{subscription_id}", server.handlePartnerSubscriptionsDelete)

	server.srv = &http.Server{
		Handler:      mux,
		Addr:         server.bind,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	logger.Log.Debugf("Starting server on %s", server.bind)
	return server.srv.ListenAndServe()
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

func (server *Server) isAuthForImport(referential *core.Referential, request *http.Request) bool {
	authToken := server.getToken(request)

	if authToken == "" {
		return false
	}

	for _, token := range referential.Tokens {
		if authToken == token {
			return true
		}
	}

	return slices.Contains(referential.ImportTokens, authToken)
}

func (server *Server) referentialSetup(response http.ResponseWriter, request *http.Request) (controller *ReferentialController) {
	defer monitoring.HandleHttpPanic(response)

	if !server.isAdmin(request) {
		http.Error(response, "Unauthorized request", http.StatusUnauthorized)
		logger.Log.Debugf("Tried to access ressource admin without autorization token:\n%v", request)
		return
	}

	response.Header().Set("Server", version.ApplicationName())
	response.Header().Set("Content-Type", "application/json")

	controller = NewReferentialController(server)
	return
}

func (server *Server) referentialModelSetup(response http.ResponseWriter, request *http.Request) (controller RestfulResource, model string) {
	defer monitoring.HandleHttpPanic(response)

	referentialSlug := request.PathValue("referential_slug")
	foundReferential := server.CurrentReferentials().FindBySlug(core.ReferentialSlug(referentialSlug))
	if foundReferential == nil {
		http.Error(response, "Referential not found", http.StatusNotFound)
		return
	}

	if !server.isAuth(foundReferential, request) {
		http.Error(response, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	model = request.PathValue("model")
	newController, ok := newWithReferentialControllerMap[model]
	if !ok {
		http.Error(response, "Invalid ressource", http.StatusBadRequest)
		return
	}

	response.Header().Set("Server", version.ApplicationName())
	response.Header().Set("Content-Type", "application/json")

	controller = newController(foundReferential)

	return
}

func (server *Server) handleReferentialModelIndex(response http.ResponseWriter, request *http.Request) {
	controller, model := server.referentialModelSetup(response, request)
	if controller == nil {
		return
	}

	logger.Log.Debugf("%s controller Index request: %v", model, request)

	params := request.URL.Query()
	controller.Index(response, params)
}

func (server *Server) handleReferentialModelShow(response http.ResponseWriter, request *http.Request) {
	controller, model := server.referentialModelSetup(response, request)
	if controller == nil {
		return
	}

	logger.Log.Debugf("%s controller Show request: %v", model, request)

	id := request.PathValue("id")
	controller.Show(response, id)
}

func (server *Server) handleReferentialModelCreate(response http.ResponseWriter, request *http.Request) {
	controller, model := server.referentialModelSetup(response, request)
	if controller == nil {
		return
	}

	logger.Log.Debugf("%s controller Delete request: %v", model, request)

	body := getRequestBody(response, request)
	if body == nil {
		return
	}

	logger.Log.Debugf("%s controller Create request: %v", model, request)

	controller.Create(response, body)
}

func (server *Server) handleReferentialModelUpdate(response http.ResponseWriter, request *http.Request) {
	controller, model := server.referentialModelSetup(response, request)
	if controller == nil {
		return
	}

	body := getRequestBody(response, request)
	if body == nil {
		return
	}

	logger.Log.Debugf("%s controller Update request: %v", model, request)

	id := request.PathValue("id")
	controller.Update(response, id, body)
}

func (server *Server) handleReferentialPartnerSave(response http.ResponseWriter, request *http.Request) {
	referentialSlug := request.PathValue("referential_slug")
	foundReferential := server.CurrentReferentials().FindBySlug(core.ReferentialSlug(referentialSlug))
	if foundReferential == nil {
		http.Error(response, "Referential not found", http.StatusNotFound)
		return
	}
	if !server.isAuth(foundReferential, request) {
		http.Error(response, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	partnerControler := NewPartnerController(foundReferential)
	response.Header().Set("Server", version.ApplicationName())
	response.Header().Set("Content-Type", "application/json")

	logger.Log.Debugf("Partner controller Save request: %v", request)

	partnerControler.(Savable).Save(response)
}

func (server *Server) handleReferentialModelDelete(response http.ResponseWriter, request *http.Request) {
	controller, model := server.referentialModelSetup(response, request)
	if controller == nil {
		return
	}

	logger.Log.Debugf("%s controller Delete request: %v", model, request)

	id := request.PathValue("id")
	controller.Delete(response, id)
}

func (server *Server) handleReferentialImport(response http.ResponseWriter, request *http.Request) {
	referentialSlug := request.PathValue("referential_slug")
	foundReferential := server.CurrentReferentials().FindBySlug(core.ReferentialSlug(referentialSlug))
	if foundReferential == nil {
		http.Error(response, "Referential not found", http.StatusNotFound)
		return
	}

	if !server.isAuthForImport(foundReferential, request) {
		http.Error(response, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	response.Header().Set("Server", version.ApplicationName())
	response.Header().Set("Content-Type", "application/json")

	logger.Log.Debugf("Import controller request: %v", request)

	controller := NewImportController(foundReferential)
	controller.serve(response, request)
}

func (server *Server) handleReferentialIndex(response http.ResponseWriter, request *http.Request) {
	controller := server.referentialSetup(response, request)

	controller.Index(response, nil)
}

func (server *Server) handleReferentialCreate(response http.ResponseWriter, request *http.Request) {
	controller := server.referentialSetup(response, request)

	controller.Create(response, request)
}

func (server *Server) handleReferentialGet(response http.ResponseWriter, request *http.Request) {
	data := request.PathValue("model")
	referentialSlug := request.PathValue("referential_slug")
	if data == "gtfs" && referentialSlug != "_referentials" {
		server.handleGtfs(response, request)
	}

	if referentialSlug == "_referentials" {
		controller := server.referentialSetup(response, request)
		controller.Show(response, data)
		return
	}

	server.handleReferentialModelIndex(response, request)
}

func (server *Server) handleReferentialUpdate(response http.ResponseWriter, request *http.Request) {
	controller := server.referentialSetup(response, request)

	id := request.PathValue("id")
	controller.Update(response, request, id)
}

func (server *Server) handleReferentialDelete(response http.ResponseWriter, request *http.Request) {
	controller := server.referentialSetup(response, request)

	id := request.PathValue("id")
	controller.Delete(response, id)
}

func (server *Server) handleReferentialSave(response http.ResponseWriter, request *http.Request) {
	controller := server.referentialSetup(response, request)

	controller.Save(response)
}

func (server *Server) handleReferentialReload(response http.ResponseWriter, request *http.Request) {
	controller := server.referentialSetup(response, request)

	id := request.PathValue("id")
	controller.reload(id, response)
}

func (server *Server) handleStatus(response http.ResponseWriter, request *http.Request) {
	defer monitoring.HandleHttpPanic(response)

	response.Header().Set("Server", version.ApplicationName())
	response.Header().Set("Content-Type", "application/json")

	controller := NewStatusController(server)
	controller.serve(response, request)
}

func (server *Server) handleTimeGet(response http.ResponseWriter, request *http.Request) {
	defer monitoring.HandleHttpPanic(response)

	response.Header().Set("Server", version.ApplicationName())
	response.Header().Set("Content-Type", "application/json")

	controller := NewTimeController(server)
	controller.get(response)
}

func (server *Server) handleTimeAdvance(response http.ResponseWriter, request *http.Request) {
	defer monitoring.HandleHttpPanic(response)

	response.Header().Set("Server", version.ApplicationName())
	response.Header().Set("Content-Type", "application/json")

	controller := NewTimeController(server)
	controller.advance(response, request)
}

func (server *Server) handleSIRILite(response http.ResponseWriter, request *http.Request) {
	defer monitoring.HandleHttpPanic(response)

	referentialSlug := request.PathValue("referential_slug")
	foundReferential := server.CurrentReferentials().FindBySlug(core.ReferentialSlug(referentialSlug))

	logger.Log.Debugf("SIRI Lite request: %v", request)

	siriLiteHandler := NewSIRILiteHandler(foundReferential, server.getToken(request))
	siriLiteHandler.serve(response, request)
}

func (server *Server) HandleSIRI(response http.ResponseWriter, request *http.Request) {
	defer monitoring.HandleHttpPanic(response)

	referentialSlug := request.PathValue("referential_slug")
	foundReferential := server.CurrentReferentials().FindBySlug(core.ReferentialSlug(referentialSlug))

	logger.Log.Debugf("SIRI request: %v", request)

	response.Header().Set("Server", version.ApplicationName())

	siriHandler := NewSIRIHandler(foundReferential)
	siriHandler.serve(response, request)
}

func (server *Server) handlePush(response http.ResponseWriter, request *http.Request) {
	defer monitoring.HandleHttpPanic(response)

	referentialSlug := request.PathValue("referential_slug")
	foundReferential := server.CurrentReferentials().FindBySlug(core.ReferentialSlug(referentialSlug))
	if foundReferential == nil {
		http.Error(response, "Referential not found", http.StatusNotFound)
		return
	}

	logger.Log.Debugf("Push request: %v", request)

	response.Header().Set("Server", version.ApplicationName())

	pushHandler := NewPushHandler(foundReferential, server.getToken(request))
	pushHandler.serve(response, request)
}

func (server *Server) handleGtfs(response http.ResponseWriter, request *http.Request) {
	defer monitoring.HandleHttpPanic(response)

	referentialSlug := request.PathValue("referential_slug")
	foundReferential := server.CurrentReferentials().FindBySlug(core.ReferentialSlug(referentialSlug))
	if foundReferential == nil {
		http.Error(response, "Referential not found", http.StatusNotFound)
		return
	}

	logger.Log.Debugf("Gtfs request: %v", request)

	response.Header().Set("Server", version.ApplicationName())

	gtfsHandler := NewGtfsHandler(foundReferential, server.getToken(request))

	resource := request.PathValue("resource")
	gtfsHandler.serve(response, request, resource)
}

func (server *Server) handleGraphql(response http.ResponseWriter, request *http.Request) {
	defer monitoring.HandleHttpPanic(response)

	referentialSlug := request.PathValue("referential_slug")
	foundReferential := server.CurrentReferentials().FindBySlug(core.ReferentialSlug(referentialSlug))
	if foundReferential == nil {
		http.Error(response, "Referential not found", http.StatusNotFound)
		return
	}

	logger.Log.Debugf("Graphql request: %v", request)

	response.Header().Set("Server", version.ApplicationName())

	graphqlHandler := NewGraphqlHandler(foundReferential, server.getToken(request))
	graphqlHandler.serve(response, request)
}

func (server *Server) partnerSubscriptionsSetup(response http.ResponseWriter, request *http.Request) (controller RestfulResource) {
	defer monitoring.HandleHttpPanic(response)

	referentialSlug := request.PathValue("referential_slug")
	foundReferential := server.CurrentReferentials().FindBySlug(core.ReferentialSlug(referentialSlug))
	if foundReferential == nil {
		http.Error(response, "Referential not found", http.StatusNotFound)
		return
	}

	if !server.isAuth(foundReferential, request) {
		http.Error(response, "Unauthorized request", http.StatusUnauthorized)
		return
	}

	controller = NewPartnerController(foundReferential)
	return
}

func (server *Server) handlePartnerSubscriptionsIndex(response http.ResponseWriter, request *http.Request) {
	controller := server.partnerSubscriptionsSetup(response, request)

	logger.Log.Debugf("Partner Subscriptions controller Index request: %v", request)

	id := request.PathValue("id")
	controller.(SubscriptionResource).SubscriptionsIndex(response, id)
}

func (server *Server) handlePartnerSubscriptionsCreate(response http.ResponseWriter, request *http.Request) {
	controller := server.partnerSubscriptionsSetup(response, request)

	body := getRequestBody(response, request)
	if body == nil {
		return
	}

	logger.Log.Debugf("Partner Subscriptions controller Create request: %v", request)

	id := request.PathValue("id")
	controller.(SubscriptionResource).SubscriptionsCreate(response, id, body)
}

func (server *Server) handlePartnerSubscriptionsDelete(response http.ResponseWriter, request *http.Request) {
	controller := server.partnerSubscriptionsSetup(response, request)

	logger.Log.Debugf("Partner Subscriptions controller Delete request: %v", request)

	subscriptionId := request.PathValue("subscription_id")
	id := request.PathValue("id")
	controller.(SubscriptionResource).SubscriptionsDelete(response, id, subscriptionId)
}
