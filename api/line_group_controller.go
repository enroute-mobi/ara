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

type LineGroupsController struct {
	referential *core.Referential
}

func NewLineGroupsController(referential *core.Referential) RestfulResource {
	return &LineGroupsController{
		referential: referential,
	}
}

func (controller *LineGroupsController) findLineGroup(identifier string) (*model.LineGroup, bool) {
	return controller.referential.Model().LineGroups().Find(model.LineGroupId(identifier))
}

func (controller *LineGroupsController) Index(response http.ResponseWriter, _params url.Values) {
	logger.Log.Debugf("LineGroup Index")
	controller.referential.Model().Lines()
	jsonBytes, _ := json.Marshal(controller.referential.Model().LineGroups().FindAll())
	response.Write(jsonBytes)
}

func (controller *LineGroupsController) Show(response http.ResponseWriter, identifier string) {
	lineGroup, ok := controller.findLineGroup(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("LineGroup not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Get lineGroup %s", identifier)

	jsonBytes, _ := lineGroup.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *LineGroupsController) Delete(response http.ResponseWriter, identifier string) {
	lineGroup, ok := controller.findLineGroup(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("LineGroup not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Delete lineGroup %s", identifier)

	jsonBytes, _ := lineGroup.MarshalJSON()
	controller.referential.Model().LineGroups().Delete(lineGroup)

	response.Write(jsonBytes)
}

func (controller *LineGroupsController) Update(response http.ResponseWriter, identifier string, body []byte) {
	lineGroup, ok := controller.findLineGroup(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("LineGroup not found: %s", identifier), http.StatusNotFound)
		return
	}

	logger.Log.Debugf("Update lineGroup %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &lineGroup)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	controller.referential.Model().LineGroups().Save(lineGroup)
	jsonBytes, _ := json.Marshal(&lineGroup)
	response.Write(jsonBytes)
}

func (controller *LineGroupsController) Create(response http.ResponseWriter, body []byte) {
	logger.Log.Debugf("Create lineGroup: %s", string(body))

	lineGroup := controller.referential.Model().LineGroups().New()

	err := json.Unmarshal(body, &lineGroup)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	if lineGroup.Id() != "" {
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	controller.referential.Model().LineGroups().Save(lineGroup)
	jsonBytes, _ := json.Marshal(&lineGroup)
	response.Write(jsonBytes)
}
