package api

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	external_models "github.com/af83/ara-external-models"
	"github.com/af83/edwig/core"
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
	if handler.referential == nil {
		http.Error(response, "Referential not found", http.StatusNotFound)
		return
	}

	// Check if request header is protobuf or return an error
	if request.Header.Get("Content-Type") != "application/x-protobuf" {
		http.Error(response, "Expected application/x-protobuf content", http.StatusUnsupportedMediaType)
		return
	}

	// Find Partner by authorization Key
	partner, ok := handler.referential.Partners().FindBySetting(core.LOCAL_CREDENTIAL, handler.token)
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

	// Check if request is gzip
	var requestReader io.Reader
	if request.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(request.Body)
		if err != nil {
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
		http.Error(response, e, http.StatusBadRequest)
		return
	}
	if len(content) == 0 {
		http.Error(response, "Empty body", http.StatusBadRequest)
		return
	}

	externalModel := &external_models.ExternalCompleteModel{}
	err = proto.Unmarshal(content, externalModel)
	if err != nil {
		e := fmt.Sprintf("Error while unmarshalling body: %v", err)
		http.Error(response, e, http.StatusBadRequest)
		return
	}

	connector.(*core.PushCollector).HandlePushNotification(externalModel)
	response.WriteHeader(http.StatusOK)
}
