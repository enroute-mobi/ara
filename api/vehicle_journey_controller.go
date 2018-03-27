package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type VehicleJourneyController struct {
	referential *core.Referential
}

func NewVehicleJourneyController(referential *core.Referential) ControllerInterface {
	return &Controller{
		restfulResource: &VehicleJourneyController{
			referential: referential,
		},
	}
}

func (controller *VehicleJourneyController) findVehicleJourney(tx *model.Transaction, identifier string) (model.VehicleJourney, bool) {
	idRegexp := "([0-9a-zA-Z-]+):([0-9a-zA-Z-:]+)"
	pattern := regexp.MustCompile(idRegexp)
	foundStrings := pattern.FindStringSubmatch(identifier)
	if foundStrings != nil {
		objectid := model.NewObjectID(foundStrings[1], foundStrings[2])
		return tx.Model().VehicleJourneys().FindByObjectId(objectid)
	}
	return tx.Model().VehicleJourneys().Find(model.VehicleJourneyId(identifier))
}

func (controller *VehicleJourneyController) Index(response http.ResponseWriter, filters url.Values) {
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	logger.Log.Debugf("VehicleJourneys Index")

	jsonBytes, _ := json.Marshal(tx.Model().VehicleJourneys().FindAll())
	response.Write(jsonBytes)
}

func (controller *VehicleJourneyController) Show(response http.ResponseWriter, identifier string) {
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	vehicleJourney, ok := controller.findVehicleJourney(tx, identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Vehicle journey not found: %s", identifier), 404)
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

	vehicleJourney, ok := controller.findVehicleJourney(tx, identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Vehicle journey not found: %s", identifier), 404)
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

	vehicleJourney, ok := controller.findVehicleJourney(tx, identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Vehicle journey not found: %s", identifier), 404)
		return
	}

	logger.Log.Debugf("Update vehicleJourney %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &vehicleJourney)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), 400)
		return
	}

	for _, obj := range vehicleJourney.ObjectIDs() {
		vj, ok := tx.Model().VehicleJourneys().FindByObjectId(obj)
		if ok && vj.Id() != vehicleJourney.Id() {
			http.Error(response, fmt.Sprintf("Invalid request: vehicleJourney %v already have an objectid %v", vj.Id(), obj.String()), 400)
			return
		}
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
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), 400)
		return
	}

	if vehicleJourney.Id() != "" {
		http.Error(response, "Invalid request", 400)
		return
	}

	for _, obj := range vehicleJourney.ObjectIDs() {
		vj, ok := tx.Model().VehicleJourneys().FindByObjectId(obj)
		if ok {
			http.Error(response, fmt.Sprintf("Invalid request: vehicleJourney %v already have an objectid %v", vj.Id(), obj.String()), 400)
			return
		}
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
