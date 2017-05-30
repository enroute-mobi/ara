package api

import (
	"encoding/json"
	"net/http"

	"github.com/af83/edwig/version"
)

type StatusController struct{}
type Status struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func NewStatusController(server *Server) ControllerInterface {
	return &StatusController{}
}

func (controller *StatusController) serve(response http.ResponseWriter, request *http.Request, requestData *RequestData) {
	status := Status{
		Status:  "ok",
		Version: version.Value(),
	}

	jsonBytes, _ := json.Marshal(status)
	response.Write(jsonBytes)
}
