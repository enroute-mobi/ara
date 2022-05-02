package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
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

func (controller *VehicleJourneyController) findVehicleJourney(identifier string) (*model.VehicleJourney, bool) {
	idRegexp := "([0-9a-zA-Z-]+):([0-9a-zA-Z-:]+)"
	pattern := regexp.MustCompile(idRegexp)
	foundStrings := pattern.FindStringSubmatch(identifier)
	if foundStrings != nil {
		objectid := model.NewObjectID(foundStrings[1], foundStrings[2])
		return controller.referential.Model().VehicleJourneys().FindByObjectId(objectid)
	}
	return controller.referential.Model().VehicleJourneys().Find(model.VehicleJourneyId(identifier))
}

func (controller *VehicleJourneyController) Index(response http.ResponseWriter, filters url.Values) {
	logger.Log.Debugf("VehicleJourneys Index")

	jsonBytes, _ := json.Marshal(controller.referential.Model().VehicleJourneys().FindAll())
	response.Write(jsonBytes)
}

func (controller *VehicleJourneyController) Show(response http.ResponseWriter, identifier string) {
	vehicleJourney, ok := controller.findVehicleJourney(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Vehicle journey not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Get vehicleJourney %s", identifier)

	jsonBytes, _ := vehicleJourney.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *VehicleJourneyController) Delete(response http.ResponseWriter, identifier string) {
	vehicleJourney, ok := controller.findVehicleJourney(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Vehicle journey not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Delete vehicleJourney %s", identifier)

	jsonBytes, _ := vehicleJourney.MarshalJSON()
	controller.referential.Model().VehicleJourneys().Delete(vehicleJourney)
	response.Write(jsonBytes)
}

func (controller *VehicleJourneyController) Update(response http.ResponseWriter, identifier string, body []byte) {
	vehicleJourney, ok := controller.findVehicleJourney(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Vehicle journey not found: %s", identifier), http.StatusNotFound)
		return
	}

	logger.Log.Debugf("Update vehicleJourney %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &vehicleJourney)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	for _, obj := range vehicleJourney.ObjectIDs() {
		vj, ok := controller.referential.Model().VehicleJourneys().FindByObjectId(obj)
		if ok && vj.Id() != vehicleJourney.Id() {
			http.Error(response, fmt.Sprintf("Invalid request: vehicleJourney %v already have an objectid %v", vj.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	controller.referential.Model().VehicleJourneys().Save(vehicleJourney)
	jsonBytes, _ := vehicleJourney.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *VehicleJourneyController) Create(response http.ResponseWriter, body []byte) {
	logger.Log.Debugf("Create vehicleJourney: %s", string(body))

	vehicleJourney := controller.referential.Model().VehicleJourneys().New()

	err := json.Unmarshal(body, &vehicleJourney)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	if vehicleJourney.Id() != "" {
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	for _, obj := range vehicleJourney.ObjectIDs() {
		vj, ok := controller.referential.Model().VehicleJourneys().FindByObjectId(obj)
		if ok {
			http.Error(response, fmt.Sprintf("Invalid request: vehicleJourney %v already have an objectid %v", vj.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	controller.referential.Model().VehicleJourneys().Save(vehicleJourney)
	jsonBytes, _ := vehicleJourney.MarshalJSON()
	response.Write(jsonBytes)
}
