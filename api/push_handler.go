package api

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	external_models "bitbucket.org/enroute-mobi/ara-external-models"
	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"github.com/golang/protobuf/proto"
)

type PushHandler struct {
	referential *core.Referential
	token       string
}

func NewPushHandler(referential *core.Referential, token string) *PushHandler {
	return &PushHandler{
		referential: referential,
		token:       token,
	}
}

func (handler *PushHandler) serve(response http.ResponseWriter, request *http.Request) {
	// Check if request header is protobuf or return an error
	if request.Header.Get("Content-Type") != "application/x-protobuf" {
		http.Error(response, "Expected application/x-protobuf content", http.StatusUnsupportedMediaType)
		return
	}

	// Find Partner by authorization Key
	if handler.token == "" {
		http.Error(response, "Invalid Authorization Token", http.StatusUnauthorized)
		return
	}
	partner, ok := handler.referential.Partners().FindByCredential(handler.token)
	if !ok {
		http.Error(response, "Invalid Authorization Token", http.StatusUnauthorized)
		return
	}

	// Find Push connector
	connector, ok := partner.Connector(core.PUSH_COLLECTOR)
	if !ok {
		http.Error(response, "Partner doesn't have a push connector", http.StatusNotImplemented)
		return
	}

	startTime := handler.referential.Clock().Now()
	message := handler.newBQMessage(string(partner.Slug()), request.RemoteAddr)
	defer audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(message)

	// Check if request is gzip
	var requestReader io.Reader
	if request.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(request.Body)
		if err != nil {
			handler.logError(message, startTime, "Can't unzip request")
			http.Error(response, "Can't unzip request", http.StatusBadRequest)
			return
		}
		defer gzipReader.Close()
		requestReader = gzipReader
	} else {
		requestReader = request.Body
	}

	// Attempt to read the body
	content, err := ioutil.ReadAll(requestReader)
	if err != nil {
		e := fmt.Sprintf("Error while reading body: %v", err)
		handler.logError(message, startTime, e)
		http.Error(response, e, http.StatusBadRequest)
		return
	}
	message.RequestSize = int64(len(content))
	if len(content) == 0 {
		handler.logError(message, startTime, "Empty body")
		http.Error(response, "Empty body", http.StatusBadRequest)
		return
	}

	externalModel := &external_models.ExternalCompleteModel{}
	err = proto.Unmarshal(content, externalModel)
	if err != nil {
		e := fmt.Sprintf("Error while unmarshalling body: %v", err)
		handler.logError(message, startTime, e)
		http.Error(response, e, http.StatusBadRequest)
		return
	}

	connector.(*core.PushCollector).HandlePushNotification(externalModel, message)
	response.WriteHeader(http.StatusOK)
}

func (handler *PushHandler) newBQMessage(slug, remoteAddress string) *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Protocol:  "push",
		Type:      "push-notification",
		Direction: "received",
		Partner:   slug,
		Status:    "OK",
		IPAddress: remoteAddress,
	}
}

func (handler *PushHandler) logError(m *audit.BigQueryMessage, startTime time.Time, format string, values ...interface{}) {
	m.ProcessingTime = handler.referential.Clock().Since(startTime).Seconds()
	m.Status = "Error"
	errorString := fmt.Sprintf(format, values...)

	m.ErrorDetails = errorString
	logger.Log.Debugf(errorString)
}
