package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
)

type ReferentialController struct {
	server *Server
}

func NewReferentialController(server *Server) ControllerInterface {
	return &Controller{
		restfulRessource: &ReferentialController{
			server: server,
		},
	}
}

func (controller *ReferentialController) Index(response http.ResponseWriter) {
	logger.Log.Debugf("Referentials Index")

	jsonBytes, _ := json.Marshal(controller.server.CurrentReferentials().FindAll())
	response.Write(jsonBytes)
}

func (controller *ReferentialController) Show(response http.ResponseWriter, identifier string) {
	referential := controller.server.CurrentReferentials().Find(core.ReferentialId(identifier))
	if referential == nil {
		http.Error(response, fmt.Sprintf("Referential not found: %s", identifier), 500)
		return
	}
	logger.Log.Debugf("Get referential %s", identifier)

	jsonBytes, _ := referential.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *ReferentialController) Delete(response http.ResponseWriter, identifier string) {
	referential := controller.server.CurrentReferentials().Find(core.ReferentialId(identifier))
	if referential == nil {
		http.Error(response, fmt.Sprintf("Referential not found: %s", identifier), 500)
		return
	}
	logger.Log.Debugf("Delete referential %s", identifier)

	jsonBytes, _ := referential.MarshalJSON()
	referential.Stop()
	controller.server.CurrentReferentials().Delete(referential)
	response.Write(jsonBytes)
}

func (controller *ReferentialController) Update(response http.ResponseWriter, identifier string, body []byte) {
	referential := controller.server.CurrentReferentials().Find(core.ReferentialId(identifier))
	if referential == nil {
		http.Error(response, fmt.Sprintf("Referential not found: %s", identifier), 500)
		return
	}

	logger.Log.Debugf("Update referential %s: %s", identifier, string(body))

	referential.Stop()
	defer referential.Start()

	apiReferential := referential.Definition()
	err := json.Unmarshal(body, apiReferential)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), 400)
		return
	}

	if !apiReferential.Validate() {
		jsonBytes, _ := json.Marshal(apiReferential)
		response.WriteHeader(http.StatusBadRequest)
		response.Write(jsonBytes)
		return
	}

	referential.SetDefinition(apiReferential)
	referential.Save()
	jsonBytes, _ := referential.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *ReferentialController) Create(response http.ResponseWriter, body []byte) {
	logger.Log.Debugf("Create referential: %s", string(body))

	referential := controller.server.CurrentReferentials().New("")
	apiReferential := referential.Definition()
	err := json.Unmarshal(body, apiReferential)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), 400)
		return
	}

	if !apiReferential.Validate() {
		jsonBytes, _ := json.Marshal(apiReferential)
		response.WriteHeader(http.StatusBadRequest)
		response.Write(jsonBytes)
		return
	}

	referential.SetDefinition(apiReferential)
	referential.Save()
	referential.Start()
	jsonBytes, _ := referential.MarshalJSON()
	response.Write(jsonBytes)
}
