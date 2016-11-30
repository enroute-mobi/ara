package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type StopAreaController struct {
	referential *core.Referential
}

func NewStopAreaController(referential *core.Referential) *Controller {
	return &Controller{
		restfulRessource: &StopAreaController{
			referential: referential,
		},
	}
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
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	stopArea, ok := tx.Model().StopAreas().Find(model.StopAreaId(identifier))
	if !ok {
		http.Error(response, fmt.Sprintf("Stop area not found: %s", identifier), 500)
		return
	}
	logger.Log.Debugf("Delete stopArea %s", identifier)

	jsonBytes, _ := stopArea.MarshalJSON()
	tx.Model().StopAreas().Delete(&stopArea)
	err := tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", 500)
		return
	}
	response.Write(jsonBytes)
}

func (controller *StopAreaController) Update(response http.ResponseWriter, identifier string, body []byte) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	stopArea, ok := tx.Model().StopAreas().Find(model.StopAreaId(identifier))
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

	tx.Model().StopAreas().Save(&stopArea)
	err = tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", 500)
		return
	}
	jsonBytes, _ := stopArea.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *StopAreaController) Create(response http.ResponseWriter, body []byte) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	logger.Log.Debugf("Create stopArea: %s", string(body))

	stopArea := tx.Model().StopAreas().New()

	err := json.Unmarshal(body, &stopArea)
	if err != nil {
		http.Error(response, "Invalid request: can't parse request body", 400)
		return
	}
	if stopArea.Id() != "" {
		http.Error(response, "Invalid request", 400)
		return
	}

	tx.Model().StopAreas().Save(&stopArea)
	err = tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", 500)
		return
	}
	jsonBytes, _ := stopArea.MarshalJSON()
	response.Write(jsonBytes)
}
