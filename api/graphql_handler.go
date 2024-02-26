package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/enroute-mobi/ara/api/gql"
	"bitbucket.org/enroute-mobi/ara/api/rah"
	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"

	graphql "github.com/graph-gophers/graphql-go"
)

type GraphqlHandler struct {
	rah.RemoteAddressHandler

	referential *core.Referential
	token       string
}

func NewGraphqlHandler(referential *core.Referential, token string) *GraphqlHandler {
	return &GraphqlHandler{
		referential: referential,
		token:       token,
	}
}

func (handler *GraphqlHandler) serve(response http.ResponseWriter, request *http.Request) {
	// Find Partner by authorization Key
	partner, ok := handler.referential.Partners().FindByCredential(handler.token)
	if !ok {
		http.Error(response, "Invalid Authorization Token", http.StatusUnauthorized)
		return
	}

	// Find Push connector
	_, ok = partner.Connector(core.GRAPHQL_SERVER)
	if !ok {
		http.Error(response, "Partner doesn't have a graphql connector", http.StatusNotImplemented)
		return
	}

	startTime := handler.referential.Clock().Now()
	message := handler.newBQMessage(string(partner.Slug()), handler.HandleRemoteAddress(request))
	defer audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(message)

	var params struct {
		Query         string                 `json:"query"`
		OperationName string                 `json:"operationName"`
		Variables     map[string]interface{} `json:"variables"`
	}
	if err := json.NewDecoder(request.Body).Decode(&params); err != nil {
		e := fmt.Sprintf("Error while reading body: %v", err)
		handler.logError(message, startTime, e)
		http.Error(response, e, http.StatusBadRequest)
		return
	}

	schema := graphql.MustParseSchema(gql.Schema, &gql.Resolver{Partner: partner})
	r := schema.Exec(request.Context(), params.Query, params.OperationName, params.Variables)
	responseJSON, err := json.Marshal(r)
	if err != nil {
		e := fmt.Sprintf("Error while reading body: %v", err)
		handler.logError(message, startTime, e)
		http.Error(response, e, http.StatusBadRequest)
		return
	}

	processingTime := handler.referential.Clock().Since(startTime)

	message.RequestRawMessage = params.Query
	message.ProcessingTime = processingTime.Seconds()

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(http.StatusOK)
	response.Write(responseJSON)
}

func (handler *GraphqlHandler) newBQMessage(slug, remoteAddress string) *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Protocol:  "graphql",
		Direction: "received",
		Status:    "OK",
		Partner:   slug,
		IPAddress: remoteAddress,
	}
}

func (handler *GraphqlHandler) logError(m *audit.BigQueryMessage, startTime time.Time, format string, values ...interface{}) {
	m.ProcessingTime = handler.referential.Clock().Since(startTime).Seconds()
	m.Status = "Error"
	errorString := fmt.Sprintf(format, values...)

	m.ErrorDetails = errorString
	logger.Log.Debugf(errorString)
}
