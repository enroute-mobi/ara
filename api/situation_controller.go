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

type SituationController struct {
	referential *core.Referential
}

func NewSituationController(referential *core.Referential) ControllerInterface {
	return &Controller{
		restfulResource: &SituationController{
			referential: referential,
		},
	}
}

func (controller *SituationController) findSituation(identifier string) (model.Situation, bool) {
	foundStrings := idPattern.FindStringSubmatch(identifier)
	if foundStrings != nil {
		objectid := model.NewObjectID(foundStrings[1], foundStrings[2])
		return controller.referential.Model().Situations().FindByObjectId(objectid)
	}
	return controller.referential.Model().Situations().Find(model.SituationId(identifier))
}

func (controller *SituationController) Index(response http.ResponseWriter, filters url.Values) {
	logger.Log.Debugf("Situations Index")

	jsonBytes, _ := json.Marshal(controller.referential.Model().Situations().FindAll())
	response.Write(jsonBytes)
}

func (controller *SituationController) Show(response http.ResponseWriter, identifier string) {
	situation, ok := controller.findSituation(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Situation not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Get situation %s", identifier)

	jsonBytes, _ := situation.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *SituationController) Delete(response http.ResponseWriter, identifier string) {
	situation, ok := controller.findSituation(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Situation not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Delete situation %s", identifier)

	jsonBytes, _ := situation.MarshalJSON()
	controller.referential.Model().Situations().Delete(&situation)
	response.Write(jsonBytes)
}

func (controller *SituationController) Update(response http.ResponseWriter, identifier string, body []byte) {
	situation, ok := controller.findSituation(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Situation not found: %s", identifier), http.StatusNotFound)
		return
	}

	logger.Log.Debugf("Update situation %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &situation)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	controller.referential.Model().Situations().Save(&situation)
	jsonBytes, _ := json.Marshal(&situation)
	response.Write(jsonBytes)
}

func (controller *SituationController) Create(response http.ResponseWriter, body []byte) {
	logger.Log.Debugf("Create situation: %s", string(body))

	situation := controller.referential.Model().Situations().New()

	err := json.Unmarshal(body, &situation)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	if situation.Id() != "" {
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	controller.referential.Model().Situations().Save(&situation)
	jsonBytes, _ := json.Marshal(&situation)
	response.Write(jsonBytes)
}
