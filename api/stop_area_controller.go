package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type StopAreaController struct {
	ControllerReferential
}

func NewStopAreaController() (controller *Controller) {
	return &Controller{
		ressourceController: &StopAreaController{},
	}
}

func (controller *StopAreaController) Ressources() string {
	return "stop_areas"
}

func (controller *StopAreaController) Index(response http.ResponseWriter) {
	logger.Log.Debugf("StopAreas Index")

	jsonBytes, _ := json.Marshal(controller.referential.Model().StopAreas().FindAll())
	response.Write(jsonBytes)
}

func (controller *StopAreaController) Show(response http.ResponseWriter, identifier string) {
	stopArea, ok := controller.referential.Model().StopAreas().Find(model.StopAreaId(identifier))
	if !ok {
		http.Error(response, fmt.Sprintf("Stop area not found: %s", identifier), 500)
		return
	}
	logger.Log.Debugf("Get stopArea %s", identifier)

	jsonBytes, _ := stopArea.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *StopAreaController) Delete(response http.ResponseWriter, identifier string) {
	stopArea, ok := controller.referential.Model().StopAreas().Find(model.StopAreaId(identifier))
	if !ok {
		http.Error(response, fmt.Sprintf("Stop area not found: %s", identifier), 500)
		return
	}
	logger.Log.Debugf("Delete stopArea %s", identifier)

	jsonBytes, _ := stopArea.MarshalJSON()
	controller.referential.Model().StopAreas().Delete(&stopArea)
	response.Write(jsonBytes)
}

func (controller *StopAreaController) Update(response http.ResponseWriter, identifier string, body []byte) {
	stopArea, ok := controller.referential.Model().StopAreas().Find(model.StopAreaId(identifier))
	if !ok {
		http.Error(response, fmt.Sprintf("Stop area not found: %s", identifier), 500)
		return
	}

	logger.Log.Debugf("Update stopArea %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &stopArea)
	if err != nil {
		http.Error(response, "Invalid request: can't parse request body", 400)
		return
	}

	controller.referential.Model().StopAreas().Save(&stopArea)
	jsonBytes, _ := stopArea.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *StopAreaController) Create(response http.ResponseWriter, body []byte) {
	logger.Log.Debugf("Create stopArea: %s", string(body))

	stopArea := controller.referential.Model().StopAreas().New()

	err := json.Unmarshal(body, &stopArea)
	if err != nil {
		http.Error(response, "Invalid request: can't parse request body", 400)
		return
	}
	if stopArea.Id() != "" {
		http.Error(response, "Invalid request", 400)
		return
	}

	controller.referential.Model().StopAreas().Save(&stopArea)
	jsonBytes, _ := stopArea.MarshalJSON()
	response.Write(jsonBytes)
}
