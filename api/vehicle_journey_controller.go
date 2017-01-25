package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type VehicleJourneyController struct {
	referential *core.Referential
}

func NewVehicleJourneyController(referential *core.Referential) *Controller {
	return &Controller{
		restfulRessource: &VehicleJourneyController{
			referential: referential,
		},
	}
}

func (controller *VehicleJourneyController) Index(response http.ResponseWriter) {
	logger.Log.Debugf("VehicleJourneys Index")

	jsonBytes, _ := json.Marshal(controller.referential.Model().VehicleJourneys().FindAll())
	response.Write(jsonBytes)
}

func (controller *VehicleJourneyController) Show(response http.ResponseWriter, identifier string) {
	vehicleJourney, ok := controller.referential.Model().VehicleJourneys().Find(model.VehicleJourneyId(identifier))
	if !ok {
		http.Error(response, fmt.Sprintf("Vehicle journey not found: %s", identifier), 500)
		return
	}
	logger.Log.Debugf("Get vehicleJourney %s", identifier)

	jsonBytes, _ := vehicleJourney.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *VehicleJourneyController) Delete(response http.ResponseWriter, identifier string) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	vehicleJourney, ok := tx.Model().VehicleJourneys().Find(model.VehicleJourneyId(identifier))
	if !ok {
		http.Error(response, fmt.Sprintf("Vehicle journey not found: %s", identifier), 500)
		return
	}
	logger.Log.Debugf("Delete vehicleJourney %s", identifier)

	jsonBytes, _ := vehicleJourney.MarshalJSON()
	tx.Model().VehicleJourneys().Delete(&vehicleJourney)
	err := tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", 500)
		return
	}
	response.Write(jsonBytes)
}

func (controller *VehicleJourneyController) Update(response http.ResponseWriter, identifier string, body []byte) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	vehicleJourney, ok := tx.Model().VehicleJourneys().Find(model.VehicleJourneyId(identifier))
	if !ok {
		http.Error(response, fmt.Sprintf("Vehicle journey not found: %s", identifier), 500)
		return
	}

	logger.Log.Debugf("Update vehicleJourney %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &vehicleJourney)
	if err != nil {
		http.Error(response, "Invalid request: can't parse request body", 400)
		return
	}

	tx.Model().VehicleJourneys().Save(&vehicleJourney)
	err = tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", 500)
		return
	}
	jsonBytes, _ := vehicleJourney.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *VehicleJourneyController) Create(response http.ResponseWriter, body []byte) {
	// New transaction
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	logger.Log.Debugf("Create vehicleJourney: %s", string(body))

	vehicleJourney := tx.Model().VehicleJourneys().New()

	err := json.Unmarshal(body, &vehicleJourney)
	if err != nil {
		http.Error(response, "Invalid request: can't parse request body", 400)
		return
	}
	if vehicleJourney.Id() != "" {
		http.Error(response, "Invalid request", 400)
		return
	}

	tx.Model().VehicleJourneys().Save(&vehicleJourney)
	err = tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", 500)
		return
	}
	jsonBytes, _ := vehicleJourney.MarshalJSON()
	response.Write(jsonBytes)
}
