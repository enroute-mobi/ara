package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"time"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
)

type VehicleController struct {
	referential *core.Referential
}

func NewVehicleController(referential *core.Referential) ControllerInterface {
	return &Controller{
		restfulResource: &VehicleController{
			referential: referential,
		},
	}
}

func (controller *VehicleController) findVehicle(tx *model.Transaction, identifier string) (model.Vehicle, bool) {
	idRegexp := "([0-9a-zA-Z-]+):([0-9a-zA-Z-:]+)"
	pattern := regexp.MustCompile(idRegexp)
	foundStrings := pattern.FindStringSubmatch(identifier)
	if foundStrings != nil {
		objectid := model.NewObjectID(foundStrings[1], foundStrings[2])
		return tx.Model().Vehicles().FindByObjectId(objectid)
	}
	return tx.Model().Vehicles().Find(model.VehicleId(identifier))
}

func (controller *VehicleController) Index(response http.ResponseWriter, filters url.Values) {
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	logger.Log.Debugf("Vehicles Index")

	stime := controller.referential.Clock().Now()
	vehicles := tx.Model().Vehicles().FindAll()
	logger.Log.Debugf("VehicleController FindAll time : %v", time.Since(stime))
	stime = controller.referential.Clock().Now()
	jsonBytes, _ := json.Marshal(vehicles)
	logger.Log.Debugf("VehicleController Json Marshal time : %v ", time.Since(stime))
	response.Write(jsonBytes)
}

func (controller *VehicleController) Show(response http.ResponseWriter, identifier string) {
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	vehicle, ok := controller.findVehicle(tx, identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Vehicle not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Get vehicle %s", identifier)

	jsonBytes, _ := vehicle.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *VehicleController) Delete(response http.ResponseWriter, identifier string) {
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	vehicle, ok := controller.findVehicle(tx, identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Vehicle not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Delete vehicle %s", identifier)

	jsonBytes, _ := vehicle.MarshalJSON()
	tx.Model().Vehicles().Delete(&vehicle)
	err := tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", http.StatusInternalServerError)
		return
	}
	response.Write(jsonBytes)
}

func (controller *VehicleController) Update(response http.ResponseWriter, identifier string, body []byte) {
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	vehicle, ok := controller.findVehicle(tx, identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Vehicle not found: %s", identifier), http.StatusNotFound)
		return
	}

	logger.Log.Debugf("Update vehicle %s: %s", identifier, string(body))

	err := json.Unmarshal(body, &vehicle)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	for _, obj := range vehicle.ObjectIDs() {
		v, ok := tx.Model().Vehicles().FindByObjectId(obj)
		if ok && v.Id() != vehicle.Id() {
			http.Error(response, fmt.Sprintf("Invalid request: vehicle %v already have an objectid %v", v.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	tx.Model().Vehicles().Save(&vehicle)
	err = tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", http.StatusInternalServerError)
		return
	}
	jsonBytes, _ := vehicle.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *VehicleController) Create(response http.ResponseWriter, body []byte) {
	tx := controller.referential.NewTransaction()
	defer tx.Close()

	logger.Log.Debugf("Create vehicle: %s", string(body))

	vehicle := tx.Model().Vehicles().New()

	err := json.Unmarshal(body, &vehicle)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	if vehicle.Id() != "" {
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	for _, obj := range vehicle.ObjectIDs() {
		v, ok := tx.Model().Vehicles().FindByObjectId(obj)
		if ok {
			http.Error(response, fmt.Sprintf("Invalid request: vehicle %v already have an objectid %v", v.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	tx.Model().Vehicles().Save(&vehicle)
	err = tx.Commit()
	if err != nil {
		logger.Log.Debugf("Transaction error: %v", err)
		http.Error(response, "Internal error", http.StatusInternalServerError)
		return
	}
	jsonBytes, _ := vehicle.MarshalJSON()
	response.Write(jsonBytes)
}
