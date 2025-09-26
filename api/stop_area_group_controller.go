package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
)

type StopAreaGroupsController struct {
	referential *core.Referential
}

func NewStopAreaGroupsController(referential *core.Referential) RestfulResource {
	return &StopAreaGroupsController{
		referential: referential,
	}
}

func (controller *StopAreaGroupsController) findStopAreaGroup(identifier string) (*model.StopAreaGroup, bool) {
	return controller.referential.Model().StopAreaGroups().Find(model.StopAreaGroupId(identifier))
}

func (controller *StopAreaGroupsController) Index(response http.ResponseWriter, _params url.Values) {
	logger.Log.Debugf("StopAreaGroup Index")
	controller.referential.Model().Lines()
	jsonBytes, _ := json.Marshal(controller.referential.Model().StopAreaGroups().FindAll())
	response.Write(jsonBytes)
}

func (controller *StopAreaGroupsController) Show(response http.ResponseWriter, identifier string, _params url.Values) {
	stopAreaGroup, ok := controller.findStopAreaGroup(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("StopAreaGroup not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Get stopAreaGroup %s", identifier)

	jsonBytes, _ := stopAreaGroup.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *StopAreaGroupsController) Delete(response http.ResponseWriter, identifier string) {
	stopAreaGroup, ok := controller.findStopAreaGroup(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("StopAreaGroup not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Delete stopAreaGroup %s", identifier)

	jsonBytes, _ := stopAreaGroup.MarshalJSON()
	controller.referential.Model().StopAreaGroups().Delete(stopAreaGroup)

	response.Write(jsonBytes)
}

func (controller *StopAreaGroupsController) Update(response http.ResponseWriter, identifier string, body []byte) {
	stopAreaGroup, ok := controller.findStopAreaGroup(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("StopAreaGroup not found: %s", identifier), http.StatusNotFound)
		return
	}

	logger.Log.Debugf("Update stopAreaGroup %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &stopAreaGroup)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	controller.referential.Model().StopAreaGroups().Save(stopAreaGroup)
	jsonBytes, _ := json.Marshal(&stopAreaGroup)
	response.Write(jsonBytes)
}

func (controller *StopAreaGroupsController) Create(response http.ResponseWriter, body []byte) {
	logger.Log.Debugf("Create stopAreaGroup: %s", string(body))

	stopAreaGroup := controller.referential.Model().StopAreaGroups().New()

	err := json.Unmarshal(body, &stopAreaGroup)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	if stopAreaGroup.Id() != "" {
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	controller.referential.Model().StopAreaGroups().Save(stopAreaGroup)
	jsonBytes, _ := json.Marshal(&stopAreaGroup)
	response.Write(jsonBytes)
}
