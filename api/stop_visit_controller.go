package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type StopVisitController struct {
	referential *core.Referential
}

func NewStopVisitController(referential *core.Referential) ControllerInterface {
	return &Controller{
		restfulRessource: &StopVisitController{
			referential: referential,
		},
	}
}

func (controller *StopVisitController) Index(response http.ResponseWriter) {
	logger.Log.Debugf("StopVisits Index")

	jsonBytes, _ := json.Marshal(controller.referential.Model().StopVisits().FindAll())
	response.Write(jsonBytes)
}

func (controller *StopVisitController) Show(response http.ResponseWriter, identifier string) {
	stopVisit, ok := controller.referential.Model().StopVisits().Find(model.StopVisitId(identifier))
	if !ok {
		http.Error(response, fmt.Sprintf("Stop visit not found: %s", identifier), 500)
		return
	}
	logger.Log.Debugf("Get stopVisit %s", identifier)

	jsonBytes, _ := stopVisit.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *StopVisitController) Delete(response http.ResponseWriter, identifier string) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	stopVisit, ok := tx.Model().StopVisits().Find(model.StopVisitId(identifier))
	if !ok {
		http.Error(response, fmt.Sprintf("Stop visit not found: %s", identifier), 500)
		return
	}
	logger.Log.Debugf("Delete stopVisit %s", identifier)

	jsonBytes, _ := stopVisit.MarshalJSON()
	tx.Model().StopVisits().Delete(&stopVisit)
	err := tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", 500)
		return
	}
	response.Write(jsonBytes)
}

func (controller *StopVisitController) Update(response http.ResponseWriter, identifier string, body []byte) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	stopVisit, ok := tx.Model().StopVisits().Find(model.StopVisitId(identifier))
	if !ok {
		http.Error(response, fmt.Sprintf("Stop visit not found: %s", identifier), 500)
		return
	}

	logger.Log.Debugf("Update stopVisit %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &stopVisit)
	if err != nil {
		http.Error(response, "Invalid request: can't parse request body", 400)
		return
	}

	tx.Model().StopVisits().Save(&stopVisit)
	err = tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", 500)
		return
	}
	jsonBytes, _ := stopVisit.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *StopVisitController) Create(response http.ResponseWriter, body []byte) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	logger.Log.Debugf("Create stopVisit: %s", string(body))

	stopVisit := tx.Model().StopVisits().New()

	err := json.Unmarshal(body, &stopVisit)
	if err != nil {
		http.Error(response, "Invalid request: can't parse request body", 400)
		return
	}
	if stopVisit.Id() != "" {
		http.Error(response, "Invalid request", 400)
		return
	}

	tx.Model().StopVisits().Save(&stopVisit)
	err = tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", 500)
		return
	}
	jsonBytes, _ := stopVisit.MarshalJSON()
	response.Write(jsonBytes)
}
