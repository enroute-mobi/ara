package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
)

type VehicleController struct {
	referential *core.Referential
}

func NewVehicleController(referential *core.Referential) RestfulResource {
	return &VehicleController{
		referential: referential,
	}
}

func (controller *VehicleController) findVehicle(identifier string) (*model.Vehicle, bool) {
	foundStrings := idPattern.FindStringSubmatch(identifier)
	if foundStrings != nil {
		code := model.NewCode(foundStrings[1], foundStrings[2])
		return controller.referential.Model().Vehicles().FindByCode(code)
	}
	return controller.referential.Model().Vehicles().Find(model.VehicleId(identifier))
}

func (controller *VehicleController) Index(response http.ResponseWriter, params url.Values) {
	logger.Log.Debugf("Vehicles Index")

	stime := controller.referential.Clock().Now()
	allVehicles := controller.referential.Model().Vehicles().FindAll()
	logger.Log.Debugf("VehicleController FindAll time : %v", controller.referential.Clock().Since(stime))

	direction := params.Get("direction")
	switch direction {
	case "desc":
		sort.Slice(allVehicles, func(i, j int) bool {
			return allVehicles[i].Id() > allVehicles[j].Id()
		})
	case "asc", "":
		sort.Slice(allVehicles, func(i, j int) bool {
			return allVehicles[i].Id() < allVehicles[j].Id()
		})
	default:
		http.Error(response, fmt.Sprintf("invalid request: query parameter \"direction\": %s", params.Get("direction")), http.StatusBadRequest)
		return
	}

	paginatedVehicles, err := paginate(allVehicles, params)
	if err != nil {
		http.Error(response, err.Error(), http.StatusBadRequest)
		return
	}

	jsonBytes, _ := json.Marshal(paginatedVehicles)
	response.Write(jsonBytes)

	logger.Log.Debugf("VehicleController Json Marshal time : %v ", controller.referential.Clock().Since(stime))
}

func (controller *VehicleController) Show(response http.ResponseWriter, identifier string, _params url.Values) {
	vehicle, ok := controller.findVehicle(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Vehicle not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Get vehicle %s", identifier)

	jsonBytes, _ := vehicle.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *VehicleController) Delete(response http.ResponseWriter, identifier string) {
	vehicle, ok := controller.findVehicle(identifier)
	if !ok {
		http.Error(response, fmt.Sprintf("Vehicle not found: %s", identifier), http.StatusNotFound)
		return
	}
	logger.Log.Debugf("Delete vehicle %s", identifier)

	jsonBytes, _ := vehicle.MarshalJSON()
	controller.referential.Model().Vehicles().Delete(vehicle)
	response.Write(jsonBytes)
}

func (controller *VehicleController) Update(response http.ResponseWriter, identifier string, body []byte) {
	vehicle, ok := controller.findVehicle(identifier)
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

	for _, obj := range vehicle.Codes() {
		v, ok := controller.referential.Model().Vehicles().FindByCode(obj)
		if ok && v.Id() != vehicle.Id() {
			http.Error(response, fmt.Sprintf("Invalid request: vehicle %v already have a code %v", v.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	controller.referential.Model().Vehicles().Save(vehicle)
	jsonBytes, _ := vehicle.MarshalJSON()
	response.Write(jsonBytes)
}

func (controller *VehicleController) Create(response http.ResponseWriter, body []byte) {
	logger.Log.Debugf("Create vehicle: %s", string(body))

	vehicle := controller.referential.Model().Vehicles().New()

	err := json.Unmarshal(body, &vehicle)
	if err != nil {
		http.Error(response, fmt.Sprintf("Invalid request: can't parse request body: %v", err), http.StatusBadRequest)
		return
	}

	if vehicle.Id() != "" {
		http.Error(response, "Invalid request", http.StatusBadRequest)
		return
	}

	for _, obj := range vehicle.Codes() {
		v, ok := controller.referential.Model().Vehicles().FindByCode(obj)
		if ok {
			http.Error(response, fmt.Sprintf("Invalid request: vehicle %v already have a code %v", v.Id(), obj.String()), http.StatusBadRequest)
			return
		}
	}

	controller.referential.Model().Vehicles().Save(vehicle)
	jsonBytes, _ := vehicle.MarshalJSON()
	response.Write(jsonBytes)
}
